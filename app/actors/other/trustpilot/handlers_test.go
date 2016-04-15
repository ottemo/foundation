package trustpilot

import (
	"testing"
)

func getGoodCredentials() (tpCredentials, bool) {
	isValidCredentials := false
	goodCredentials := tpCredentials{
		username:  "",
		password:  "",
		apiKey:    "",
		apiSecret: "",
	}
	return goodCredentials, isValidCredentials
}

func TestNotAuthorized(t *testing.T) {
	emptyCredentials := tpCredentials{}
	badCredentials := tpCredentials{
		username:  "username",
		password:  "password",
		apiKey:    "key",
		apiSecret: "secret",
	}

	tests := []tpCredentials{
		emptyCredentials,
		badCredentials,
	}

	expectedErrMsg := "Non 200 response while trying to get trustpilot access token: StatusCode:401 Unauthorized"
	for _, cred := range tests {
		token, err := getAccessToken(cred)

		if token != "" || err.Error() != expectedErrMsg {
			t.Error("expected empty token, and a 403 response status")
		}
	}
}

func TestGoodCredentials(t *testing.T) {
	cred, haveCreds := getGoodCredentials()
	if !haveCreds {
		return
	}

	token, err := getAccessToken(cred)
	if token == "" || err != nil {
		t.Error("expected a valid token and no error", err)
	}
}

func TestProductReviewLink(t *testing.T) {
	businessID := "ABC" // TODO: WHAT DO I DO FOR TESTING HERE
	productReviewData := ProductReview{
		Consumer: ProductReviewConsumer{
			Email: "test@ottemo.io",
			Name:  "Jon Dough",
		},
		Products: []ProductReviewProduct{
			ProductReviewProduct{
				ProductURL: "",
				ImageURL:   "",
				Name:       "",
				Sku:        "",
				Brand:      "",
			},
		},
		ReferenceID: "abc123",
		Locale:      requestLocale,
	}

	// Testing bad credentials
	badToken := "badToken"
	productReviewLink, err := getProductReviewLink(productReviewData, businessID, badToken)
	errMsg := "Non 200 response while trying to get trustpilot review link: StatusCode:401 Unauthorized"
	if productReviewLink != "" || err.Error() != errMsg {
		t.Error("expected an authorization error when testing product review link with a bad token, businessID")
	}

	// Testing good credentials
	// cred, haveCreds := getGoodCredentials()
	// if !haveCreds {
	// 	return
	// }
	// token, _ = getAccessToken(cred)
	// productReviewLink, err := getProductReviewLink(productReviewData, businessID, accessToken)

}
