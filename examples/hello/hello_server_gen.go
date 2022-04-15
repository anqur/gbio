package hello

import (
	"fmt"
	"net/http"

	"github.com/anqur/gbio/pkg/endpoints"
	"github.com/anqur/gbio/pkg/logging"
	"github.com/anqur/gbio/pkg/servers"
)

func internalServerError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func ok(w http.ResponseWriter, d []byte, ctx http.Header) {
	w.Header().Add("Content-Type", "application/json")
	for k, vs := range ctx {
		for _, v := range vs {
			w.Header().Add(k, v)
		}
	}
	_, _ = w.Write(d)
}

func RegisterGreeting(i Greeting, opts ...endpoints.Option) servers.Option {
	ep := new(servers.Endpoint)
	ep.Tag = endpoints.DefaultTag

	for _, opt := range opts {
		opt(&ep.Endpoint)
	}

	ep.Name = fmt.Sprintf("hello.Greeting/%s", ep.Tag)
	ep.BaseURI = fmt.Sprintf("/Greeting/%s/SayHi", ep.Tag)
	ep.Handler = func(w http.ResponseWriter, r *http.Request) {
		req, err := NewDecoder(r.Body, r.Header).SelfIntro()
		if err != nil {
			internalServerError(w, err)
			return
		}
		d, ctx, err := (&OkReplyEncoder{i.SayHi(req)}).Marshal()
		if err != nil {
			internalServerError(w, err)
			return
		}
		ok(w, d, ctx)
		logging.Info.Println("Access:", r.RemoteAddr, r.RequestURI)
	}

	return func(s *servers.Server) { s.Register(ep) }
}

func RegisterAdmin(i Admin, opts ...endpoints.Option) servers.Option {
	ep := new(servers.Endpoint)
	ep.Tag = endpoints.DefaultTag

	for _, opt := range opts {
		opt(&ep.Endpoint)
	}

	ep.Name = fmt.Sprintf("hello.Admin/%s", ep.Tag)
	ep.BaseURI = fmt.Sprintf("/Admin/%s/HiAdmin", ep.Tag)
	ep.Handler = func(w http.ResponseWriter, r *http.Request) {
		req, err := NewDecoder(r.Body, r.Header).ImAdmin()
		if err != nil {
			internalServerError(w, err)
			return
		}
		d, ctx, err := (&ReplyEncoder{i.HiAdmin(req)}).Marshal()
		if err != nil {
			internalServerError(w, err)
			return
		}
		ok(w, d, ctx)
		logging.Info.Println("Access:", r.RemoteAddr, r.RequestURI)
	}

	return func(s *servers.Server) { s.Register(ep) }
}
