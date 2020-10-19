package client

import (
	"net/url"
	"path"
)

// Client is the presto client
type Client struct {
	url string
}

// New creates a new presto client
func New(url string) Client {
	return Client{
		url: url,
	}
}

func (c Client) withPath(fragment string) (string, error) {
	u, err := url.Parse(c.url)
	if err != nil {
		return "", err
	}
	u.Path = path.Join(u.Path, fragment)
	return u.String(), nil
}
