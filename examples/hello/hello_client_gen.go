package hello

import (
	"github.com/anqur/gbio"
)

type Tx struct {
	gbio.Tx
}

func Do() *Tx {
	return &Tx{Tx: gbio.Tx{Cl: gbio.DefaultClient}}
}

func With(cl *gbio.Client) *Tx { return &Tx{Tx: gbio.Tx{Cl: cl}} }

func (c *Tx) SayHi(req *SelfIntro) *OkReply {
	httpReq, err := c.Request(
		serviceKey,
		"/Greeting/SayHi",
		&SelfIntroEncoder{req},
	)
	if err != nil {
		c.Error = err
		return nil
	}

	httpResp, err := c.Cl.HttpClient().Do(httpReq)
	if err != nil {
		c.Error = err
		return nil
	}

	resp, err := NewDecoder(httpResp.Body, httpResp.Header).OkReply()
	if err != nil {
		c.Error = err
		return nil
	}

	return resp
}

func (c *Tx) HiAdmin(req *ImAdmin) Reply {
	httpReq, err := c.Request(
		serviceKey,
		"/Admin/HiAdmin",
		&ImAdminEncoder{req},
	)
	if err != nil {
		c.Error = err
		return nil
	}

	httpResp, err := c.Cl.HttpClient().Do(httpReq)
	if err != nil {
		c.Error = err
		return nil
	}

	resp, err := NewDecoder(httpResp.Body, httpResp.Header).Reply()
	if err != nil {
		c.Error = err
		return nil
	}

	return resp
}
