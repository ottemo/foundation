package media

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/media"
	"github.com/ottemo/foundation/utils"
)

// init makes package self-initialization routine
func init() {
	api.RegisterOnRestServiceStart(setupAPI)
	app.OnAppStart(onAppStart)

	utils.RegisterTemplateFunction("media", mediaTemplateDirective)
}

func onAppStart() error {
	mediaStorageInstance, err := media.GetMediaStorage()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	mediaStorage = mediaStorageInstance
	return nil
}

// mediaTemplateDirective - for adding image to pages
//   use {{media "mediaName" .}} to fetch image URL
func mediaTemplateDirective(args ...interface{}) (string, error) {
	mediaName := ""
	if len(args) > 0 {
		mediaName = utils.InterfaceToString(args[0])
	}
	imagePath, err := mediaStorage.GetMediaPath(ConstStorageModel, ConstStorageObject, ConstStorageType)
	if err != nil {
		return "", env.ErrorDispatch(err)
	}

	foundationURL := app.GetFoundationURL(imagePath + mediaName)

	return "<img src=\"" + foundationURL + "\" />", nil
}
