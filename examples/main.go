package main

import (
	"fmt"
	"strings"

	"github.com/anqur/gbio"

	"github.com/anqur/gbio/examples/hello"
)

type HelloService struct{}

func (HelloService) SayHi(i *hello.SelfIntro) *hello.OkReply {
	return &hello.OkReply{Message: fmt.Sprintf("Hi, %s!", i.Name)}
}

func (HelloService) HiAdmin(i *hello.ImAdmin) hello.Reply {
	tk := strings.TrimPrefix(i.Authorization, "Bearer ")
	if tk != "s3cr3t" {
		return &hello.ErrReply{
			Code:  hello.Unauthorized,
			Error: "nah you're not admin",
		}
	}
	return &hello.OkReply{Message: "Hi, admin!"}
}

func main() {
	gbio.UseMux(hello.Mux(new(HelloService)))
	if err := gbio.ListenAndServe(); err != nil {
		panic(err)
	}
}
