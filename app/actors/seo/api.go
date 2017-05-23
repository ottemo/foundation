package seo

import (
	"os"
	"strings"
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/category"
	"github.com/ottemo/foundation/app/models/cms"
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/app/models/seo"
	"io"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	service := api.GetRestService()

	service.GET("seo/items", APIListSEOItems)
	service.GET("seo/attributes", api.IsAdminHandler(APIListSeoAttributes))

	service.GET("seo/url", APIGetSEOItem)
	service.GET("seo/canonical/:id", APIGetSEOItemByID)

	service.GET("seo/sitemap", APIGenerateSitemap)
	service.GET("seo/sitemap/sitemap.xml", APIGetSitemap)

	// Admin Only
	service.POST("seo/item", api.IsAdminHandler(APICreateSEOItem))
	service.PUT("seo/item/:itemID", api.IsAdminHandler(APIUpdateSEOItem))
	service.DELETE("seo/item/:itemID", api.IsAdminHandler(APIDeleteSEOItem))

	return nil
}

// APIListSEOItems returns a list registered SEO records
func APIListSEOItems(context api.InterfaceApplicationContext) (interface{}, error) {

	// retrieve collection model
	seoItemCollectionModel, err := GetSEOItemCollectionModel()
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	// filters handle
	if err := models.ApplyFilters(context, seoItemCollectionModel.GetDBCollection()); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "8f41ce59-fd0d-4755-b5f8-a11613adf9bc", err.Error())
	}

	// check "count" request
	if context.GetRequestArgument(api.ConstRESTActionParameter) == "count" {
		return seoItemCollectionModel.GetDBCollection().Count()
	}

	// limit parameter handle
	if err := seoItemCollectionModel.ListLimit(models.GetListLimit(context)); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "9b662ffc-1ef4-4f5f-ac97-552554321536", err.Error())
	}

	// extra parameter handle
	if err := models.ApplyExtraAttributes(context, seoItemCollectionModel); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "7607d46c-9700-4380-875f-56cce8c550cf", err.Error())
	}

	listItems, err := seoItemCollectionModel.List()
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	return listItems, nil
}

// APIListSeoAttributes returns a list of seo item attributes
func APIListSeoAttributes(context api.InterfaceApplicationContext) (interface{}, error) {

	seoItemModel, err := seo.GetSEOItemModel()
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	return seoItemModel.GetAttributesInfo(), nil
}

// APIListSEOItemsAlt returns a list registered SEO records
func APIListSEOItemsAlt(context api.InterfaceApplicationContext) (interface{}, error) {

	collection, err := db.GetCollection(ConstCollectionNameURLRewrites)
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	// If you give us a url to match we are only going to return one item
	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	for key := range requestData {
		switch key {
		case "url":
			if err := collection.AddFilter("url", "=", context.GetRequestArgument("url")); err != nil {
				_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ce480db1-58b1-4189-b7ec-c5a10c7197ac", err.Error())
			}
		}
	}

	records, err := collection.Load()
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	return records, nil
}

// APIGetSEOItem returns SEO item for a specified url
//   - SEO url should be specified in "url" argument
func APIGetSEOItem(context api.InterfaceApplicationContext) (interface{}, error) {

	collection, err := db.GetCollection(ConstCollectionNameURLRewrites)
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	// No starting slashes pls
	specifiedURL := context.GetRequestArgument("url")
	specifiedURL = strings.Trim(specifiedURL, "/")
	if err := collection.AddFilter("url", "=", specifiedURL); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "3902720b-ad39-4b17-a8bc-d80989147808", err.Error())
	}
	records, err := collection.Load()

	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	return records, nil
}

// APIGetSEOItemByID returns SEO item for a specified id
func APIGetSEOItemByID(context api.InterfaceApplicationContext) (interface{}, error) {

	// checking request context
	//-------------------------
	seoItemID := context.GetRequestArgument("id")
	if seoItemID == "" {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "2070ec2c-43c6-4a98-98aa-6334e684a23a", "Required field 'id' is blank or absend.")
	}

	// operation
	//-------------------------
	seoItemModel, err := seo.GetSEOItemModel()
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	err = seoItemModel.Load(seoItemID)
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	return seoItemModel.ToHashMap(), nil
}

