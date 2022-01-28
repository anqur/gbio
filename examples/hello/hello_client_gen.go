package hello

import (
	"fmt"

	"github.com/anqur/gbio/clients"
	"github.com/anqur/gbio/endpoints"
)

type (
	tx struct{ cl *clients.Client }

	greetingTx struct{ *clients.Endpoint }
	adminTx    struct{ *clients.Endpoint }
)

var Tx = &tx{cl: clients.Default}

func With(cl *clients.Client) *tx { return &tx{cl: cl} }

func (t *tx) Greeting(opts ...endpoints.Option) *greetingTx {
	ep := &clients.Endpoint{Cl: t.cl}
	ep.Tag = endpoints.DefaultTag

	for _, opt := range opts {
		opt(&ep.Endpoint)
	}

	ep.Name = fmt.Sprintf("hello.Greeting/%s", ep.Tag)
	ep.BaseURI = fmt.Sprintf("/Greeting/%s/SayHi", ep.Tag)
	return &greetingTx{Endpoint: ep}
}

func (t *tx) Admin(opts ...endpoints.Option) *adminTx {
	ep := &clients.Endpoint{Cl: t.cl}
	ep.Tag = endpoints.DefaultTag

	for _, opt := range opts {
		opt(&ep.Endpoint)
	}

	ep.Name = fmt.Sprintf("hello.Admin/%s", ep.Tag)
	ep.BaseURI = fmt.Sprintf("/Admin/%s/HiAdmin", ep.Tag)
	return &adminTx{Endpoint: ep}
}

func (t *greetingTx) SayHi(req *SelfIntro) *OkReply {
	httpReq, err := t.Cl.Request(t.Name, t.BaseURI, &SelfIntroEncoder{req})
	if err != nil {
		t.Error = err
		return nil
	}

	httpResp, err := t.Cl.HttpClient().Do(httpReq)
	if err != nil {
		t.Error = err
		return nil
	}

	resp, err := NewDecoder(httpResp.Body, httpResp.Header).OkReply()
	if err != nil {
		t.Error = err
		return nil
	}

	return resp
}

func (t *adminTx) HiAdmin(req *ImAdmin) Reply {
	httpReq, err := t.Cl.Request(t.Name, t.BaseURI, &ImAdminEncoder{req})
	if err != nil {
		t.Error = err
		return nil
	}

	httpResp, err := t.Cl.HttpClient().Do(httpReq)
	if err != nil {
		t.Error = err
		return nil
	}

	resp, err := NewDecoder(httpResp.Body, httpResp.Header).Reply()
	if err != nil {
		t.Error = err
		return nil
	}

	return resp
}
