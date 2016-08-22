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

// DefaultSEOItem is a default implementer of InterfaceSEOItem
type DefaultSEOItem struct {
	id string

	Url string
	Rewrite string

	Type string
	Title string
	MetaKeywords string
	MetaDescription string
}
