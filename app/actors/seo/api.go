package seo

import (
	"github.com/ottemo/foundation/app"
	"os"
	"strings"
	"time"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models/category"
	"github.com/ottemo/foundation/app/models/cms"
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	service := api.GetRestService()

	service.GET("seo/items", APIListSEOItems)

	service.GET("seo/url", APIGetSEOItem)
	service.GET("seo/url/:url", APIGetSEOItem)
	service.GET("seo/canonical/:id", APIGetSEOItemByID)

	service.GET("seo/sitemap", APIGenerateSitemap)
	service.GET("seo/sitemap/sitemap.xml", APIGetSitemap)

	// Admin Only
	service.POST("seo/item", api.IsAdmin(APICreateSEOItem))
	service.PUT("seo/item/:itemID", api.IsAdmin(APIUpdateSEOItem))
	service.DELETE("seo/item/:itemID", api.IsAdmin(APIDeleteSEOItem))

	return nil
}

// APIListSEOItems returns a list registered SEO records
func APIListSEOItems(context api.InterfaceApplicationContext) (interface{}, error) {

	collection, err := db.GetCollection(ConstCollectionNameURLRewrites)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	collection.SetResultColumns("url", "type", "rewrite")
	records, err := collection.Load()

	return records, env.ErrorDispatch(err)
}

// APIListSEOItemsAlt returns a list registered SEO records
func APIListSEOItemsAlt(context api.InterfaceApplicationContext) (interface{}, error) {

	collection, err := db.GetCollection(ConstCollectionNameURLRewrites)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// If you give us a url to match we are only going to return one item
	requestData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for key := range requestData {
		switch key {
		case "url":
			collection.AddFilter("url", "=", context.GetRequestArgument("url"))
		}
	}

	records, err := collection.Load()

	return records, env.ErrorDispatch(err)
}

// APIGetSEOItem returns SEO item for a specified url
//   - SEO url should be specified in "url" argument
func APIGetSEOItem(context api.InterfaceApplicationContext) (interface{}, error) {

	collection, err := db.GetCollection(ConstCollectionNameURLRewrites)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// No starting slashes pls
	specifiedURL := context.GetRequestArgument("url")
	specifiedURL = strings.Trim(specifiedURL, "/")
	collection.AddFilter("url", "=", specifiedURL)
	records, err := collection.Load()

	return records, env.ErrorDispatch(err)
}

// APIGetSEOItemByID returns SEO item for a specified id
func APIGetSEOItemByID(context api.InterfaceApplicationContext) (interface{}, error) {

	collection, err := db.GetCollection(ConstCollectionNameURLRewrites)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	id := context.GetRequestArgument("id")
	records, err := collection.LoadByID(id)

	return records, env.ErrorDispatch(err)
}

// APIUpdateSEOItem updates existing SEO item
//   - SEO item id should be specified in "itemID" argument
func APIUpdateSEOItem(context api.InterfaceApplicationContext) (interface{}, error) {

	postValues, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	collection, err := db.GetCollection(ConstCollectionNameURLRewrites)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	urlRewriteID := context.GetRequestArgument("itemID")
	record, err := collection.LoadByID(urlRewriteID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// if rewrite 'url' was changed - checking new value for duplicates
	//-----------------------------------------------------------------
	if urlValue, present := postValues["url"]; present && urlValue != record["url"] {
		urlValue := utils.InterfaceToString(urlValue)

		collection.AddFilter("url", "=", urlValue)
		recordsNumber, err := collection.Count()
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}
		if recordsNumber > 0 {
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
		return nil, env.ErrorDispatch(err)
	}

	return record, nil
}

// APICreateSEOItem creates a new SEO item
//   - "url" and "rewrite" attributes are required
func APICreateSEOItem(context api.InterfaceApplicationContext) (interface{}, error) {

	postValues, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	if !utils.KeysInMapAndNotBlank(postValues, "url", "rewrite") {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "1a3901a1-48d5-4055-bacb-2b02681bbb71", "'url' and 'rewrite' context should be specified")
	}

	valueURL := utils.InterfaceToString(postValues["url"])
	valueRewrite := utils.InterfaceToString(postValues["rewrite"])

	// looking for duplicated 'url'
	//-----------------------------
	collection, err := db.GetCollection(ConstCollectionNameURLRewrites)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	collection.AddFilter("url", "=", valueURL)
	recordsNumber, err := collection.Count()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	if recordsNumber > 0 {
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
		return nil, env.ErrorDispatch(err)
	}

	err = collection.DeleteByID(context.GetRequestArgument("itemID"))
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return "ok", nil
}

