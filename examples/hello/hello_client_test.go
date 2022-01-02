package hello_test

import (
	"fmt"
	"testing"

	"github.com/anqur/gbio/examples/hello"
)

func TestSayHi(t *testing.T) {
	tx := hello.Do()
	r := tx.SayHi(&hello.SelfIntro{Name: "Anqur"})
	if err := tx.Error; err != nil {
		t.Fatal(err)
	}
	fmt.Println(r.Message)
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
	fmt.Println(m.Message)
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
	fmt.Println(e.Code, e.Error)
}
