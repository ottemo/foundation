package trustpilot

import (
	"fmt"
	"testing"
)

func getGoodCredentials() (tpCredentials, bool) {
	isValidCredentials := true
	goodCredentials := tpCredentials{
		username:  "engineering@ottemo.io",
		password:  "Ottemo2016!",
		apiKey:    "8N3P9QOvhrC4ej2PawPRoRD6O4BjjsWM",
		apiSecret: "LeiZXj9t7kf16dDG",
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
	businessID := "54f078980000ff00057db649" // TODO: WHAT DO I DO FOR TESTING HERE
	productReviewData := ProductReview{
		Consumer: ProductReviewConsumer{
			Email: "adam+tp-test@ottemo.io",
			Name:  "Jon Dough",
		},
		Products: []ProductReviewProduct{
			ProductReviewProduct{
				ProductURL: "https://karigran.com/facial-oil/the-system",
				ImageURL:   "https://karigran.com/media/image/Product/5511ff1cd4a2560a1400000e/Kari-Gran-Facial-Oil-The-KG-System-Product-Desktop_1448658076_515x515.png",
				Name:       "The KG System",
				Sku:        "KGSYS",
				Brand:      "Karigran",
			},
		},
		ReferenceID: "abc123",
		Locale:      requestLocale,
	}

	// Testing bad credentials
	// badToken := "badToken"
	// productReviewLink, err := getProductReviewLink(productReviewData, businessID, badToken)
	// errMsg := "Non 200 response while trying to get trustpilot review link: StatusCode:401 Unauthorized"
	// if productReviewLink != "" || err.Error() != errMsg {
	// 	t.Error("expected an authorization error when testing product review link with a bad token, businessID")
	// }

	// Testing good credentials
	cred, haveCreds := getGoodCredentials()
	if !haveCreds {
		return
	}
	accessToken, _ := getAccessToken(cred)
	productReviewLink, err := getProductReviewLink(productReviewData, businessID, accessToken)
	fmt.Println(productReviewLink, err)
	if err != nil {
		t.Error("expected a product review link back: ")
	}
}
