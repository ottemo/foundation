package shipstation

import "github.com/ottemo/foundation/env"

// buildTrackingUrl Assemble a tracking url from the carrier and tracking number
// - http://verysimple.com/2011/07/06/ups-tracking-url/
// - http://eshipguy.com/tracking/
func buildTrackingUrl(carrier string, trackingNumber string) string {
	var trackingUrl string

	// These are the recognized carriers that shipstation might pass back
	trackingMap := map[string]string{
		"AccessWorldwide":               "",
		"APC":                           "",
		"Asendia":                       "",
		"AustraliaPost":                 "",
		"BrokersWorldWide":              "",
		"CanadaPost":                    "",
		"DHL":                           "http://track.dhl-usa.com/TrackByNbr.asp?ShipmentNumber=",
		"DHLCanada":                     "",
		"DHLGlobalMail":                 "",
		"FedEx":                         "http://www.fedex.com/Tracking?action=track&tracknumbers=",
		"FedExCanada":                   "",
		"FedExInternationalMailService": "",
		"FirstMile":                     "",
		"Globegistics":                  "",
		"IMEX":                          "",
		"LoneStar":                      "",
		"Newgistics":                    "",
		"OnTrac":                        "",
		"Other":                         "",
		"RoyalMail":                     "",
		"UPS":                           "http://wwwapps.ups.com/WebTracking/track?track=yes&trackNums=",
		"UPSMI":                         "https://www.ups-mi.net/packageID/packageid.aspx?pid=",
		"USPS":                          "https://tools.usps.com/go/TrackConfirmAction_input?qtc_tLabels1=",
	}

	trackingUrl = trackingMap[carrier]
	if trackingUrl == "" {
		env.LogError(env.ErrorNew(ConstErrorModule, ConstErrorLevel, "1799137F-0FAD-4200-BAFF-954C30FCE674", "We don't have a tracking url set up for a certain carrier: "+carrier))
	} else {
		trackingUrl += trackingNumber
	}

	return trackingUrl
}
