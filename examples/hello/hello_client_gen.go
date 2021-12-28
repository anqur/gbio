package hello

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	urls "net/url"
)

type Client struct {
	h *http.Client
	u *urls.URL
}

func NewClient(url string) (*Client, error) {
	u, err := urls.Parse(url)
	if err != nil {
		return nil, err
	}
	return &Client{
		// TODO: Options to pass the HTTP client.
		h: http.DefaultClient,
		u: u,
	}, nil
}

func (c *Client) Call() *Call {
	return &Call{
		cl: c,
		u:  &urls.URL{Scheme: c.u.Scheme, Host: c.u.Host},
	}
}

type Call struct {
	cl    *Client
	u     *urls.URL
	Error error
}

type taggedJustHi struct {
	discriminator
	*JustHi
}

type taggedSelfIntro struct {
	discriminator
	*SelfIntro
}

func (c *Call) SayHi(greeting Greeting) *Reply {
	var d []byte
	switch r := greeting.(type) {
	case JustHi:
		tagged := &taggedJustHi{JustHi: &r}
		tagged.Tag = "JustHi"
		d, c.Error = json.Marshal(tagged)
	case *JustHi:
		tagged := &taggedJustHi{JustHi: r}
		tagged.Tag = "JustHi"
		d, c.Error = json.Marshal(tagged)

	case SelfIntro:
		tagged := &taggedSelfIntro{SelfIntro: &r}
		tagged.Tag = "SelfIntro"
		d, c.Error = json.Marshal(tagged)
	case *SelfIntro:
		tagged := &taggedSelfIntro{SelfIntro: r}
		tagged.Tag = "SelfIntro"
		d, c.Error = json.Marshal(tagged)
	}
	if c.Error != nil {
		return nil
	}

	// TODO: Marshal context.

	c.u.Path = "/SayHi"
	var r *http.Response
	r, c.Error = c.cl.h.Post(
		c.u.String(),
		"application/json",
		bytes.NewReader(d),
	)
	if c.Error != nil {
		return nil
	}

	defer func() { _ = r.Body.Close() }()
	d, c.Error = io.ReadAll(r.Body)
	if c.Error != nil {
		return nil
	}

	var resp Reply
	if c.Error = json.Unmarshal(d, &resp); c.Error != nil {
		return nil
	}

	return &resp
}
