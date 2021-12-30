package hello_test

import (
	"fmt"
	"testing"

	"github.com/anqur/gbio"
	etcd "go.etcd.io/etcd/client/v3"

	"github.com/anqur/gbio/examples/hello"
)

func TestClient(t *testing.T) {
	gbio.UseLookupRegistry(&etcd.Config{
		Endpoints:            nil,
		DialTimeout:          0,
		DialKeepAliveTime:    0,
		DialKeepAliveTimeout: 0,
		MaxCallSendMsgSize:   0,
		MaxCallRecvMsgSize:   0,
		TLS:                  nil,
		Username:             "",
		Password:             "",
		RejectOldCluster:     false,
		DialOptions:          nil,
		Context:              nil,
		Logger:               nil,
		LogConfig:            nil,
		PermitWithoutStream:  false,
	})
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