// APIGetSitemap returns SEO records based sitemap (auto re-generating it if needed)
//   - result is not a JSON but "text/xml"
func APIGetSitemap(context api.InterfaceApplicationContext) (interface{}, error) {

	// if sitemap expied - generating new one
	info, err := os.Stat(ConstSitemapFilePath)
	if err != nil || (time.Now().Unix()-info.ModTime().Unix() >= ConstSitemapExpireSec) {
		return APIGenerateSitemap(context)
	}

	// using generated otherwise
	sitemapFile, err := os.Open(ConstSitemapFilePath)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	defer sitemapFile.Close()

	if responseWriter := context.GetResponseWriter(); responseWriter != nil {
		context.SetResponseContentType("text/xml")

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

	return nil, env.ErrorDispatch(err)
}

// APIGenerateSitemap generates a new sitemap based on SEO records
//   - generates sitemap any time called (no cache used)
//   - result is not a JSON but "text/xml"
func APIGenerateSitemap(context api.InterfaceApplicationContext) (interface{}, error) {

	context.SetResponseContentType("text/xml")
	responseWriter := context.GetResponseWriter()

	// creating sitemap file
	sitemapFile, err := os.Create(ConstSitemapFilePath)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	defer sitemapFile.Close()

	// writer to write into HTTP and file simultaneously
	newline := []byte("\n")
	writeLine := func(line []byte) {
		if responseWriter != nil {
			responseWriter.Write(line)
			responseWriter.Write(newline)
		}

		sitemapFile.Write(line)
		sitemapFile.Write(newline)
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
	rewritesCollection.SetResultColumns("url")
	rewritesCollection.Iterate(iteratorFunc)

	rewritesCollection.SetResultColumns("rewrite")

	// Product pages
	rewriteType = "product"
	rewritesCollection.AddFilter("type", "=", rewriteType)

	productCollectionModel, _ := product.GetProductCollectionModel()
	dbProductCollection := productCollectionModel.GetDBCollection()
	dbProductCollection.SetResultColumns("_id")
	dbProductCollection.AddFilter("_id", "nin", rewritesCollection)
	dbProductCollection.Iterate(iteratorFunc)

	// Category pages
	rewriteType = "category"
	rewritesCollection.ClearFilters()
	rewritesCollection.AddFilter("type", "=", rewriteType)

	categoryCollectionModel, _ := category.GetCategoryCollectionModel()
	dbCategoryCollection := categoryCollectionModel.GetDBCollection()
	dbCategoryCollection.SetResultColumns("_id")
	dbCategoryCollection.AddFilter("_id", "nin", rewritesCollection)
	dbCategoryCollection.Iterate(iteratorFunc)

	// Cms pages
	rewriteType = "page"
	rewritesCollection.ClearFilters()
	rewritesCollection.AddFilter("type", "=", rewriteType)

	cmsPageCollectionModel, _ := cms.GetCMSPageCollectionModel()
	dbCMSPageCollection := cmsPageCollectionModel.GetDBCollection()
	dbCMSPageCollection.SetResultColumns("_id")
	dbCMSPageCollection.AddFilter("_id", "nin", rewritesCollection)
	dbCMSPageCollection.Iterate(iteratorFunc)

	writeLine([]byte("</urlset>"))

	return nil, nil
}
