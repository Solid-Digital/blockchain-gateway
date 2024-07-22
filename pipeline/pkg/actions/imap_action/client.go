package imap_action

import (
	"fmt"
	"github.com/emersion/go-imap"
	move "github.com/emersion/go-imap-move"
	imapclient "github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"github.com/pkg/errors"
	"github.com/prometheus/common/log"
	"github.com/unchainio/interfaces/logger"
	"io"
	"io/ioutil"
)

type Client struct {
	logger logger.Logger
	cfg    *Config
	Client *client
}

type client struct {
	*imapclient.Client
	*move.MoveClient
}

func NewClient(logger logger.Logger, cfg *Config) (*Client, error) {
	c, err := imapclient.DialTLS(fmt.Sprintf("%s%s", cfg.Domain, cfg.Port), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to dial imap server")
	}

	client := &Client{
		logger: logger,
		cfg:    cfg,
		Client: &client{
			Client:     c,
			MoveClient: move.NewClient(c),
		},
	}

	err = c.Login(cfg.Username, cfg.Password)
	if err != nil {
		return nil, errors.Wrap(err, "failed to login")
	}

	err = client.createFailedMailbox()
	if err != nil {
		return nil, err
	}

	mbox, err := client.Client.Select("INBOX", false)
	if err != nil {
		return nil, err
	}
	logger.Debugf("Started client - mailbox contains %v messages", mbox.Messages)

	return client, nil
}

func (c *Client) createFailedMailbox() error {
	ch := make(chan *imap.MailboxInfo, 1)

	err := c.Client.List("", "Failed", ch)
	if err != nil {
		return errors.Wrap(err, "")
	}

	mailbox := <-ch

	// mailbox Failed already exists
	if mailbox != nil {
		c.logger.Debugf("mailbox 'Failed' already exists")
		return nil
	}

	c.logger.Debugf("Creating mailbox 'Failed'")

	err = c.Client.Create("Failed")
	if err != nil {
		return errors.Wrap(err, "failed to create inbox 'Failed'")
	}

	return nil
}

func (c *Client) MarkMessageAsRead(seqNum uint32) error {
	seqSet := new(imap.SeqSet)
	seqSet.AddNum(seqNum)

	item := imap.FormatFlagsOp(imap.AddFlags, true)
	flags := []interface{}{imap.SeenFlag}

	err := c.Client.Store(seqSet, item, flags, nil)
	if err != nil {
		return errors.Wrapf(err, "could not mark mail with seqNum %v as seen", seqNum)
	}

	return nil
}

func (c *Client) MoveFailedMessage(seqNum uint32) error {
	seqSet := new(imap.SeqSet)
	seqSet.AddNum(seqNum)

	err := c.Client.MoveWithFallback(seqSet, "Failed")
	if err != nil {
		return errors.Wrapf(err, "could not move mail with seqNum %v to the failed inbox", seqNum)
	}

	return nil
}

func (c *Client) GetNewMessageAttachments() (map[string]interface{}, error) {
	// Create search criteria based on from address and subject
	criteria := &imap.SearchCriteria{
		Text:         []string{c.cfg.FromFilter, c.cfg.SubjectFilter},
		WithoutFlags: []string{imap.SeenFlag},
	}

	seqNums, err := c.Client.Search(criteria)
	if err != nil {
		c.logger.Errorf("%v", err)
		return nil, errors.Wrap(err, "error while searching")
	}
	c.logger.Printf("New mail(s) with seq#: %v", seqNums)

	if len(seqNums) == 0 {
		return nil, nil
	}

	seqset := new(imap.SeqSet)
	seqset.AddNum(seqNums...)

	imapMessages := make(chan *imap.Message, 10)

	var section imap.BodySectionName
	section.Peek = true
	items := []imap.FetchItem{section.FetchItem()}

	go func() {
		err := c.Client.Fetch(seqset, items, imapMessages)
		if err != nil {
			log.Errorf("%v", err)
		}
	}()

	output := map[string]interface{}{
		"messages": map[uint32]interface{}{},
	}

	// handle individual messages
	for msg := range imapMessages {
		if msg == nil {
			c.logger.Errorf("could not find a msg")
			return nil, err
		}

		r := msg.GetBody(&section)
		if r == nil {
			c.logger.Errorf("server did not return msg body")
			return nil, err
		}

		mr, err := mail.CreateReader(r)
		if err != nil {
			c.logger.Errorf("could not create mail reader")
			return nil, err
		}

		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			} else if err != nil {
				c.logger.Errorf("error reading next part of message reader")
				return nil, err
			}

			switch h := p.Header.(type) {
			case *mail.AttachmentHeader:
				filename, _ := h.Filename()
				c.logger.Debugf("Found attachment: %v", filename)

				bodyBytes, err := ioutil.ReadAll(p.Body)
				if err != nil {
					return output, err
				}

				// handle individual message's bodybytes
				messages := output["messages"].(map[uint32]interface{})
				messages[msg.SeqNum] = bodyBytes
				output["messages"] = messages
			default:
				continue
			}
		}
	}

	return output, nil
}
