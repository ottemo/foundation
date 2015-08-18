// Package session is a default implementation of InterfaceSessionService and InterfaceSession
// declared in "github.com/ottemo/foundation/api" package
package session

import (
	"sync"
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstSessionLifeTime          = 26280000 // session idle period before expire (in sec); set to 365 days
	ConstSessionUpdateTime        = 10       // '0' - immediate mode, '>0' - update timer mode (in sec)
	ConstSessionKeepInMemoryItems = 1000     // limits application sessions array for "immediate mode", '0' - unlimited

	ConstErrorModule = "api/session"
	ConstErrorLevel  = env.ConstErrorLevelService

	ConstStorageFolder = "./var/session/"
	ConstCryptSession  = false
)

// Package global variables
var (
	SessionService api.InterfaceSessionService
)

// DefaultSessionService is a basic implementer of InterfaceSessionService declared in
// "github.com/ottemo/foundation/api" package
type DefaultSessionService struct {
	Sessions      map[string]*DefaultSessionContainer // active sessions set
	sessionsMutex sync.RWMutex                        // syncronization on Sessions variable modification

	// package supports "memcache", "redis", "memsession" build tags to change default (filesystem) storage location
	Storage InterfaceServiceStorage
}

// DefaultSession is a default implementer of InterfaceSession declared in
// "github.com/ottemo/foundation/api" package
type DefaultSession string

// DefaultSessionContainer is a structure to hold session related information
type DefaultSessionContainer struct {
	id        DefaultSession
	Data      map[string]interface{}
	UpdatedAt time.Time
}

// InterfaceServiceStorage session storage layer for a session service
type InterfaceServiceStorage interface {
	GetStorageName() string

	LoadSession(sessionID string) (*DefaultSessionContainer, error)
	FlushSession(sessionID string) error
}
