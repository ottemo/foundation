package basebuild

import (
	_ "github.com/ottemo/foundation/env/config"   // system configuration service
	_ "github.com/ottemo/foundation/env/errorbus" // error bus service
	_ "github.com/ottemo/foundation/env/eventbus" // event bus service
	_ "github.com/ottemo/foundation/env/ini"      // ini configuration service
	_ "github.com/ottemo/foundation/env/logger"   // file-based logger service

	_ "github.com/ottemo/foundation/api/rest"      // RESTful API service
	_ "github.com/ottemo/foundation/api/session"   // API session management service
	_ "github.com/ottemo/foundation/impex"         // import/export service
	_ "github.com/ottemo/foundation/media/fsmedia" // storage manager service

	_ "github.com/ottemo/foundation/app/actors/category"        // category model implementation
	_ "github.com/ottemo/foundation/app/actors/cms"             // cms page/block models implementation
	_ "github.com/ottemo/foundation/app/actors/product"         // product model implementation
	_ "github.com/ottemo/foundation/app/actors/product/review"  // product reviews support
	_ "github.com/ottemo/foundation/app/actors/visitor"         // visitor model implementation
	_ "github.com/ottemo/foundation/app/actors/visitor/address" // visitor address model implementation

	_ "github.com/ottemo/foundation/app/actors/cart"     // checkout cart model implementation
	_ "github.com/ottemo/foundation/app/actors/checkout" // checkout model implementation
	_ "github.com/ottemo/foundation/app/actors/order"    // purchase order model implementation
	_ "github.com/ottemo/foundation/app/actors/stock"    // stock management model implementation

	_ "github.com/ottemo/foundation/app/actors/payment/authorizenet" // Authorize.Net payment method
	_ "github.com/ottemo/foundation/app/actors/payment/checkmo"      // "Check Money Order" payment method
	_ "github.com/ottemo/foundation/app/actors/payment/paypal"       // PayPal payment method

	_ "github.com/ottemo/foundation/app/actors/shipping/fedex"    // FedEx shipping method
	_ "github.com/ottemo/foundation/app/actors/shipping/flatrate" // "Flat Rate" shipping method
	_ "github.com/ottemo/foundation/app/actors/shipping/usps"     // USPS shipping method

	_ "github.com/ottemo/foundation/app/actors/discount" // coupon based discounts
	_ "github.com/ottemo/foundation/app/actors/tax"      // shipping tax rates

	_ "github.com/ottemo/foundation/app/actors/rts" // real time statictic service
	_ "github.com/ottemo/foundation/app/actors/seo" // url rewrites support
)
