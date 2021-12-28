package hello

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

type helloDelegate struct {
	http.ServeMux

	s Hello
}

type discriminator struct {
	Tag string `json:"_t"`
}

var emptyStrVal = reflect.ValueOf("")

func parseCtxValue(ty reflect.Type, v []string) (reflect.Value, error) {
	switch ty.Kind() {
	case reflect.String:
		if len(v) > 0 {
			return reflect.ValueOf(v[0]), nil
		}
		return emptyStrVal, nil
	case reflect.Slice:
		if ty.Elem().Kind() != reflect.String {
			return reflect.Value{}, fmt.Errorf("expected `[]string` for context value slice, found %q", ty.Kind())
		}
		return reflect.ValueOf(v), nil
	}
	return reflect.Value{}, fmt.Errorf("expected `string` or `[]string` as context value type, found %q", ty.Kind())
}

func unmarshalCtx(r *http.Request, req interface{}) error {
	ty := reflect.TypeOf(req)
	val := reflect.ValueOf(req)
	if ty.Kind() == reflect.Ptr {
		ty = ty.Elem()
		val = val.Elem()
	}
	if ty.Kind() != reflect.Struct {
		return fmt.Errorf("expected struct model, found %q", ty.Kind())
	}

	for i := 0; i < ty.NumField(); i++ {
		tyField := ty.Field(i)
		valField := val.Field(i)
		if tyField.Anonymous {
			anon := valField.Interface()
			if tyField.Type.Kind() == reflect.Struct {
				anon = valField.Addr().Interface()
			}
			if err := unmarshalCtx(r, anon); err != nil {
				return err
			}
			continue
		}

		ctxKey := tyField.Tag.Get("ctx")
		if ctxKey == "" {
			continue
		}

		k := http.CanonicalHeaderKey(ctxKey)
		v, ok := r.Header[k]
		if !ok {
			continue
		}

		val, err := parseCtxValue(tyField.Type, v)
		if err != nil {
			return err
		}

		valField.Set(val)
	}

	return nil
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

	if err := unmarshalCtx(r, req); err != nil {
		return nil, err
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
