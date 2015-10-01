// Package friendmail is an extension which provides ability to send email to friend
package friendmail

import (
	"github.com/dchest/captcha"
	"github.com/ottemo/foundation/env"
	"sync"
	"time"
)

// Package global constants
const (
	ConstCollectionNameFriendMail = "friend_mail"

	ConstConfigPathFriendMail              = "general.friendmail"
	ConstConfigPathFriendMailEmailTemplate = "general.friendmail.template"
	ConstConfigPathFriendMailEmailSubject  = "general.friendmail.subject"

	ConstErrorModule = "friendmail"
	ConstErrorLevel  = env.ConstErrorLevelActor

	ConstMaxCaptchaItems = 100000 // captcha maximum amount (to prevent memory leaks)
	ConstCaptchaLifeTime = 300    // seconds generated captcha works (5 min)
)

var (
	captchaValuesMutex sync.RWMutex         // synchronization on captchaValues variable
	captchaValues      map[string]time.Time // global variable to track generated captcha codes
	captchaStore       captcha.Store
)
