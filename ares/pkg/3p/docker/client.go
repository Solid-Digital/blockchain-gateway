package docker

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/docker/cli/cli/streams"

	"github.com/docker/docker/api/types"

	dc "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/unchainio/interfaces/logger"
	"github.com/unchainio/pkg/iferr"
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	*dc.Client
	cfg *Config
	log logger.Logger
}

func NewClient(log logger.Logger, cfg *Config) (*Client, func(), error) {
	log.Debugf("Initializing docker...\n")
	var c *http.Client

	if cfg.TLS != nil {
		cert, err := cfg.TLS.X509KeyPair()
		if err == nil {
			tlsConfig := &tls.Config{
				Certificates:       []tls.Certificate{cert},
				InsecureSkipVerify: true,
			}

			c = &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: tlsConfig,
				},
			}
		}
	}

	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}

	dockerClient, err := dc.NewClientWithOpts(dc.WithAPIVersionNegotiation(), dc.WithHost(cfg.Host), dc.WithHTTPClient(c), dc.WithHTTPHeaders(defaultHeaders)) // Engine 1.12

	if err != nil {
		return nil, nil, err
	}

	log.Debugf("connected to docker on %s\n", cfg.Host)

	client := &Client{
		Client: dockerClient,
		cfg:    cfg,
		log:    log,
	}

	cleanup := func() {
		iferr.Warn(client.Close())
	}

	return client, cleanup, nil
}

func (c *Client) ImagePush(ctx context.Context, image string, options types.ImagePushOptions) (io.ReadCloser, error) {
	responseBody, err := c.Client.ImagePush(ctx, image, options)

	if err != nil {
		return nil, err
	}

	defer responseBody.Close()
	buf := new(bytes.Buffer)

	tr := io.TeeReader(responseBody, buf)

	out := streams.NewOut(os.Stdout)

	err = jsonmessage.DisplayJSONMessagesToStream(tr, out, nil)

	if err != nil {
		return nil, err
	}

	return ioutil.NopCloser(buf), err
}

func (c *Client) ImagePull(ctx context.Context, image string, options types.ImagePullOptions) (io.ReadCloser, error) {
	responseBody, err := c.Client.ImagePull(ctx, image, options)

	if err != nil {
		return nil, err
	}

	defer responseBody.Close()
	buf := new(bytes.Buffer)

	tr := io.TeeReader(responseBody, buf)

	out := streams.NewOut(os.Stdout)

	err = jsonmessage.DisplayJSONMessagesToStream(tr, out, nil)

	if err != nil {
		return nil, err
	}

	return ioutil.NopCloser(buf), err
}
