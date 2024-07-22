package elastic

import (
	"github.com/olivere/elastic"
	"github.com/unchainio/pkg/errors"
)

type Client struct {
	*elastic.Client
	cfg *Config
}

func NewClient(cfg *Config) (c *Client, err error) {
	//url := "http://37.139.18.39:9200"
	client, err := elastic.NewSimpleClient(elastic.SetURL(cfg.URL), elastic.SetBasicAuth(cfg.User, cfg.Pass))
	if err != nil {
		return nil, errors.Wrap(err, "falied to initialize elastic client")
	}
	return &Client{client, cfg}, nil
}
