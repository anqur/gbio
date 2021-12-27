package hello

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type helloDelegate struct {
	http.ServeMux

	s Hello
}

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
	case "SelfIntro":
		var v SelfIntro
		err = json.Unmarshal(d, &v)
		req = &v
	default:
		return nil, fmt.Errorf("unknown message tag %q", tag.Tag)
	}
	return req, err
}

func internalServerError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(fmt.Sprintf("%q", err.Error())))
}

func ok(w http.ResponseWriter, v interface{}) {
	d, _ := json.Marshal(v)
	w.Write(d)
	w.Header().Add("Content-Type", "application/json")
}

func (d *helloDelegate) SayHi(r Greeting) *Reply {
	return d.s.SayHi(r)
}

func NewServer(s Hello, addr string) *http.Server {
	d := &helloDelegate{s: s}
	d.HandleFunc("/SayHi", func(w http.ResponseWriter, r *http.Request) {
		req, err := unmarshal(r)
		if err != nil {
			internalServerError(w, err)
			return
		}
		ok(w, d.s.SayHi(req))
	})
	return &http.Server{
		Addr:    addr,
		Handler: d,
	}
}
