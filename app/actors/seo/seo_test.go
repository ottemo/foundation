package seo_test

import (
	"testing"

	"github.com/ottemo/foundation/test"
	"github.com/ottemo/foundation/app/models/seo"
	"github.com/ottemo/foundation/utils"
)

// Basic SEOItem test: check model save/load
func TestSEO(t *testing.T) {
	err := test.StartAppInTestingMode()
	if err != nil {
		t.Error(err)
		return
	}

	seoItemData, err := utils.DecodeJSONToStringKeyMap(`{
		"url": "url value",
		"rewrite": "rewrite value",
		"type": "type value",
		"title": "title value",
		"meta keywords": "keywords, value",
		"meta_description": "description value"
	}`)
	if err != nil {
		t.Error(err)
		return
	}

	seoItemModel, err := seo.GetSEOItemModel()
	if err != nil {
		t.Error(err)
		return
	}

	err = seoItemModel.FromHashMap(seoItemData)
	if err != nil {
		t.Error(err)
		return
	}

	err = seoItemModel.Save()
	if err != nil {
		t.Error(err)
		return
	}
	defer seoItemModel.Delete()

	seoItemID := seoItemModel.GetID()

	seoItemTestModel, err := seo.GetSEOItemModel()
	if err != nil {
		t.Error(err)
		return
	}

	err = seoItemTestModel.Load(seoItemID)
	if err != nil {
		t.Error(err)
		return
	}

	if seoItemTestModel.GetURL() != seoItemModel.GetURL() {
		t.Error("fail: url")
	}

	if seoItemTestModel.GetRewrite() != seoItemModel.GetRewrite() {
		t.Error("fail: rewrite")
	}

	if seoItemTestModel.GetType() != seoItemModel.GetType() {
		t.Error("fail: type")
	}

	if seoItemTestModel.GetTitle() != seoItemModel.GetTitle() {
		t.Error("fail: title")
	}

	if seoItemTestModel.GetMetaDescription() != seoItemModel.GetMetaDescription() {
		t.Error("fail: meta description")
	}

	if seoItemTestModel.GetMetaKeywords() != seoItemModel.GetMetaKeywords() {
		t.Error("fail: meta keywords")
	}
}
