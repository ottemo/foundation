package mailchimp

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/test"

	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestMailchimpSubscribe(tst *testing.T) {
	if err := test.StartAppInTestingMode(); err != nil {
		tst.Error(err)
	}

	//set the configuration to allow mailchimp
	var config = env.GetConfig()
	config.SetValue(ConstConfigPathMailchimpEnabled, true)
	config.SetValue(ConstConfigPathMailchimpAPIKey, "23dbf42618e8f43e624a6dd89de9bd46-us12")
	config.SetValue(ConstConfigPathMailchimpBaseURL, "https://us12.api.mailchimp.com/3.0/")

	rand.Seed(time.Now().UTC().UnixNano())
	testRegistration := Registration{
		EmailAddress: fmt.Sprintf("test+%d@myottemotest.com", rand.Int()),
		Status:       "subscribed",
		MergeFields: map[string]string{
			"FNAME": "Test",
			"LNAME": "User",
		},
	}

	if err := Subscribe("b9537d1e65", testRegistration); err != nil {
		tst.Error(err)
	}

}
