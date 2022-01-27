package hello

import (
	"fmt"

	"github.com/anqur/gbio"
)

type (
	tx struct{ cl *gbio.Client }

	greetingTx struct{ *gbio.ClientEndpoint }
	adminTx    struct{ *gbio.ClientEndpoint }
)

var Tx = &tx{cl: gbio.DefaultClient}

func With(cl *gbio.Client) *tx { return &tx{cl: cl} }

func (t *tx) Greeting(opts ...gbio.EndpointOption) *greetingTx {
	ep := &gbio.ClientEndpoint{Cl: t.cl}
	ep.Tag = gbio.DefaultTag

	for _, opt := range opts {
		opt(&ep.Endpoint)
	}

	ep.Name = fmt.Sprintf("hello.Greeting/%s", ep.Tag)
	ep.BaseURI = fmt.Sprintf("/Greeting/%s/SayHi", ep.Tag)
	return &greetingTx{ClientEndpoint: ep}
}

func (t *tx) Admin(opts ...gbio.EndpointOption) *adminTx {
	ep := &gbio.ClientEndpoint{Cl: t.cl}
	ep.Tag = gbio.DefaultTag

	for _, opt := range opts {
		opt(&ep.Endpoint)
	}

	ep.Name = fmt.Sprintf("hello.Admin/%s", ep.Tag)
	ep.BaseURI = fmt.Sprintf("/Admin/%s/HiAdmin", ep.Tag)
	return &adminTx{ClientEndpoint: ep}
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
