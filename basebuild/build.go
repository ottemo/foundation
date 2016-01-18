package basebuild

import (
	_ "github.com/ottemo/foundation/env/config"   // System Configuration service
	_ "github.com/ottemo/foundation/env/cron"     // Schedule service
	_ "github.com/ottemo/foundation/env/errorbus" // Error Bus service
	_ "github.com/ottemo/foundation/env/eventbus" // Event Bus service
	_ "github.com/ottemo/foundation/env/ini"      // INI Configuration service
	_ "github.com/ottemo/foundation/env/logger"   // File-based Logging service

	_ "github.com/ottemo/foundation/api/rest"      // RESTful API service
	_ "github.com/ottemo/foundation/api/session"   // Session Management service
	_ "github.com/ottemo/foundation/impex"         // Import/Export service
	_ "github.com/ottemo/foundation/media/fsmedia" // Media Storage service

	_ "github.com/ottemo/foundation/app/actors/category"        // Category module
	_ "github.com/ottemo/foundation/app/actors/cms"             // CMS Page/Block module
	_ "github.com/ottemo/foundation/app/actors/product"         // Product module
	_ "github.com/ottemo/foundation/app/actors/product/review"  // Product Reviews module
	_ "github.com/ottemo/foundation/app/actors/visitor"         // Visitor module
	_ "github.com/ottemo/foundation/app/actors/visitor/address" // Visitor Address module
	_ "github.com/ottemo/foundation/app/actors/visitor/token"   // Visitor Token module

	_ "github.com/ottemo/foundation/app/actors/cart"         // Shopping Cart module
	_ "github.com/ottemo/foundation/app/actors/checkout"     // Checkout module
	_ "github.com/ottemo/foundation/app/actors/order"        // Purchase Order module
	_ "github.com/ottemo/foundation/app/actors/stock"        // Stock Management module
	_ "github.com/ottemo/foundation/app/actors/subscription" // subscription extension

	_ "github.com/ottemo/foundation/app/actors/payment/authorizenet" // Authorize.Net payment method
	_ "github.com/ottemo/foundation/app/actors/payment/checkmo"      // "Check Money Order" payment method
	_ "github.com/ottemo/foundation/app/actors/payment/paypal"       // PayPal payment method

	_ "github.com/ottemo/foundation/app/actors/shipping/fedex"    // FedEx shipping method
	_ "github.com/ottemo/foundation/app/actors/shipping/flatrate" // "Flat Rate" shipping method
	_ "github.com/ottemo/foundation/app/actors/shipping/usps"     // USPS shipping method

	_ "github.com/ottemo/foundation/app/actors/discount/coupon"   // Coupon based discounts
	_ "github.com/ottemo/foundation/app/actors/discount/giftcard" // Gift Cards
	_ "github.com/ottemo/foundation/app/actors/tax"               // Tax Rates

	_ "github.com/ottemo/foundation/app/actors/rts" // Real Time Statistics service
	_ "github.com/ottemo/foundation/app/actors/seo" // URL Rewrite support

	_ "github.com/ottemo/foundation/app/actors/other/friendmail" // email friend extension
	_ "github.com/ottemo/foundation/app/actors/other/grouping"   // products grouping extension
	_ "github.com/ottemo/foundation/app/actors/other/mailchimp"  // MailChimp integration
	_ "github.com/ottemo/foundation/app/actors/other/quickbooks" // QuickBooks exporting extension
	_ "github.com/ottemo/foundation/app/actors/other/trustpilot" // TrustPilot integration
)
