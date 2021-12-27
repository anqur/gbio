package hello

type JustHi struct{}

type SelfIntro struct {
	Name string
}

type Greeting interface {
	isGreeting()
}

func (JustHi) isGreeting()    {}
func (SelfIntro) isGreeting() {}

type Reply struct {
	Message string
}

type Hello interface {
	SayHi(*Greeting) *Reply
}
