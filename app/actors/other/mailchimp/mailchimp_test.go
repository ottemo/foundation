package mailchimp_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/ottemo/foundation/app/actors/other/mailchimp"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/test"
	"github.com/ottemo/foundation/db"
)

func TestMailchimpSubscribe(t *testing.T) {
	// start app
	err := test.StartAppInTestingMode()
	if err != nil {
		t.Error(err)
	}

	db.RegisterOnDatabaseStart(func () error {
		testMailchimpSubscribe(t)
		return nil
	})
}

func testMailchimpSubscribe(tst *testing.T) {
	//set the configuration to allow mailchimp
	var config = env.GetConfig()
	if err := config.SetValue(mailchimp.ConstConfigPathMailchimpEnabled, true); err != nil {
		tst.Error(err)
	}
	if err := config.SetValue(mailchimp.ConstConfigPathMailchimpAPIKey, "23dbf42618e8f43e624a6dd89de9bd46-us12"); err != nil {
		tst.Error(err)
	}
	if err := config.SetValue(mailchimp.ConstConfigPathMailchimpBaseURL, "https://us12.api.mailchimp.com/3.0/"); err != nil {
		tst.Error(err)
	}

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
