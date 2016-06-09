package fsmedia

import (
	"os"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/media"
	"github.com/ottemo/foundation/utils"
)

// init makes package self-initialization routine
func init() {
	instance := new(FilesystemMediaStorage)

	if err := media.RegisterMediaStorage(instance); err == nil {
		instance.imageSizes = make(map[string]string)
		instance.setupWaitCnt = 3

		env.RegisterOnConfigIniStart(instance.setupOnIniConfigStart)
		env.RegisterOnConfigStart(instance.setupConfig)
		db.RegisterOnDatabaseStart(instance.setupOnDatabaseStart)

		api.RegisterOnRestServiceStart(setupAPI)

		// process of resizing images on media start
		//		media.RegisterOnMediaStorageStart(instance.ResizeAllMediaImages)
	}
}

// setupCheckDone performs callback event if setup was done
func (it *FilesystemMediaStorage) setupCheckDone() {

	// so, we are not sure on events sequence order
	if it.setupWaitCnt--; it.setupWaitCnt == 0 {
		err := media.OnMediaStorageStart()
		if err != nil {
			env.ErrorDispatch(err)
		}
	}
}

// setupOnIniConfigStart is a initialization based on ini config service
func (it *FilesystemMediaStorage) setupOnIniConfigStart() error {

	var storageFolder = ConstMediaDefaultFolder

	if iniConfig := env.GetIniConfig(); iniConfig != nil {
		if iniValue := iniConfig.GetValue("media.fsmedia.folder", "?"+ConstMediaDefaultFolder); iniValue != "" {
			storageFolder = iniValue
		}

		if iniValue := iniConfig.GetValue("media.resize.images.onfly", "false"); iniValue != "" {
			resizeImagesOnFly = utils.InterfaceToBool(iniValue)
		}
	}

	err := os.MkdirAll(storageFolder, os.ModePerm)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	it.storageFolder = storageFolder

	if it.storageFolder != "" && it.storageFolder[len(it.storageFolder)-1] != '/' {
		it.storageFolder += "/"
	}

	it.setupCheckDone()

	return nil
}

// setupOnDatabaseStart is a initialization based on config service
func (it *FilesystemMediaStorage) setupOnDatabaseStart() error {

	dbEngine := db.GetDBEngine()
	if dbEngine == nil {
		return env.ErrorNew(ConstErrorModule,
			env.ConstErrorLevelStartStop,
			"b50f7128-a32f-4d92-866e-1ee35ba079df",
			"Unable to find database engine specified in configuration file to start.")
	}

	dbCollection, err := dbEngine.GetCollection(ConstMediaDBCollection)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	dbCollection.AddColumn("model", db.ConstTypeVarchar, true)
	dbCollection.AddColumn("object", db.ConstTypeVarchar, true)
	dbCollection.AddColumn("type", db.ConstTypeVarchar, true)
	dbCollection.AddColumn("media", db.ConstTypeVarchar, false)
	dbCollection.AddColumn("created_at", db.ConstTypeDatetime, false)
	it.setupCheckDone()

	return nil
}
