package hello

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/anqur/gbio"
)

func Do() *Doer {
	return &Doer{cl: gbio.DefaultClient}
}

func With(cl *gbio.Client) *Doer {
	return &Doer{cl: cl}
}

type Doer struct {
	cl    *gbio.Client
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

func (c *Doer) SayHi(greeting Greeting) *Reply {
	// TODO: Impose code generation on marshalling context.
	ctx := make(map[string]string)

	var d []byte
	switch r := greeting.(type) {
	case JustHi:
		tagged := &taggedJustHi{JustHi: &r}
		tagged.Tag = "JustHi"
		d, c.Error = json.Marshal(tagged)
		ctx["x-request-id"] = r.ReqID
	case *JustHi:
		tagged := &taggedJustHi{JustHi: r}
		tagged.Tag = "JustHi"
		d, c.Error = json.Marshal(tagged)
		ctx["x-request-id"] = r.ReqID

	case SelfIntro:
		tagged := &taggedSelfIntro{SelfIntro: &r}
		tagged.Tag = "SelfIntro"
		d, c.Error = json.Marshal(tagged)
		ctx["x-request-id"] = r.ReqID
	case *SelfIntro:
		tagged := &taggedSelfIntro{SelfIntro: r}
		tagged.Tag = "SelfIntro"
		d, c.Error = json.Marshal(tagged)
		ctx["x-request-id"] = r.ReqID
	}
	if c.Error != nil {
		return nil
	}

	var url string
	url, c.Error = c.cl.LookupEndpoint("hello.Hello")
	if c.Error != nil {
		return nil
	}

	var req *http.Request
	req, c.Error = http.NewRequest(
		http.MethodPost,
		url+"/SayHi",
		bytes.NewReader(d),
	)
	if c.Error != nil {
		return nil
	}

	req.Header.Add("Content-Type", "application/json")
	for k, v := range ctx {
		req.Header.Add(k, v)
	}

	var r *http.Response
	r, c.Error = c.cl.HttpClient().Do(req)
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
