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

// InterfaceOrderItem represents interface to access business layer implementation of SEO item object
type InterfaceSEOEngine interface {
	GetObjectSEO(seoType string, objectID string)
	SetObjectSEO(seoType string, objectID string, seoItem InterfaceSEOItem) error

	GetUrlSEO(urlPattern string) []InterfaceSEOItem


}

// InterfaceOrderItem represents interface to access business layer implementation of SEO item object
type InterfaceSEOItem interface {
	GetID() string
	SetID(newID string) error

	GetObjectID() string
	SetObjectID(objectID string) error

	GetType() string
	SetType(newType string) error

	GetURL() string
	SetURL(newURL string) error

	models.InterfaceModel
	models.InterfaceObject
	models.InterfaceStorable
}