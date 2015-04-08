package category

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/media"
)

// AddMedia adds new media assigned to category
func (it *DefaultCategory) AddMedia(mediaType string, mediaName string, content []byte) error {
	categoryID := it.GetID()
	if categoryID == "" {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "85650715-3acf-4e47-a365-c6e8911d9118", "category id not set")
	}

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return mediaStorage.Save(it.GetModelName(), categoryID, mediaType, mediaName, content)
}

// RemoveMedia removes media assigned to category
func (it *DefaultCategory) RemoveMedia(mediaType string, mediaName string) error {
	categoryID := it.GetID()
	if categoryID == "" {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "87bb383a-cf35-48e0-9d50-ad517ed2e8f9", "category id not set")
	}

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return mediaStorage.Remove(it.GetModelName(), categoryID, mediaType, mediaName)
}

// ListMedia lists media assigned to category
func (it *DefaultCategory) ListMedia(mediaType string) ([]string, error) {
	var result []string

	categoryID := it.GetID()
	if categoryID == "" {
		return result, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "45b1ebde-3dd0-4c6c-9960-fddd89f4907f", "category id not set")
	}

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return result, env.ErrorDispatch(err)
	}

	return mediaStorage.ListMedia(it.GetModelName(), categoryID, mediaType)
}

// GetMedia returns content of media assigned to category
func (it *DefaultCategory) GetMedia(mediaType string, mediaName string) ([]byte, error) {
	categoryID := it.GetID()
	if categoryID == "" {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "5f5d3c33-de82-4580-a6e7-f5c45e9281e5", "category id not set")
	}

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return mediaStorage.Load(it.GetModelName(), categoryID, mediaType, mediaName)
}

// GetMediaPath returns relative location of media assigned to category in media storage
func (it *DefaultCategory) GetMediaPath(mediaType string) (string, error) {
	categoryID := it.GetID()
	if categoryID == "" {
		return "", env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0055f93a-5d10-41db-8d93-ea2bb4bee216", "category id not set")
	}

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	return mediaStorage.GetMediaPath(it.GetModelName(), categoryID, mediaType)
}
