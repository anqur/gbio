package hello

type SelfIntro struct {
	Name string
}

type ImAdmin struct {
	Authorization string `json:"-"`
}

type Code int

const (
	OK Code = iota
	InvalidParam
	Unauthorized
)

type Reply interface{ isReply() }
type OkReply struct{ Message string }
type ErrReply struct {
	Code  Code
	Error string
}

func (OkReply) isReply()  {}
func (ErrReply) isReply() {}

type Greeting interface {
	SayHi(*SelfIntro) *OkReply
}

type Admin interface {
	HiAdmin(*ImAdmin) Reply
}

type Hello interface {
	Greeting
	Admin
}
