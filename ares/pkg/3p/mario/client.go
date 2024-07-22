package mario

import (
	"bitbucket.org/unchain/ares/pkg/ares"
	"github.com/google/wire"
	"github.com/unchainio/interfaces/logger"
	"github.com/unchainio/pkg/errors"

	"github.com/go-resty/resty/v2"
)

var ImageBuilderSet = wire.NewSet(NewClient, wire.Bind(new(ares.ImageBuilder), new(Client)))

// Component defines a single component filename and location.
type Config struct {
	URL string
}

type Client struct {
	http *resty.Client
	log  logger.Logger
	cfg  *Config
}

// NewClient returns a new instantiated client
func NewClient(log logger.Logger, cfg *Config) *Client {
	client := resty.New()
	client.SetHostURL(cfg.URL)

	return &Client{
		http: client,
		log:  log,
		cfg:  cfg,
	}
}

// BuildImage manifest *BuildManifest
func (c *Client) BuildImage(manifest *ares.BuildManifest) error {
	r, err := c.http.R().SetBody(manifest).Post("build/")
	if err != nil {
		return errors.Wrap(err, "")
	}

	if r.IsError() {
		return errors.New(string(r.Body()))
	}

	c.log.Debugf(string(r.Body()))

	return nil
}
