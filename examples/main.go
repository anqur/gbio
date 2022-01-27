package main

import (
	"fmt"
	"strings"

	"github.com/anqur/gbio"

	"github.com/anqur/gbio/examples/hello"
)

type Greeting struct{}

func (Greeting) SayHi(i *hello.SelfIntro) *hello.OkReply {
	return &hello.OkReply{Message: fmt.Sprintf("Hi, %s!", i.Name)}
}

type GreetingV2 struct{}

func (GreetingV2) SayHi(i *hello.SelfIntro) *hello.OkReply {
	return &hello.OkReply{Message: fmt.Sprintf("Aloha, %s!", i.Name)}
}

type Admin struct{}

func (Admin) HiAdmin(i *hello.ImAdmin) hello.Reply {
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
	gbio.UseServer(
		hello.RegisterGreeting(new(Greeting)),
		hello.RegisterGreeting(new(GreetingV2), gbio.WithTag("v2")),
		hello.RegisterAdmin(new(Admin)),
	)
	if err := gbio.ListenAndServe(); err != nil {
		panic(err)
	}
}
