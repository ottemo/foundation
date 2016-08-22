package seo

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstModelNameSEOItem = "SEOItem"

	ConstErrorModule = "seo"
	ConstErrorLevel  = env.ConstErrorLevelModel
)

// InterfaceSEOEngine represents interface to access business layer implementation of SEO engine
type InterfaceSEOEngine interface {
	GetSEO(seoType string, objectID string, urlPattern string) []InterfaceSEOItem
}

// InterfaceSEOItem represents interface to access business layer implementation of SEO item object
type InterfaceSEOItem interface {
	GetURL() string
	SetURL(newURL string) error

	GetRewrite() string

	GetTitle() string
	GetType() string
	GetMetaKeywords() string
	GetMetaDescription() string

	models.InterfaceModel
	models.InterfaceObject
	models.InterfaceStorable
}
