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

type Reply interface{ isReply() int }
type OkReply struct{ Message string }
type ErrReply struct {
	Code  Code
	Error string
}

func (OkReply) isReply() int  { return 1 }
func (ErrReply) isReply() int { return 2 }

type Greeting interface {
	SayHi(*SelfIntro) *OkReply
}

type Admin interface {
	HiAdmin(*ImAdmin) Reply
}
