package fsmedia

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/media"
	"github.com/ottemo/foundation/utils"
)

// configures package related API endpoint routines
func setupAPI() error {

	var err error

	err = api.GetRestService().RegisterAPI("media", api.GET, APIGetMediaInfo)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// APIGetMediaInfo will resize all images if the params of the request contain 'resizeAll' with a value of true
func APIGetMediaInfo(context api.InterfaceApplicationContext) (interface{}, error) {
	// TODO: add example api call or add this to Apiary - jwv

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	requestParams := context.GetRequestArguments()

	resizeAll := utils.GetFirstMapValue(requestParams, "resizeAll", "resizeImages", "resizeAllImages")

	if resizeAll != nil && utils.InterfaceToBool(resizeAll) {
		mediaStorage, err := media.GetMediaStorage()
		if err != nil {
			return nil, err
		}

		err = mediaStorage.ResizeAllMediaImages()
		if err != nil {
			return nil, err
		}
	}

	return "ok", nil
}
