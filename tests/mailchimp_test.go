package tests

import (
	"fmt"
	"github.com/ottemo/foundation/app/actors/other/mailchimp"
	"math/rand"
	"testing"
	"time"
	"github.com/ottemo/foundation/env"
)

func TestMailchimpSubscribe(tst *testing.T) {
	if err := StartAppInTestingMode(); err != nil {
		tst.Error(err)
	}

	//set the configuration to allow mailchimp
	var config = env.GetConfig();
	config.SetValue(mailchimp.MailchimpEnabledConfig, true)
	config.SetValue(mailchimp.MailchimpApiKeyConfig,"23dbf42618e8f43e624a6dd89de9bd46-us12")
	config.SetValue(mailchimp.MailchimpBaseUrlConfig,"https://us12.api.mailchimp.com/3.0/")

	rand.Seed(time.Now().UTC().UnixNano())
	testRegistration := mailchimp.Registration{
		EmailAddress: fmt.Sprintf("test+%d@myottemotest.com", rand.Int()),
		Status:       "subscribed",
		MergeFields: map[string]string{
			"FNAME": "Test",
			"LNAME": "User",
		},
	}

	if err := mailchimp.Subscribe("b9537d1e65", testRegistration); err != nil {
		tst.Error(err)
	}

}
