package codec

import "net/http"

type Encoder interface {
	Marshal() ([]byte, http.Header, error)
}
