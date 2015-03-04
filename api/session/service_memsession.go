// +build memsession

// "service_memsession.go" is a blank session storage - "memsession" build tag should be specified in order to use it
// (session instances holds only on memory without flushing to longer term storage)

package session

import (
	"github.com/ottemo/foundation/api"
)

// init makes package self-initialization routine
func init() {
	SessionService = InitDefaultSessionService()

	// service registration within system
	api.RegisterSessionService(SessionService)
}
