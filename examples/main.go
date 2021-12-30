package main

import (
	"fmt"

	"github.com/anqur/gbio"

	"github.com/anqur/gbio/examples/hello"
)

type HelloService struct{}

func (HelloService) SayHi(g hello.Greeting) *hello.Reply {
	fmt.Printf("%+v\n", g)
	if i, ok := g.(*hello.SelfIntro); ok {
		return &hello.Reply{Message: fmt.Sprintf("Hi, %s!", i.Name)}
	}
	return &hello.Reply{Message: "Hi, stranger!"}
}

func main() {
	gbio.UseMux(hello.Mux(new(HelloService)))
	if err := gbio.ListenAndServe(); err != nil {
		panic(err)
	}
}
