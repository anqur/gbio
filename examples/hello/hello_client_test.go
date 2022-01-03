package hello_test

import (
	"testing"

	"github.com/anqur/gbio/examples/hello"
)

func TestSayHi(t *testing.T) {
	tx := hello.Do()
	r := tx.SayHi(&hello.SelfIntro{Name: "Anqur"})
	if err := tx.Error; err != nil {
		t.Fatal(err)
	}
	if r.Message != "Hi, Anqur!" {
		t.Fatal(r)
	}
}

func TestHiAdminOK(t *testing.T) {
	tx := hello.Do()
	r := tx.HiAdmin(&hello.ImAdmin{Authorization: "Bearer s3cr3t"})
	if err := tx.Error; err != nil {
		t.Fatal(err)
	}
	m, ok := r.(*hello.OkReply)
	if !ok {
		t.Fatal(r)
	}
	if m.Message != "Hi, admin!" {
		t.Fatal(m)
	}
}

func TestHiAdminFailed(t *testing.T) {
	tx := hello.Do()
	r := tx.HiAdmin(new(hello.ImAdmin))
	if err := tx.Error; err != nil {
		t.Fatal(err)
	}
	e, ok := r.(*hello.ErrReply)
	if !ok {
		t.Fatal(r)
	}
	if e.Code != hello.Unauthorized {
		t.Fatal(e)
	}
}
