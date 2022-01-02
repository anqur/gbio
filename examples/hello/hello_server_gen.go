package hello

import (
	"fmt"
	"net/http"
)

type helloMux struct {
	http.ServeMux

	s Hello
}

func (*helloMux) ServiceName() string { return serviceKey }

type discriminator struct {
	Tag string `json:"_t"`
}

func internalServerError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte(fmt.Sprintf("%q", err.Error())))
}

func ok(w http.ResponseWriter, d []byte, ctx http.Header) {
	w.Header().Add("Content-Type", "application/json")
	_, _ = w.Write(d)
	for k, vs := range ctx {
		for _, v := range vs {
			w.Header().Add(k, v)
		}
	}
}

func (d *helloMux) SayHi(r *SelfIntro) *OkReply {
	return d.s.SayHi(r)
}

func (d *helloMux) HiAdmin(r *ImAdmin) Reply {
	return d.s.HiAdmin(r)
}

func Mux(s Hello) http.Handler {
	d := &helloMux{s: s}
	d.HandleFunc("/Greeting/SayHi", func(w http.ResponseWriter, r *http.Request) {
		req, err := NewDecoder(r.Body, r.Header).SelfIntro()
		if err != nil {
			internalServerError(w, err)
			return
		}
		d, ctx, err := (&OkReplyEncoder{d.s.SayHi(req)}).Marshal()
		if err != nil {
			internalServerError(w, err)
			return
		}
		ok(w, d, ctx)
	})
	d.HandleFunc("/Admin/HiAdmin", func(w http.ResponseWriter, r *http.Request) {
		req, err := NewDecoder(r.Body, r.Header).ImAdmin()
		if err != nil {
			internalServerError(w, err)
			return
		}
		d, ctx, err := (&ReplyEncoder{d.s.HiAdmin(req)}).Marshal()
		if err != nil {
			internalServerError(w, err)
			return
		}
		ok(w, d, ctx)
	})
	return d
}
