package json

import (
	"github.com/julienschmidt/httprouter"
)

type DefaultService struct {
	ListenOn string
	Router   *httprouter.Router
	Handlers map[string]httprouter.Handle
}
