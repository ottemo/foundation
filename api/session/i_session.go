package session

import (
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// GetID returns current session id
func (it *DefaultSession) GetID() string {
	if it == nil {
		env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0467d401-030a-4ca4-9540-a4e2e72cc736", "nil session instance")
		return ""
	}

	return it.id
}

// Get returns session value by a given key or nil - if not set
func (it *DefaultSession) Get(key string) interface{} {
	// checking current instance
	if it == nil {
		env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0467d401-030a-4ca4-9540-a4e2e72cc736", "nil session instance")
		return nil
	}

	// immediate mode
	if ConstSessionUpdateTime == 0 {
		err := sessionService.updateSession(it.id, it)
		if err != nil {
			env.ErrorDispatch(err)
		}
	}

	// requested session operation
	if value, ok := it.Data[key]; ok == true {
		return value
	}
	return nil
}

// Set assigns value to session key
func (it *DefaultSession) Set(key string, value interface{}) {
	// checking current instance
	if it == nil {
		env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0467d401-030a-4ca4-9540-a4e2e72cc736", "nil session instance")
	}

	// immediate mode
	if ConstSessionUpdateTime == 0 {
		err := sessionService.updateSession(it.id, it)
		if err != nil {
			env.ErrorDispatch(err)
		}
	}

	// requested session operation
	it.Data[key] = value

	it.UpdatedAt = time.Now()

	// checking case session was already flushed
	if sessionInstance, present := sessionService.Sessions[it.GetID()]; !present || sessionInstance != it {
		sessionService.Sessions[it.GetID()] = it
	}

	// immediate mode
	if ConstSessionUpdateTime == 0 {
		err := it.Save()
		if err != nil {
			env.ErrorDispatch(err)
		}
	}
}

// SetModified updates last modification time to current moment
func (it *DefaultSession) SetModified() {
	it.UpdatedAt = time.Now()

	// checking case session was already flushed
	if sessionInstance, present := sessionService.Sessions[it.GetID()]; !present || sessionInstance != it {
		sessionService.Sessions[it.GetID()] = it
	}
}

// Close clears session data
func (it *DefaultSession) Close() error {
	// checking current instance
	if it == nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0467d401-030a-4ca4-9540-a4e2e72cc736", "nil session instance")
	}

	sessionID := it.GetID()

	// making system event
	eventData := map[string]interface{}{"session": it, "sessionID": sessionID}
	env.Event("session.close", eventData)

	// removing cache references
	if _, present := sessionService.Sessions[sessionID]; present {
		sessionService.sessionsMutex.Lock()
		delete(sessionService.Sessions, sessionID)
		sessionService.sessionsMutex.Unlock()
	}

	// removing storage references
	fileName := ConstStorageFolder + sessionID
	if _, err := os.Stat(fileName); os.IsExist(err) {
		err := os.Remove(fileName)
		if err != nil {
			env.ErrorDispatch(err)
		}
	}

	return nil
}

// Save stores session for a long period
func (it *DefaultSession) Save() error {
	// checking current instance
	if it == nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0467d401-030a-4ca4-9540-a4e2e72cc736", "nil session instance")
	}

	// checking session id
	sessionID := it.GetID()
	if sessionID == "" {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "af9b84a1-8c82-443c-919f-8dd75f956887", "session id is blank")
	}

	// saving session data
	sessionFile, err := os.OpenFile(ConstStorageFolder+sessionID, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)
	defer func() {
		sessionFile.Close()
		os.Chtimes(sessionFile.Name(), it.UpdatedAt, it.UpdatedAt)
	}()

	if err != nil {
		return env.ErrorDispatch(err)
	}

	var writer io.Writer = sessionFile
	if ConstCryptSession {
		writer, err = utils.EncryptWriter(writer)
		if err != nil {
			return env.ErrorDispatch(err)
		}
	}

	jsonEncoder := json.NewEncoder(writer)
	err = jsonEncoder.Encode(it)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// Load restores session from a long period storage
func (it *DefaultSession) Load(sessionID string) error {

	// checking session id
	if sessionID == "" {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "af9b84a1-8c82-443c-919f-8dd75f956887", "session id is blank")
	}

	// checking current instance
	if it == nil {
		it = new(DefaultSession)
	}
	it.id = sessionID

	// loading session data to instance
	filename := ConstStorageFolder + sessionID
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "363cd5a8-1a3d-4163-a7d3-cb96dbaff01c", "session "+sessionID+" not found")
	}

	// checking for expired session
	if time.Now().Sub(fileInfo.ModTime()).Seconds() >= ConstSessionLifeTime {
		err := os.Remove(filename)
		if err != nil {
			return err
		}
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "7aee9352-a08b-420a-a725-7f32a17495a8", "session "+sessionID+" expired")
	}

	// if not expired
	sessionFile, err := os.OpenFile(filename, os.O_RDONLY, 0660)
	defer sessionFile.Close()

	var reader io.Reader = sessionFile
	if ConstCryptSession {
		reader, err = utils.EncryptReader(reader)
		if err != nil {
			return env.ErrorDispatch(err)
		}
	}

	jsonDecoder := json.NewDecoder(reader)
	err = jsonDecoder.Decode(it)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	it.UpdatedAt = fileInfo.ModTime()

	// updating loaded if needed
	if it.Data == nil {
		it.Data = make(map[string]interface{})
	}

	if it.UpdatedAt.IsZero() {
		it.UpdatedAt = time.Now()
	}

	// updating sessions cache
	sessionService.sessionsMutex.Lock()
	sessionService.Sessions[sessionID] = it
	sessionService.sessionsMutex.Unlock()

	return nil
}
