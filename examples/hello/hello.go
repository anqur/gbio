package hello

type BaseGreeting struct {
	ReqID string `json:"-"`
}

type JustHi struct {
	BaseGreeting
}

type SelfIntro struct {
	BaseGreeting

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
	SayHi(Greeting) *Reply
}
