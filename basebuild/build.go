package basebuild

import (
	_ "github.com/ottemo/commerce/env/config"   // System Configuration service
	_ "github.com/ottemo/commerce/env/cron"     // Schedule service
	_ "github.com/ottemo/commerce/env/errorbus" // Error Bus service
	_ "github.com/ottemo/commerce/env/eventbus" // Event Bus service
	_ "github.com/ottemo/commerce/env/ini"      // INI Configuration service
	_ "github.com/ottemo/commerce/env/logger"   // File-based Logging service

	_ "github.com/ottemo/commerce/api/context"   // Context runtime transfer service
	_ "github.com/ottemo/commerce/api/rest"      // RESTful API service
	_ "github.com/ottemo/commerce/api/session"   // Session Management service
	_ "github.com/ottemo/commerce/impex"         // Import/Export service
	_ "github.com/ottemo/commerce/media/fsmedia" // Media Storage service
	_ "github.com/ottemo/commerce/env/otto"      // Otto - JS like scripting language

	_ "github.com/ottemo/commerce/app/actors/category"        // Category module
	_ "github.com/ottemo/commerce/app/actors/cms"             // CMS Page/Block module
	_ "github.com/ottemo/commerce/app/actors/product"         // Product module
	_ "github.com/ottemo/commerce/app/actors/product/review"  // Product Reviews module
	_ "github.com/ottemo/commerce/app/actors/swatch"          // Product Reviews module
	_ "github.com/ottemo/commerce/app/actors/visitor"         // Visitor module
	_ "github.com/ottemo/commerce/app/actors/visitor/address" // Visitor Address module
	_ "github.com/ottemo/commerce/app/actors/visitor/token"   // Visitor Token module

	_ "github.com/ottemo/commerce/app/actors/cart"         // Shopping Cart module
	_ "github.com/ottemo/commerce/app/actors/checkout"     // Checkout module
	_ "github.com/ottemo/commerce/app/actors/order"        // Purchase Order module
	_ "github.com/ottemo/commerce/app/actors/stock"        // Stock Management module
	_ "github.com/ottemo/commerce/app/actors/subscription" // subscription extension
	_ "github.com/ottemo/commerce/app/actors/xdomain"      // XDomain support module

	_ "github.com/ottemo/commerce/app/actors/payment/authorizenet" // Authorize.Net payment method
	// _ "github.com/ottemo/commerce/app/actors/payment/braintree"    // Braintree payment method
	_ "github.com/ottemo/commerce/app/actors/payment/checkmo"      // "Check Money Order" payment method
	_ "github.com/ottemo/commerce/app/actors/payment/paypal"       // PayPal payment method
	// _ "github.com/ottemo/commerce/app/actors/payment/stripe"       // Stripe payment method

	_ "github.com/ottemo/commerce/app/actors/shipping/fedex"      // FedEx
	_ "github.com/ottemo/commerce/app/actors/shipping/flatrate"   // Flat Rate
	_ "github.com/ottemo/commerce/app/actors/shipping/flatweight" // Flat Weight
	_ "github.com/ottemo/commerce/app/actors/shipping/usps"       // USPS

	_ "github.com/ottemo/commerce/app/actors/discount/coupon"    // Coupon based discounts
	_ "github.com/ottemo/commerce/app/actors/discount/giftcard"  // Gift Cards
	_ "github.com/ottemo/commerce/app/actors/discount/saleprice" // Sale Price
	_ "github.com/ottemo/commerce/app/actors/tax"                // Tax Rates

	_ "github.com/ottemo/commerce/app/actors/reporting" // Reporting
	_ "github.com/ottemo/commerce/app/actors/rts"       // Real Time Statistics service
	_ "github.com/ottemo/commerce/app/actors/seo"       // URL Rewrite support

	_ "github.com/ottemo/commerce/app/actors/other/emma"         // Emma integration
	_ "github.com/ottemo/commerce/app/actors/other/friendmail"   // email friend extension
	_ "github.com/ottemo/commerce/app/actors/other/grouping"     // products grouping extension
	_ "github.com/ottemo/commerce/app/actors/other/mailchimp"    // MailChimp integration
	_ "github.com/ottemo/commerce/app/actors/other/shipstation"  // Shipstation integration
	_ "github.com/ottemo/commerce/app/actors/other/trustpilot"   // TrustPilot integration
	_ "github.com/ottemo/commerce/app/actors/other/vantagepoint" // VantagePoint integration

	_ "github.com/ottemo/commerce/app/actors/blog" // Blog module
)
