package hello_test

import (
	"fmt"
	"testing"

	"github.com/anqur/gbio/examples/hello"
)

func TestClient(t *testing.T) {
	tx := hello.Do()
	r := tx.SayHi(&hello.SelfIntro{
		BaseGreeting: hello.BaseGreeting{ReqID: "lolwtf"},
		Name:         "Anqur",
	})
	if err := tx.Error; err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", r)
}
