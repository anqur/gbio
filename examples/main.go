package main

import (
	"fmt"

	"github.com/anqur/gbio/examples/hello"
)

type HelloService struct{}

func (HelloService) SayHi(g hello.Greeting) *hello.Reply {
	if i, ok := g.(*hello.SelfIntro); ok {
		return &hello.Reply{Message: fmt.Sprintf("Hi, %s!", i.Name)}
	}
	return &hello.Reply{Message: "Hi, stranger!"}
}

func main() {
	if err := hello.
		NewServer(new(HelloService), ":8080").
		ListenAndServe(); err != nil {
		panic(err)
	}
}
