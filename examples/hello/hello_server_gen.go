package hello

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type helloMux struct {
	http.ServeMux

	s Hello
}

func (*helloMux) ServiceName() string { return "hello.Hello" }

type discriminator struct {
	Tag string `json:"_t"`
}

func unmarshal(r *http.Request) (Greeting, error) {
	defer func() { _ = r.Body.Close() }()
	d, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var tag discriminator
	if err := json.Unmarshal(d, &tag); err != nil {
		return nil, err
	}
	var req Greeting
	switch tag.Tag {
	case "JustHi":
		var v JustHi
		err = json.Unmarshal(d, &v)
		req = &v
		v.BaseGreeting.ReqID = r.Header.Get("x-request-id")
	case "SelfIntro":
		var v SelfIntro
		err = json.Unmarshal(d, &v)
		req = &v
		v.BaseGreeting.ReqID = r.Header.Get("x-request-id")
	default:
		return nil, fmt.Errorf("unknown message tag %q", tag.Tag)
	}
	return req, err
}

func internalServerError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte(fmt.Sprintf("%q", err.Error())))
}

func ok(w http.ResponseWriter, v interface{}) {
	d, _ := json.Marshal(v)
	_, _ = w.Write(d)
	w.Header().Add("Content-Type", "application/json")
}

func (d *helloMux) SayHi(r Greeting) *Reply {
	return d.s.SayHi(r)
}

func Mux(s Hello) http.Handler {
	d := &helloMux{s: s}
	d.HandleFunc("/SayHi", func(w http.ResponseWriter, r *http.Request) {
		req, err := unmarshal(r)
		if err != nil {
			internalServerError(w, err)
			return
		}
		ok(w, d.s.SayHi(req))
	})
	return d
}

func NewServer(s Hello, addr string) *http.Server {
	return &http.Server{Addr: addr, Handler: Mux(s)}
}