// APIUpdateSEOItem updates existing SEO item
//   - SEO item id should be specified in "itemID" argument
func APIUpdateSEOItem(context api.InterfaceApplicationContext) (interface{}, error) {

	postValues, err := api.GetRequestContentAsMap(context)
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	collection, err := db.GetCollection(ConstCollectionNameURLRewrites)
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	urlRewriteID := context.GetRequestArgument("itemID")
	record, err := collection.LoadByID(urlRewriteID)
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	// if rewrite 'url' was changed - checking new value for duplicates
	//-----------------------------------------------------------------
	if urlValue, present := postValues["url"]; present && urlValue != record["url"] {
		urlValue := utils.InterfaceToString(urlValue)

		if err := collection.AddFilter("url", "=", urlValue); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "4183d389-8cca-4be7-b3a4-abf10fe06500", err.Error())
		}
		recordsNumber, err := collection.Count()
		if err != nil {
			context.SetResponseStatusInternalServerError()
			return nil, env.ErrorDispatch(err)
		}
		if recordsNumber > 0 {
			context.SetResponseStatusInternalServerError()
			return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "c2a2e89d-b358-4c3b-9b65-4d161188b592", "rewrite for url '"+urlValue+"' already exists")
		}

		record["url"] = urlValue
	}

	// updating other attributes
	//--------------------------
	attributes := []string{"type", "rewrite", "title", "meta_keywords", "meta_description"}
	for _, attribute := range attributes {
		if value, present := postValues[attribute]; present {
			record[attribute] = value
		}
	}

	// saving updates
	//---------------
	_, err = collection.Save(record)
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	return record, nil
}

// APICreateSEOItem creates a new SEO item
//   - "url" and "rewrite" attributes are required
func APICreateSEOItem(context api.InterfaceApplicationContext) (interface{}, error) {

	postValues, err := api.GetRequestContentAsMap(context)
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	if !utils.KeysInMapAndNotBlank(postValues, "url", "rewrite") {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "1a3901a1-48d5-4055-bacb-2b02681bbb71", "'url' and 'rewrite' context should be specified")
	}

	valueURL := utils.InterfaceToString(postValues["url"])
	valueRewrite := utils.InterfaceToString(postValues["rewrite"])

	// looking for duplicated 'url'
	//-----------------------------
	collection, err := db.GetCollection(ConstCollectionNameURLRewrites)
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	if err := collection.AddFilter("url", "=", valueURL); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ef9c3c87-10d0-4341-8bee-8dba22487905", err.Error())
	}
	recordsNumber, err := collection.Count()
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}
	if recordsNumber > 0 {
		context.SetResponseStatusBadRequest()
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "77987a83-3420-4baf-99f0-af9c47689d3b", "rewrite for url '"+valueURL+"' already exists")
	}

	// making new record and storing it
	//---------------------------------
	newRecord := map[string]interface{}{
		"url":              valueURL,
		"type":             "",
		"rewrite":          valueRewrite,
		"title":            nil,
		"meta_keywords":    nil,
		"meta_description": nil,
	}

	attributes := []string{"type", "title", "meta_keywords", "meta_description"}
	for _, attribute := range attributes {
		if value, present := postValues[attribute]; present {
			newRecord[attribute] = value
		}
	}

	newID, err := collection.Save(newRecord)
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	newRecord["_id"] = newID

	return newRecord, nil
}

// APIDeleteSEOItem deletes specified SEO item
//   - SEO item id should be specified in "itemID" argument
func APIDeleteSEOItem(context api.InterfaceApplicationContext) (interface{}, error) {

	collection, err := db.GetCollection(ConstCollectionNameURLRewrites)
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	err = collection.DeleteByID(context.GetRequestArgument("itemID"))
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	return "ok", nil
}

// APIGetSitemap returns SEO records based sitemap (auto re-generating it if needed)
//   - result is not a JSON but "text/xml"
func APIGetSitemap(context api.InterfaceApplicationContext) (interface{}, error) {

	// if sitemap expired - generate new one
	info, err := os.Stat(ConstSitemapFilePath)
	if err != nil || (time.Now().Unix()-info.ModTime().Unix() >= ConstSitemapExpireSec) {
		return APIGenerateSitemap(context)
	}

	// using generated otherwise
	sitemapFile, err := os.Open(ConstSitemapFilePath)
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}
	defer func(c io.Closer){
		if err := c.Close(); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "b9b14fb0-43c6-434c-b075-f73428e9285e", err.Error())
		}
	}(sitemapFile)

	if responseWriter := context.GetResponseWriter(); responseWriter != nil {
		if err := context.SetResponseContentType("text/xml"); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "4ad383f8-065f-420a-a197-9011d88f8efa", err.Error())
		}

		buffer := make([]byte, 1024)
		for {
			n, err := sitemapFile.Read(buffer)
			if err != nil || n == 0 {
				break
			}

			n, err = responseWriter.Write(buffer[0:n])
			if err != nil || n == 0 {
				break
			}
		}
	}

	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}

	return nil, nil
}

