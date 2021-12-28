package hello

import (
	"fmt"
	"testing"
)

func TestClient(t *testing.T) {
	cl, err := NewClient("http://localhost:8080")
	if err != nil {
		t.Fatal(err)
	}

	call := cl.Call()
	r := call.SayHi(&SelfIntro{
		BaseGreeting: BaseGreeting{ReqID: "lolwtf"},
		Name:         "Anqur",
	})
	if err := call.Error; err != nil {
		t.Fatal(err)
	}
	fmt.Println(r)
}
