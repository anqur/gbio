package hello

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/anqur/gbio"
)

const serviceKey = "hello.Hello"

type tagger struct {
	Tag int `json:"_t"`
}

type Decoder struct {
	r   io.ReadCloser
	ctx http.Header
}

func NewDecoder(r io.ReadCloser, ctx http.Header) *Decoder {
	return &Decoder{r: r, ctx: ctx}
}

type SelfIntroEncoder struct {
	m *SelfIntro
}

func (e *SelfIntroEncoder) Marshal() ([]byte, http.Header, error) {
	d, err := json.Marshal(e.m)
	if err != nil {
		return nil, nil, err
	}
	return d, nil, nil
}

func (d *Decoder) SelfIntro() (*SelfIntro, error) {
	defer func() { _ = d.r.Close() }()
	buf, err := io.ReadAll(d.r)
	if err != nil {
		return nil, err
	}
	var ret SelfIntro
	if err := json.Unmarshal(buf, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}

type ImAdminEncoder struct {
	m *ImAdmin
}

func (r *ImAdminEncoder) Marshal() ([]byte, http.Header, error) {
	d, err := json.Marshal(r.m)
	if err != nil {
		return nil, nil, err
	}
	ctx := make(http.Header, 1)
	ctx.Add("Authorization", r.m.Authorization)
	return d, ctx, nil
}

func (d *Decoder) ImAdmin() (*ImAdmin, error) {
	defer func() { _ = d.r.Close() }()
	buf, err := io.ReadAll(d.r)
	if err != nil {
		return nil, err
	}
	var ret ImAdmin
	if err := json.Unmarshal(buf, &ret); err != nil {
		return nil, err
	}
	ret.Authorization = d.ctx.Get("Authorization")
	return &ret, nil
}

type OkReplyEncoder struct {
	m *OkReply
}

func (e *OkReplyEncoder) Marshal() ([]byte, http.Header, error) {
	d, err := json.Marshal(e.m)
	if err != nil {
		return nil, nil, err
	}
	return d, nil, nil
}

func (d *Decoder) OkReply() (*OkReply, error) {
	defer func() { _ = d.r.Close() }()
	buf, err := io.ReadAll(d.r)
	if err != nil {
		return nil, err
	}
	var ret OkReply
	if err := json.Unmarshal(buf, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}

type ReplyEncoder struct {
	m Reply
}

type taggedOkReply struct {
	tagger
	*OkReply
}

type taggedErrReply struct {
	tagger
	*ErrReply
}

func (e *ReplyEncoder) Marshal() ([]byte, http.Header, error) {
	var (
		d   []byte
		err error
	)
	tag := tagger{Tag: e.m.isReply()}
	switch v := e.m.(type) {
	case OkReply:
		tagged := &taggedOkReply{tagger: tag, OkReply: &v}
		d, err = json.Marshal(tagged)
	case *OkReply:
		tagged := &taggedOkReply{tagger: tag, OkReply: v}
		d, err = json.Marshal(tagged)

	case ErrReply:
		tagged := &taggedErrReply{tagger: tag, ErrReply: &v}
		d, err = json.Marshal(tagged)
	case *ErrReply:
		tagged := &taggedErrReply{tagger: tag, ErrReply: v}
		d, err = json.Marshal(tagged)

	default:
		err = fmt.Errorf("%w: %+v", gbio.ErrCodecBadMsgType, v)
	}
	return d, nil, err
}

func (d *Decoder) Reply() (Reply, error) {
	defer func() { _ = d.r.Close() }()
	buf, err := io.ReadAll(d.r)
	if err != nil {
		return nil, err
	}

	var tag tagger
	if err := json.Unmarshal(buf, &tag); err != nil {
		return nil, err
	}

	var resp Reply
	switch t := tag.Tag; t {
	case 1:
		var resp0 OkReply
		err = json.Unmarshal(buf, &resp0)
		resp = &resp0

	case 2:
		var resp0 ErrReply
		err = json.Unmarshal(buf, &resp0)
		resp = &resp0

	default:
		err = fmt.Errorf("%w: %q", gbio.ErrCodecBadMsgTag, t)
	}
	if err != nil {
		return nil, err
	}

	return resp, err
}
