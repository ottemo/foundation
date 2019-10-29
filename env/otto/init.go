package otto

import (
	"fmt"
	"github.com/ottemo/commerce/api"
	"github.com/ottemo/commerce/app"
	"github.com/ottemo/commerce/app/models"
	"github.com/ottemo/commerce/app/models/cart"
	"github.com/ottemo/commerce/app/models/category"
	"github.com/ottemo/commerce/app/models/checkout"
	"github.com/ottemo/commerce/app/models/cms"
	"github.com/ottemo/commerce/app/models/order"
	"github.com/ottemo/commerce/app/models/product"
	"github.com/ottemo/commerce/app/models/seo"
	"github.com/ottemo/commerce/app/models/stock"
	"github.com/ottemo/commerce/app/models/subscription"
	"github.com/ottemo/commerce/app/models/visitor"
	"github.com/ottemo/commerce/db"
	"github.com/ottemo/commerce/media"

	"github.com/ottemo/commerce/env"
	"github.com/ottemo/commerce/impex"
	"github.com/ottemo/commerce/utils"
)

// init performs package self-initialization
func init() {
	api.RegisterOnRestServiceStart(setupAPI)

	engine = new(ScriptEngine)
	engine.instances = make(map[string]*Script)
	engine.mappings = make(map[string]interface{})

	engine.Set("printf", fmt.Sprintf)
	engine.Set("print", fmt.Sprint)

	engine.Set("api.makeContext", makeApplicationContext)
	engine.Set("api.getHandler", apiHandler)
	engine.Set("api.call", apiCall)

	engine.Set("app.getVersion", app.GetVerboseVersion)
	engine.Set("app.getStorefrontURL", app.GetStorefrontURL)
	engine.Set("app.getDashboardURL", app.GetDashboardURL)
	engine.Set("app.getCommerceURL", app.GetcommerceURL)

	engine.Set("app.sendMail", app.SendMail)
	engine.Set("app.sendMailEx", app.SendMailEx)

	engine.Set("model.get", models.LoadModelByID)
	engine.Set("model.getModel", models.GetModel)
	engine.Set("model.getModels", models.GetDeclaredModels)

	engine.Set("model.test", func(x interface{ GetName() string }) string { return x.GetName() })

	engine.Set("model.product.get", product.LoadProductByID)
	engine.Set("model.product.getModel", product.GetProductModel)
	engine.Set("model.product.getCollection", product.GetProductCollectionModel)
	engine.Set("model.product.getStock", product.GetRegisteredStock())

	engine.Set("model.stock.getModel", stock.GetStockModel)
	engine.Set("model.stock.getCollection", stock.GetStockCollectionModel)

	engine.Set("model.visitor.get", visitor.LoadVisitorByID)
	engine.Set("model.visitor.getModel", visitor.GetVisitorModel)
	engine.Set("model.visitor.getCollection", visitor.GetVisitorCollectionModel)
	engine.Set("model.visitor.getAddressModel", visitor.GetVisitorAddressModel)
	engine.Set("model.visitor.getCardModel", visitor.GetVisitorCardModel)

	engine.Set("model.subscription.get", subscription.LoadSubscriptionByID)
	engine.Set("model.subscription.getModel", subscription.GetSubscriptionModel)
	engine.Set("model.subscription.getCollection", subscription.GetSubscriptionCollectionModel)
	engine.Set("model.subscription.getOptionValues", subscription.GetSubscriptionOptionValues)
	engine.Set("model.subscription.getCronExpr", subscription.GetSubscriptionCronExpr)
	engine.Set("model.subscription.getPeriodValue", subscription.GetSubscriptionPeriodValue)
	engine.Set("model.subscription.isEnabled", subscription.IsSubscriptionEnabled)

	engine.Set("model.seo.get", seo.LoadSEOItemByID)
	engine.Set("model.seo.getModel", seo.GetSEOItemModel)
	engine.Set("model.seo.getCollection", seo.GetSEOItemCollectionModel)
	engine.Set("model.seo.getEngine", seo.GetRegisteredSEOEngine)

	engine.Set("model.order.get", order.LoadOrderByID)
	engine.Set("model.order.getModel", order.GetOrderModel)
	engine.Set("model.order.getCollection", order.GetOrderCollectionModel)
	engine.Set("model.order.getItemCollection", order.GetOrderItemCollectionModel)
	engine.Set("model.order.getItemsForOrders", order.GetItemsForOrders)
	engine.Set("model.order.getOrdersCreatedBetween", order.GetOrdersCreatedBetween)
	engine.Set("model.order.getOrdersUpdatedBetween", order.GetFullOrdersUpdatedBetween)

	engine.Set("model.cms.getBlockById", cms.LoadCMSBlockByID)
	engine.Set("model.cms.getBlockByIdentifier", cms.LoadCMSBlockByIdentifier)
	engine.Set("model.cms.getBlockModel", cms.GetCMSBlockModel)
	engine.Set("model.cms.getBlockCollection", cms.GetCMSBlockCollectionModel)
	engine.Set("model.cms.getPageById", cms.LoadCMSPageByID)
	engine.Set("model.cms.getPageByIdentifier", cms.LoadCMSPageByIdentifier)
	engine.Set("model.cms.getPageModel", cms.GetCMSPageModel)
	engine.Set("model.cms.getPageCollection", cms.GetCMSPageCollectionModel)

	engine.Set("model.category.get", category.LoadCategoryByID)
	engine.Set("model.category.getModel", category.GetCategoryModel)
	engine.Set("model.category.getCollection", category.GetCategoryCollectionModel)

	engine.Set("model.checkout.getModel", checkout.GetCheckoutModel)
	engine.Set("model.checkout.getPaymentMethodByCode", checkout.GetPaymentMethodByCode)
	engine.Set("model.checkout.getShippingMethodByCode", checkout.GetShippingMethodByCode)
	engine.Set("model.checkout.getPaymentMethods", checkout.GetRegisteredPaymentMethods)
	engine.Set("model.checkout.getShippingMethods", checkout.GetRegisteredShippingMethods)
	engine.Set("model.checkout.validateAddress", checkout.ValidateAddress)

	engine.Set("model.cart.get", cart.LoadCartByID)
	engine.Set("model.cart.getModel", cart.GetCartModel)
	engine.Set("model.cart.getCartForVisitor", cart.GetCartForVisitor)

	engine.Set("db.getEngine", db.GetDBEngine)
	engine.Set("db.getCollection", db.GetCollection)

	engine.Set("env.getConfig", env.GetConfig)
	engine.Set("env.getIniConfig", env.GetIniConfig)
	engine.Set("env.getLogger", env.GetLogger)
	engine.Set("env.getScheduler", env.GetScheduler)
	engine.Set("env.getErrorBus", env.GetErrorBus)
	engine.Set("env.getEventBus", env.GetEventBus)
	engine.Set("env.getScriptEngine", env.GetScriptEngine)

	engine.Set("env.configValue", env.ConfigGetValue)
	engine.Set("env.iniValue", env.IniValue)
	engine.Set("env.log", env.Log)
	engine.Set("env.logError", env.LogError)
	engine.Set("env.logEvent", env.LogEvent)
	engine.Set("env.getErrorLevel", env.ErrorLevel)
	engine.Set("env.getErrorCode", env.ErrorCode)
	engine.Set("env.getErrorMessage", env.ErrorMessage)
	engine.Set("env.error.registerListener", env.ErrorRegisterListener)
	engine.Set("env.error.dispatch", env.ErrorRegisterListener)
	engine.Set("env.error.registerListener", env.ErrorRegisterListener)
	engine.Set("env.error.dispatch", env.ErrorRegisterListener)
	engine.Set("env.error.new", env.ErrorNew)
	engine.Set("env.event.registerListener", env.EventRegisterListener)
	engine.Set("env.event.dispatch", env.Event)

	engine.Set("media.get", media.GetMediaStorage)

	engine.Set("impex.importCSV", impex.ImportCSVData)
	engine.Set("impex.importCSVScript", impex.ImportCSVScript)
	engine.Set("impex.mapToCSV", impex.MapToCSV)
	engine.Set("impex.CSVToMap", impex.CSVToMap)

	engine.Set("utils.cryptToURLString", utils.CryptToURLString)
	engine.Set("utils.decryptURLString", utils.DecryptURLString)
	engine.Set("utils.passwordEncode", utils.PasswordEncode)
	engine.Set("utils.passwordCheck", utils.PasswordCheck)
	engine.Set("utils.encryptData", utils.EncryptData)
	engine.Set("utils.decryptData", utils.DecryptData)

	engine.Set("utils.isZeroTime", utils.IsZeroTime)
	engine.Set("utils.isMD5", utils.IsMD5)
	engine.Set("utils.isAmongStr", utils.IsAmongStr)
	engine.Set("utils.isInArray", utils.IsInArray)
	engine.Set("utils.isInListStr", utils.IsInListStr)
	engine.Set("utils.isBlank", utils.CheckIsBlank)
	engine.Set("utils.stringToFloat", utils.StringToFloat)
	engine.Set("utils.stringToInteger", utils.StringToInteger)
	engine.Set("utils.stringToType", utils.StringToType)
	engine.Set("utils.interfaceToBool", utils.InterfaceToBool)
	engine.Set("utils.interfaceToFloat64", utils.InterfaceToFloat64)
	engine.Set("utils.interfaceToInt", utils.InterfaceToInt)
	engine.Set("utils.interfaceToMap", utils.InterfaceToMap)
	engine.Set("utils.interfaceToString", utils.InterfaceToString)
	engine.Set("utils.interfaceToStringArray", utils.InterfaceToStringArray)
	engine.Set("utils.interfaceToTime", utils.InterfaceToTime)
	engine.Set("utils.interfaceToMap", utils.InterfaceToMap)

	engine.Set("utils.getTemplateFunctions", utils.GetTemplateFunctions)
	engine.Set("utils.registerTemplateFunction", utils.RegisterTemplateFunction)
	engine.Set("utils.textTemplate", utils.TextTemplate)

	engine.Set("utils.timezones", utils.TimeZones)
	engine.Set("utils.parseTimeZone", utils.ParseTimeZone)
	engine.Set("utils.makeTZTime", utils.MakeTZTime)
	engine.Set("utils.timeToUTCTime", utils.TimeToUTCTime)

	engine.Set("utils.getPointer", utils.GetPointer)
	engine.Set("utils.syncGet", utils.SyncGet)
	engine.Set("utils.syncSet", utils.SyncSet)
	engine.Set("utils.syncMutex", utils.SyncMutex)
	engine.Set("utils.syncLock", utils.SyncLock)
	engine.Set("utils.syncUnlock", utils.SyncUnlock)

	engine.Set("utils.sortMapByKeys", utils.SortMapByKeys)

	engine.Set("utils.encodeToJSONString", utils.EncodeToJSONString)
	engine.Set("utils.decodeJSONToArray", utils.DecodeJSONToArray)
	engine.Set("utils.decodeJSONToInterface", utils.DecodeJSONToInterface)
	engine.Set("utils.DecodeJSONToStringKeyMap", utils.DecodeJSONToStringKeyMap)

	engine.Set("utils.KeysInMapAndNotBlank", utils.KeysInMapAndNotBlank)
	engine.Set("utils.GetFirstMapValue", utils.GetFirstMapValue)
	engine.Set("utils.Explode", utils.Explode)
	engine.Set("utils.Round", utils.Round)
	engine.Set("utils.RoundPrice", utils.RoundPrice)
	engine.Set("utils.SplitQuotedStringBy", utils.SplitQuotedStringBy)
	engine.Set("utils.MatchMapAValuesToMapB", utils.MatchMapAValuesToMapB)
	engine.Set("utils.EscapeRegexSpecials", utils.EscapeRegexSpecials)
	engine.Set("utils.ValidEmailAddress", utils.ValidEmailAddress)
	engine.Set("utils.Clone", utils.Clone)
	engine.Set("utils.StrToSnakeCase", utils.StrToSnakeCase)
	engine.Set("utils.StrToCamelCase", utils.StrToCamelCase)
	engine.Set("utils.MapGetPathValue", utils.MapGetPathValue)
	engine.Set("utils.MapSetPathValue", utils.MapSetPathValue)

	env.RegisterScriptEngine("Otto", engine)
}