// APIGenerateSitemap generates a new sitemap based on SEO records
//   - generates sitemap any time called (no cache used)
//   - result is not a JSON but "text/xml"
func APIGenerateSitemap(context api.InterfaceApplicationContext) (interface{}, error) {

	if err := context.SetResponseContentType("text/xml"); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "f3c9b7e7-d0fd-4f7d-9a42-86a78f828231", err.Error())
	}
	responseWriter := context.GetResponseWriter()

	// creating sitemap file
	sitemapFile, err := os.Create(ConstSitemapFilePath)
	if err != nil {
		context.SetResponseStatusInternalServerError()
		return nil, env.ErrorDispatch(err)
	}
	defer func(c io.Closer){
		if err := c.Close(); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "5a85194f-7963-42a2-9d07-c97c463ae66b", err.Error())
		}
	}(sitemapFile)

	// writer to write into HTTP and file simultaneously
	newline := []byte("\n")
	writeLine := func(line []byte) {
		if responseWriter != nil {
			if _, err := responseWriter.Write(line); err != nil {
				_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "71fc4613-626a-4720-aac2-a13a188cf51e", err.Error())
			}
			if _, err := responseWriter.Write(newline); err != nil {
				_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "a8b5b411-d1a3-4481-abd0-e8fff0e1f642", err.Error())
			}
		}

		if _, err := sitemapFile.Write(line); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "26265d8c-3293-4319-b1f3-5f1073b07cbb", err.Error())
		}
		if _, err := sitemapFile.Write(newline); err != nil {
			_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "50959d06-083d-4176-9aef-69f27b6cffc5", err.Error())
		}
	}

	// sitemap file preparations
	writeLine([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>"))
	writeLine([]byte("<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">"))

	baseURL := app.GetStorefrontURL("")
	rewriteType := ""

	// per database record iterator
	iteratorFunc := func(record map[string]interface{}) bool {
		pageURL := ""
		if rewriteType == "" {
			pageURL = baseURL + utils.InterfaceToString(record["url"])
		} else {
			pageURL = baseURL + rewriteType + "/" + utils.InterfaceToString(record["_id"])
		}

		writeLine([]byte("  <url><loc>" + pageURL + "</loc></url>"))

		return true
	}

	// Re-writed pages
	rewritesCollection, err := db.GetCollection(ConstCollectionNameURLRewrites)
	if err := rewritesCollection.SetResultColumns("url"); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "33b433a2-e50e-4ae4-a509-8ba147779202", err.Error())
	}
	if err := rewritesCollection.Iterate(iteratorFunc); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "5fec2d5d-388b-4ae6-8bb9-7f19dd1eccc8", err.Error())
	}

	if err := rewritesCollection.SetResultColumns("rewrite"); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "4f55972e-323f-4cc4-81e1-517b5d94f8d5", err.Error())
	}

	// Product pages
	rewriteType = "product"
	if err := rewritesCollection.AddFilter("type", "=", rewriteType); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "1add17b7-e505-404a-8905-04332be27ff6", err.Error())
	}

	productCollectionModel, _ := product.GetProductCollectionModel()
	dbProductCollection := productCollectionModel.GetDBCollection()
	if err := dbProductCollection.SetResultColumns("_id"); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "063b59b0-37da-4c9a-b13b-a1277c110598", err.Error())
	}
	if err := dbProductCollection.AddFilter("_id", "nin", rewritesCollection); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "ee74ae99-a295-47de-9fea-366045a75aca", err.Error())
	}
	if err := dbProductCollection.Iterate(iteratorFunc); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "510d47f9-66a6-4e49-a05f-016e3099aaf4", err.Error())
	}

	// Category pages
	rewriteType = "category"
	if err := rewritesCollection.ClearFilters(); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "1517d748-9784-49be-8021-d3d99d68c52e", err.Error())
	}
	if err := rewritesCollection.AddFilter("type", "=", rewriteType); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "09452416-ec3c-4c48-b952-83ce528e20c3", err.Error())
	}

	categoryCollectionModel, _ := category.GetCategoryCollectionModel()
	dbCategoryCollection := categoryCollectionModel.GetDBCollection()
	if err := dbCategoryCollection.SetResultColumns("_id"); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "85e91980-eb5f-4f44-94b2-ae3cda6a7420", err.Error())
	}
	if err := dbCategoryCollection.AddFilter("_id", "nin", rewritesCollection); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "dd334247-273d-4dcd-81d6-8e2a68c5c09c", err.Error())
	}
	if err := dbCategoryCollection.Iterate(iteratorFunc); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "cd3958c9-3ba5-45fc-881d-86e5c6ffab12", err.Error())
	}

	// Cms pages
	rewriteType = "page"
	if err := rewritesCollection.ClearFilters(); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "eb511139-07dd-439e-926e-3107bac96281", err.Error())
	}
	if err := rewritesCollection.AddFilter("type", "=", rewriteType); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "3f742f50-489d-41f8-b6f6-9dbbe204a7a9", err.Error())
	}

	cmsPageCollectionModel, _ := cms.GetCMSPageCollectionModel()
	dbCMSPageCollection := cmsPageCollectionModel.GetDBCollection()
	if err := dbCMSPageCollection.SetResultColumns("_id"); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "1fae19a4-d2a3-47f4-a028-128a87fbf925", err.Error())
	}
	if err := dbCMSPageCollection.AddFilter("_id", "nin", rewritesCollection); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "484f4837-8cf7-49aa-922f-1c872642b3db", err.Error())
	}
	if err := dbCMSPageCollection.Iterate(iteratorFunc); err != nil {
		_ = env.ErrorNew(ConstErrorModule, ConstErrorLevel, "7b6cd323-5a7b-4c93-86ba-ea0833af69f2", err.Error())
	}

	writeLine([]byte("</urlset>"))

	return nil, nil
}
