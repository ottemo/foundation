package friendmail

import (
	"github.com/dchest/captcha"
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"time"
)

// init makes package self-initialization routine
func init() {
	db.RegisterOnDatabaseStart(setupDB)
	env.RegisterOnConfigStart(setupConfig)
	api.RegisterOnRestServiceStart(setupAPI)

	captchaValues = make(map[string]time.Time)

	captchaStore = captcha.NewMemoryStore(ConstMaxCaptchaItems, time.Second*ConstCaptchaLifeTime)
	captcha.SetCustomStore(captchaStore)

	// starting timer for garbage collector
	if ConstCaptchaLifeTime > 0 {
		timerInterval := time.Second * ConstCaptchaLifeTime
		ticker := time.NewTicker(timerInterval)
		go func() {
			for _ = range ticker.C {
				captchaValuesMutex.Lock()

				currentTime := time.Now()
				for key, value := range captchaValues {
					if currentTime.Sub(value).Seconds() >= ConstCaptchaLifeTime {
						captchaStore.Get(key, true)
						delete(captchaValues, key)
					}
				}

				captchaValuesMutex.Unlock()
			}
		}()
	}
}

// setupDB prepares system database for package usage
func setupDB() error {

	if collection, err := db.GetCollection(ConstCollectionNameFriendMail); err == nil {
		if err := collection.AddColumn("date", db.ConstTypeID, true); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "52acc483-0132-4252-acf3-9e965eefd73b", err.Error())
		}
		if err := collection.AddColumn("email", db.ConstTypeVarchar, true); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "53288f58-6cc8-4ae5-913d-d8f735a66678", err.Error())
		}
		if err := collection.AddColumn("data", db.ConstTypeJSON, false); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "c18bbf88-a104-4154-b8eb-0ccb6ae7a7e9", err.Error())
		}
	} else {
		return env.ErrorDispatch(err)
	}

	return nil
}
