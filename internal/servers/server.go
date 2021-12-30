package servers

import (
	"net/http"

	"github.com/anqur/gbio/internal/registries"
)

type Server struct {
	http.Server

	Reg *registries.Registry
}

type Service interface {
	ServiceName() string
}
