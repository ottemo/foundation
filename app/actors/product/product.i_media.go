package product

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/media"
)

// AddMedia adds new media assigned to product
func (it *DefaultProduct) AddMedia(mediaType string, mediaName string, content []byte) error {
	productID := it.GetID()
	if productID == "" {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "85650715-3acf-4e47-a365-c6e8911d9118", "product id not set")
	}

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return mediaStorage.Save(it.GetModelName(), productID, mediaType, mediaName, content)
}

// RemoveMedia removes media assigned to product
func (it *DefaultProduct) RemoveMedia(mediaType string, mediaName string) error {
	productID := it.GetID()
	if productID == "" {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "87bb383a-cf35-48e0-9d50-ad517ed2e8f9", "product id not set")
	}

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return mediaStorage.Remove(it.GetModelName(), productID, mediaType, mediaName)
}

// ListMedia lists media assigned to product
func (it *DefaultProduct) ListMedia(mediaType string) ([]string, error) {
	var result []string

	productID := it.GetID()
	if productID == "" {
		return result, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "45b1ebde-3dd0-4c6c-9960-fddd89f4907f", "product id not set")
	}

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return result, env.ErrorDispatch(err)
	}

	return mediaStorage.ListMedia(it.GetModelName(), productID, mediaType)
}

// GetMedia returns content of media assigned to product
func (it *DefaultProduct) GetMedia(mediaType string, mediaName string) ([]byte, error) {
	productID := it.GetID()
	if productID == "" {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "5f5d3c33-de82-4580-a6e7-f5c45e9281e5", "product id not set")
	}

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return mediaStorage.Load(it.GetModelName(), productID, mediaType, mediaName)
}

// GetMediaPath returns relative location of media assigned to product in media storage
func (it *DefaultProduct) GetMediaPath(mediaType string) (string, error) {
	productID := it.GetID()
	if productID == "" {
		return "", env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0055f93a-5d10-41db-8d93-ea2bb4bee216", "product id not set")
	}

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	return mediaStorage.GetMediaPath(it.GetModelName(), productID, mediaType)
}
