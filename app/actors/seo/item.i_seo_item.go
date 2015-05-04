package seo

import (
	"github.com/ottemo/foundation/app/models/seo"
	"github.com/ottemo/foundation/env"
)

// GetObjectID returns SEO item associated object ID
func (it *DefaultSEOItem) GetObjectID() string {
	return it.Rewrite
}

// SetObjectID specifies associated object ID for current SEO item
func (it *DefaultSEOItem) SetObjectID(objectID string) error {
	it.Rewrite = objectID
	return nil
}

// GetType returns SEOType of current SEO item
func (it *DefaultSEOItem) GetType() string {
	return it.Type
}

// SetType specifies SEOType for current SEO item
func (it *DefaultSEOItem) SetType(newType string) error {
	if seo.IsSEOType(newType) {
		it.Type = newType
		return nil
	}
	return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e6af9084-fddc-45bf-a90c-9c3d6ff88a57", "unknown seo type")
}

// GetURL returns URL of current SEO item
func (it *DefaultSEOItem) GetURL() string {
	return it.URL
}

// SetURL specifies URL for current SEO item
func (it *DefaultSEOItem) SetURL(newURL string) error {
	it.URL = newURL
	return nil
}
