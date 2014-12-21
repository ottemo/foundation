package product

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/media"
)

// AddMedia adds new media assigned to product
func (it *DefaultProduct) AddMedia(mediaType string, mediaName string, content []byte) error {
	productID := it.GetID()
	if productID == "" {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "856507153acf4e47a365c6e8911d9118", "product id not set")
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
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "87bb383acf3548e09d50ad517ed2e8f9", "product id not set")
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
		return result, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "45b1ebde3dd04c6c9960fddd89f4907f", "product id not set")
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
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "5f5d3c33de824580a6e7f5c45e9281e5", "product id not set")
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
		return "", env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0055f93a5d1041db8d93ea2bb4bee216", "product id not set")
	}

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	return mediaStorage.GetMediaPath(it.GetModelName(), productID, mediaType)
}
