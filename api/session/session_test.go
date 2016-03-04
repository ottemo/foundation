package session

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/utils"
	"math/rand"
	"testing"
)

// TestSessionsConcurrency tests synchronisation mechanisms between go routines
func TestSessionsConcurrency(t *testing.T) {
	var sessions []api.InterfaceSession

	const sessionsNumber = 100
	const routinesNumber = 1000

	sync := make(chan bool)

	// preparing set of sessions to work with
	for i := 0; i < sessionsNumber; i++ {
		session, err := api.NewSession()
		if err != nil {
			t.Error(err)
		}
		sessions = append(sessions, session)
	}

	// making go-routines storm
	for i := 0; i < routinesNumber; i++ {
		go func() {
			session := sessions[rand.Intn(sessionsNumber)]
			for i := 0; i < 100; i++ {
				key := utils.InterfaceToString(i)
				session.Set(key, utils.InterfaceToInt(session.Get(key))+1)
			}

			SessionService.GC()

			sync <- true
		}()
	}

	// waiting till all routines finishes their stuff
	finished := 0
	for finished = 0; finished < routinesNumber; <-sync {
		finished++
	}

	// closing all the sessions
	for i := 0; i < sessionsNumber; i++ {
		sessions[i].Close()
	}
}

// BenchmarkIoOperations evaluates the performance of Set/Get operations
func BenchmarkIoOperations(b *testing.B) {
	session, err := api.NewSession()
	if err != nil {
		b.Fail()
	}

	for i := 0; i < 100; i++ {
		session.Set("test", i)

		if result := session.Get("test"); result != i {
			b.Error("assigned value not matches:", result, "!=", i)
		}
	}
}
