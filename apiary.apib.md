FORMAT: 1A
HOST: http://dev.ottemo.io:3000

# Ottemo Foundation API
Foundation is the api that powers the [Ottemo Storefront](http://www.ottemo.io/),
a ridiculously fast eCommerce solution.


# Group Products
app/actors/product/api.go

- [GET]      product/:productID
- [POST]     product
- [PUT]      product/:productID
- [DELETE]   product/:productID
- [GET]      product/:productID/media/:mediaType/:mediaName
- [GET]      product/:productID/media/:mediaType
- [POST]     product/:productID/media/:mediaType/:mediaName
- [DELETE]   product/:productID/media/:mediaType/:mediaName
- [GET]      product/:productID/mediapath/:mediaType
- [GET]      product/:productID/related
- [GET]      products
- [GET]      products/attributes
- [POST]     products/attribute
- [PUT]      products/attribute/:attribute
- [DELETE]   products/attribute/:attribute

app/actors/product/review/api.go

- [GET]      product/:productID/reviews
- [POST]     product/:productID/review
- [POST]     product/:productID/ratedreview/:stars
- [DELETE]   product/:productID/review/:reviewID
- [GET]      product/:productID/rating

# Group Categories
app/actors/category/api.go

- [GET]      categories
- [GET]      categories/tree
- [GET]      categories/attributes
- [POST]     category
- [PUT]      category/:categoryID
- [DELETE]   category/:categoryID
- [GET]      category/:categoryID
- [GET]      category/:categoryID/layers
- [GET]      category/:categoryID/products
- [POST]     category/:categoryID/product/:productID
- [DELETE]   category/:categoryID/product/:productID
- [GET]      category/:categoryID/media/:mediaType/:mediaName
- [GET]      category/:categoryID/media/:mediaType
- [POST]     category/:categoryID/media/:mediaType/:mediaName
- [DELETE]   category/:categoryID/media/:mediaType/:mediaName
- [GET]      category/:categoryID/mediapath/:mediaType

# Group CMS

## Blocks
app/actors/cms/block/api.go

- [GET]      cms/blocks
- [GET]      cms/block/:id
- [POST]     cms/block
- [PUT]      cms/block/:id
- [DELETE]   cms/block/:id
- [GET]      cms/blocks/attributes

## Pages
app/actors/cms/page/api.go

- [GET]      cms/pages
- [GET]      cms/page/:id
- [POST]     cms/page
- [PUT]      cms/page/:id
- [DELETE]   cms/page/:id
- [GET]      cms/pages/attributes

## Images
app/actors/cms/images/api.go

- [GET]      cms/images
- [POST]     cms/images
- [DELETE]   cms/images/:id

# Group Discounts
This refers to all capabilities to reduce the price of a product or cart at checkout.

## Coupons
app/actors/discount/coupon/api.go

- [GET]      coupons
- [POST]     coupons
- [GET]      csv/coupons
- [POST]     csv/coupons
- [POST]     cart/coupons
- [DELETE]   cart/coupons/:code
- [GET]      coupons/:id
- [PUT]      coupons/:id
- [DELETE]   coupons/:id

## Gift Cards
app/actors/discount/giftcard/api.go

- [GET]      giftcards
- [GET]      giftcards/:giftcode
- [GET]      giftcards/:giftcode/apply
- [GET]      giftcards/:giftcode/neglect

# Group Unsorted

## app/actors/payment/authorizenet/api.go

- [POST]     authorizenet/receipt
- [POST]     authorizenet/relay

## app/actors/payment/paypal/api.go

- [GET]      paypal/success
- [GET]      paypal/cancel

## app/actors/checkout/api.go

- [GET]      checkout
- [GET]      checkout/payment/methods
- [GET]      checkout/shipping/methods
- [PUT]      checkout/shipping/address
- [PUT]      checkout/billing/address
- [PUT]      checkout/payment/method/:method
- [PUT]      checkout/shipping/method/:method/:rate
- [PUT]      checkout
- [POST]     checkout/submit


## app/actors/order/api.go

- [GET]      orders/attributes
- [GET]      orders
- [GET]      order/:orderID
- [POST]     order
- [PUT]      order/:orderID
- [DELETE]   order/:orderID

## app/actors/tax/api.go

- [GET]      taxes/csv
- [POST]     taxes/csv


## app/actors/stock/api.go

- [GET]      stock/:productID
- [POST]     stock/:productID/:qty
- [PUT]      stock/:productID/:qty
- [DELETE]   stock/:productID
- [POST]     product/:productID/stock

## app/actors/seo/api.go

- [GET]      seo/items
- [POST]     seo/item
- [PUT]      seo/item/:itemID
- [DELETE]   seo/item/:itemID
- [GET]      seo/url/:url
- [GET]      seo/sitemap
- [GET]      seo/sitemap/sitemap.xml

## app/actors/rts/api.go

- [POST]     rts/visit
- [GET]      rts/visits

    ```
    {
      "error": null,
      "redirect": "",
      "result": {
        "total": {
          "today": 100,
          "yesterday": 200,
          "week": 900
        },
        "unique": {
          "today": 10,
          "yesterday": 20,
          "week": 90
        }
      }
    }
    ```

- [GET]      rts/visits/detail/:from/:to

    ```
    {
      "error": null,
      "redirect": "",
      "result": [
        [1431734400, 0],
        [1431820800, 0],
        [1431907200, 0],
        [1431993600, 0],
        [1432080000, 29],
        [1432166400, 60],
        [1432252800, 38]
      ]
    }
    ```

- [GET]      rts/sales

    ```
    {
      "error": null,
      "redirect": "",
      "result": {
        "sales": {
          "today": 5050.50,
          "yesterday": 20100.00,
          "week": 800300.00
        },
        "orders": {
          "today": 100,
          "yesterday": 200,
          "week": 900
        }
      }
    }
    ```

- [GET]      rts/sales/detail/:from/:to

    ```
    {
      "error": null,
      "redirect": "",
      "result": [
        [1431734400, 1],
        [1431820800, 12],
        [1431907200, 23],
        [1431993600, 34],
        [1432080000, 29],
        [1432166400, 60],
        [1432252800, 38]
      ]
    }
    ```

- [GET]      rts/bestsellers

    ```
    {
      "error": null,
      "redirect": "",
      "result": [
        {
          "count": 67,
          "image": "image/Product/5488485b49c43d4283000067/charge_slate_front.png",
          "name": "Charge",
          "pid": "5488485b49c43d4283000067"
        },
        {
          "count": 56,
          "image": "image/Product/5488485d49c43d428300006b/chargehr_plum_front.png",
          "name": "Charge HR",
          "pid": "5488485d49c43d428300006b"
        }
        // ... returns 10 products
      ]
    }
    ```

- [GET]      rts/referrers
- [GET]      rts/conversion
- [GET]      rts/visits/realtime

## app/actors/cart/api.go

- [GET]      cart
- [POST]     cart/item
- [PUT]      cart/item/:itemIdx/:qty
- [DELETE]   cart/item/:itemIdx

## app/actors/visitor/address/api.go

- [POST]     visitor/:visitorID/address
- [PUT]      visitor/:visitorID/address/:addressID
- [DELETE]   visitor/:visitorID/address/:addressID
- [GET]      visitor/:visitorID/addresses
- [GET]      visitors/addresses/attributes
- [DELETE]   visitors/address/:addressID
- [PUT]      visitors/address/:addressID
- [GET]      visitors/address/:addressID
- [POST]     visit/address
- [PUT]      visit/address/:addressID
- [DELETE]   visit/address/:addressID
- [GET]      visit/addresses
- [GET]      visit/address/:addressID

## app/actors/visitor/api.go

- [POST]     visitor
- [PUT]      visitor/:visitorID
- [DELETE]   visitor/:visitorID
- [GET]      visitor/:visitorID
- [GET]      visitors
- [GET]      visitors/attributes
- [DELETE]   visitors/attribute/:attribute
- [PUT]      visitors/attribute/:attribute
- [POST]     visitors/attribute
- [POST]     visitors/register
- [GET]      visitors/validate/:key
- [GET]      visitors/invalidate/:email
- [GET]      visitors/forgot-password/:email
- [POST]     visitors/mail
- [GET]      visit
- [PUT]      visit
- [GET]      visit/logout
- [POST]     visit/login
- [POST]     visit/login-facebook
- [POST]     visit/login-google
- [GET]      visit/orders
- [GET]      visit/order/:orderID

## app/api.go

- [GET]      app/login
- [POST]     app/login
- [GET]      app/logout
- [GET]      app/rights
- [GET]      app/status

## impex/api.go

- [GET]      impex/models
- [GET]      impex/import/status
- [GET]      impex/export/:model
- [POST]     impex/import/:model
- [POST]     impex/import
- [POST]     impex/test/import
- [POST]     impex/test/mapping

## env/config/api.go

- [GET]      config/groups
- [GET]      config/item/:path
- [GET]      config/values
- [GET]      config/values/refresh
- [GET]      config/value/:path
- [POST]     config/value/:path
- [PUT]      config/value/:path
- [DELETE]   config/value/:path

# Group Cron
Cron is a utility to schedule tasks.  These tasks maybe scheduled for
a specific time, they may be repeatable or intended to be run immediately.

It is important to understand several concepts.

<b>Task</b> - a job which can be scheduled to run at a specific time
<b>Schedule</b> - a listing of all active tasks, when they will be executed and their respective metadata

The API allows you to:
        * Obtain a list of the currently scheduled tasks
        * Create a task to be run on a schedule
        * Obtain a list of possible tasks to be scheduled
        * Enable a task to be run on a schedule
        * Disable a task
        * Update the specified task
        * Run the specified task now

## env/cron/api.go

- [GET]     cron/schedule
- [POST]    cron/task
- [GET]     cron/task
- [GET]     cron/task/enable/:taskIndex
- [GET]     cron/task/disable/:taskIndex
- [PUT]     cron/task/:taskIndex
- [GET]     cron/task/run/:taskIndex
