package harbor

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"net/url"

	"fmt"

	"strings"

	"github.com/unchainio/pkg/errors"
	"github.com/unchainio/pkg/xapi"
)

type Client struct {
	cfg    *Config
	client *xapi.Client
}

func NewClient(cfg *Config) (*Client, error) {
	// Done so that other services can check if harbor is enabled or not via some form of `if harbor == nil {}`
	if cfg == nil {
		return nil, nil
	}

	tp := &xapi.BasicAuthTransport{
		Username: cfg.Username,
		Password: cfg.Password,
	}

	client, err := xapi.NewClient(
		cfg.URL,
		xapi.WithClient(tp.Client()),
	)

	if err != nil {
		return nil, err
	}

	return &Client{
		cfg:    cfg,
		client: client,
	}, nil
}

func (c *Client) CreateProject(name string, public bool) (int, error) {
	var err error

	data := map[string]interface{}{
		"project_name": name,
		"metadata": map[string]interface{}{
			"public": strconv.FormatBool(public),
		},
	}

	req, err := c.client.NewRequest(http.MethodPost, "/api/projects", data)

	if err != nil {
		return 0, err
	}

	out := make(map[string]json.RawMessage)

	resp, cleanup, err := c.client.Do(context.Background(), req, &out)
	if err != nil {
		return 0, err
	}
	defer cleanup()

	location, err := resp.Location()

	if err != nil {
		return 0, errors.Wrap(err, "failed to get location")
	}

	idStr := strings.TrimPrefix(location.Path, "/api/projects/")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		return 0, errors.Wrap(err, "failed to parse ID")
	}

	return id, nil
}

func (c *Client) UpsertProject(name string, public bool) (int, error) {
	p, err := c.FindProjects().Name(name).One()

	if err != nil {
		return c.CreateProject(name, public)
	}

	return p.ID, nil
}

type FindProjectsQuery struct {
	client *Client

	name     string
	public   *bool
	owner    string
	page     *int
	pageSize *int
}

func (q *FindProjectsQuery) String() string {
	b := strings.Builder{}

	if q.name != "" {
		b.WriteString("name: '")
		b.WriteString(q.name)
		b.WriteString("'; ")
	}

	if q.public != nil {
		b.WriteString("public: '")
		b.WriteString(strconv.FormatBool(*q.public))
		b.WriteString("'; ")
	}

	if q.owner != "" {
		b.WriteString("owner: '")
		b.WriteString(q.owner)
		b.WriteString("'; ")
	}

	if q.page != nil {
		b.WriteString("page: '")
		b.WriteString(strconv.FormatInt(int64(*q.page), 10))
		b.WriteString("'; ")
	}

	if q.pageSize != nil {
		b.WriteString("pageSize: '")
		b.WriteString(strconv.FormatInt(int64(*q.pageSize), 10))
		b.WriteString("'; ")
	}

	return b.String()
}

func (q *FindProjectsQuery) Name(name string) *FindProjectsQuery {
	q.name = name

	return q
}

func (q *FindProjectsQuery) Public(public bool) *FindProjectsQuery {
	q.public = &public

	return q
}

func (q *FindProjectsQuery) Owner(owner string) *FindProjectsQuery {
	q.owner = owner

	return q
}

func (q *FindProjectsQuery) Page(page int) *FindProjectsQuery {
	q.page = &page

	return q
}

func (q *FindProjectsQuery) PageSize(pageSize int) *FindProjectsQuery {
	q.pageSize = &pageSize

	return q
}

func (c *Client) FindProjects() *FindProjectsQuery {
	return &FindProjectsQuery{
		client: c,
	}
}

func (q *FindProjectsQuery) All() ([]Project, error) {
	var err error

	path, err := url.Parse("/api/projects")
	params := path.Query()

	if q.name != "" {
		params.Set("name", q.name)
	}

	if q.public != nil {
		params.Set("public", strconv.FormatBool(*q.public))
	}

	if q.owner != "" {
		params.Set("owner", q.owner)
	}

	if q.page != nil {
		params.Set("page", strconv.FormatInt(int64(*q.page), 10))
	}

	if q.pageSize != nil {
		params.Set("page_size", strconv.FormatInt(int64(*q.pageSize), 10))
	}

	path.RawQuery = params.Encode()

	req, err := q.client.client.NewRequest(http.MethodGet, path.String(), nil)

	if err != nil {
		return nil, err
	}

	var projects []Project

	_, cleanup, err := q.client.client.Do(context.Background(), req, &projects)
	if err != nil {
		return nil, err
	}
	defer cleanup()

	return projects, nil
}

func (q *FindProjectsQuery) One() (*Project, error) {
	var err error

	projects, err := q.All()

	if err != nil {
		return nil, err
	}

	for _, project := range projects {
		if project.Name == q.name {
			return &project, nil
		}
	}

	return nil, errors.Errorf("failed to find project with query: %s", q.String())
}

func (c *Client) RemoveProject(name string) error {
	var err error

	project, err := c.FindProjects().Name(name).One()

	if err != nil {
		return err
	}

	path := fmt.Sprintf("/api/projects/%d", project.ID)

	req, err := c.client.NewRequest(http.MethodDelete, path, nil)

	if err != nil {
		return err
	}

	_, cleanup, err := c.client.Do(context.Background(), req, nil)

	if err != nil {
		return err
	}

	defer cleanup()

	return nil
}
