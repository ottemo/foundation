package category

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/media"
)

// AddMedia adds new media assigned to category
func (it *DefaultCategory) AddMedia(mediaType string, mediaName string, content []byte) error {
	categoryID := it.GetID()
	if categoryID == "" {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "4ce5dc3e-db9d-44df-ba82-29daf7fabf9a", "category id not set")
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
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "05d2b165-ad76-4376-81c2-e9cc3fb8eac5", "category id not set")
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
		return result, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e9a39d29-bc54-42d4-aa65-c0eb043be822", "category id not set")
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
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "dca7abba-3814-49a8-8e23-80607e6ad1ee", "category id not set")
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
		return "", env.ErrorNew(ConstErrorModule, ConstErrorLevel, "23685b86-37e9-4151-86fd-f1c7d662b450", "category id not set")
	}

	mediaStorage, err := media.GetMediaStorage()
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	return mediaStorage.GetMediaPath(it.GetModelName(), categoryID, mediaType)
}
