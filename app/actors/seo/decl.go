// Package seo implements a set of API intended to provide SEO optimizations
package seo

import (
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstCollectionNameURLRewrites = "url_rewrites"

	ConstSitemapFilePath  = "sitemap.xml"
	ConstSitemapExpireSec = 60 * 60 * 24

	ConstErrorModule = "seo"
	ConstErrorLevel  = env.ConstErrorLevelActor
)

// DefaultSEOEngine is a default implementer of InterfaceSEOEngine
type DefaultSEOEngine struct{}

// DefaultSEOItem is a default implementer of InterfaceSEOItem
type DefaultSEOItem struct {
	id string

	URL     string
	Type    string
	Rewrite string

	Title           string
	MetaKeywords    string
	MetaDescription string
}
