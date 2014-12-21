package impex

import (
	"io/ioutil"
	"net/http"
	"path"
	"strings"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// CheckModelImplements checks that model support InterfaceObject and InterfaceStorable interfaces
func CheckModelImplements(modelName string, neededInterfaces []string) (models.InterfaceModel, error) {
	cmdModel, err := models.GetModel(modelName)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for _, interfaceName := range neededInterfaces {
		ok := true
		switch interfaceName {
		case "InterfaceStorable":
			_, ok = cmdModel.(models.InterfaceStorable)
		case "InterfaceObject":
			_, ok = cmdModel.(models.InterfaceObject)
		case "InterfaceListable":
			_, ok = cmdModel.(models.InterfaceListable)
		case "InterfaceCollection":
			_, ok = cmdModel.(models.InterfaceCollection)
		case "InterfaceCustomAttributes":
			_, ok = cmdModel.(models.InterfaceCustomAttributes)
		case "InterfaceMedia":
			_, ok = cmdModel.(models.InterfaceMedia)
		}

		if !ok {
			return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "40d041484bd94e698748088b85bb7e0f", "model "+modelName+" not implements "+interfaceName)
		}
	}

	return cmdModel, nil
}

// TODO: make command parameters standardized parser to split required/optional parameters and get then in one function call

// ArgsGetAsNamed collects arguments into map, unnamed arguments will go as position index
func ArgsGetAsNamed(args []string, includeIndexes bool) map[string]string {
	result := make(map[string]string)
	for idx, arg := range args {
		splited := utils.SplitQuotedStringBy(arg, '=', ':')
		if len(splited) > 1 {
			key := splited[0]
			key = strings.Trim(strings.TrimSpace(key), "\"'`")

			value := strings.Join(splited[1:], " ")
			value = strings.Trim(strings.TrimSpace(value), "\"'`")

			result[key] = value
		} else {
			if includeIndexes {
				result[utils.InterfaceToString(idx)] = strings.Trim(strings.TrimSpace(arg), "\"'")
			}
		}
	}
	return result
}

// ArgsFindWorkingModel looks for model mention among command attributes
func ArgsFindWorkingModel(args []string, neededInterfaces []string) (models.InterfaceModel, error) {
	var result models.InterfaceModel
	var err error

	namedArgs := ArgsGetAsNamed(args, true)
	for _, argKey := range []string{"model", "1"} {
		if argValue, present := namedArgs[argKey]; present {
			result, err = CheckModelImplements(argValue, neededInterfaces)
			if err == nil {
				return result, nil
			}
		}
	}

	return nil, err
}

// ArgsFindIDKey looks for object identifier mention among command attributes
func ArgsFindIDKey(args []string) string {
	namedArgs := ArgsGetAsNamed(args, false)
	for _, checkingKey := range []string{"idKey", "id", "_id"} {
		if argValue, present := namedArgs[checkingKey]; present {
			return argValue
		}
	}
	return ""
}

// ArgsFindWorkingAttributes looks for model attributes inclusion/exclusion mention among command attributes
func ArgsFindWorkingAttributes(args []string) map[string]bool {
	result := make(map[string]bool)
	namedArgs := ArgsGetAsNamed(args, false)

	for _, argKey := range []string{"skip", "ignore", "use", "include", "attributes"} {
		if argValue, present := namedArgs[argKey]; present {
			for _, attributeName := range strings.Split(argValue, ",") {
				attributeName = strings.TrimSpace(attributeName)

				switch argKey {
				case "skip", "ignore":
					result[attributeName] = false
				default:
					result[attributeName] = true
				}
			}
		}
	}
	return result
}

// Init is a INSERT command initialization routine
func (it *ImportCmdInsert) Init(args []string, exchange map[string]interface{}) error {

	workingModel, err := ArgsFindWorkingModel(args, []string{"InterfaceStorable", "InterfaceObject"})
	if err != nil {
		return env.ErrorDispatch(err)
	}

	it.model = workingModel
	it.attributes = ArgsFindWorkingAttributes(args)

	namedArgs := ArgsGetAsNamed(args, false)
	if _, present := namedArgs["--skipErrors"]; present {
		it.skipErrors = true
	}

	return nil
}

// Process is a INSERT command processor
func (it *ImportCmdInsert) Process(itemData map[string]interface{}, input interface{}, exchange map[string]interface{}) (interface{}, error) {

	// preparing model
	//-----------------
	if it.model == nil {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "eb15163bbca941a792bcc646fe4eba53", "INSERT command have no assigned model to work on")
	}
	cmdModel, err := it.model.New()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	modelAsObject := cmdModel.(models.InterfaceObject)
	modelAsStorable := cmdModel.(models.InterfaceStorable)

	// filling model attributes
	//--------------------------
	for attribute, value := range itemData {
		if useAttribute, wasMentioned := it.attributes[attribute]; !wasMentioned || useAttribute {
			err := modelAsObject.Set(attribute, value)
			if err != nil && !it.skipErrors {
				return nil, err
			}
		}
	}

	// storing model
	//---------------
	err = modelAsStorable.Save()
	if err != nil {
		err = env.ErrorDispatch(err)

		if !it.skipErrors {
			return cmdModel, err
		}
	}

	return cmdModel, nil
}

// Init is a UPDATE command initialization routines
func (it *ImportCmdUpdate) Init(args []string, exchange map[string]interface{}) error {
	workingModel, err := ArgsFindWorkingModel(args, []string{"InterfaceStorable", "InterfaceObject"})
	if err != nil {
		return env.ErrorDispatch(err)
	}

	it.model = workingModel
	it.attributes = ArgsFindWorkingAttributes(args)
	it.idKey = ArgsFindIDKey(args)

	if it.model == nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "0ceb92d576fc4022866f6f905d299ab8", "INSERT command have no assigned model to work on")
	}

	if it.idKey == "" {
		it.idKey = "_id"
	}

	return nil
}

// Process is a UPDATE command processor
func (it *ImportCmdUpdate) Process(itemData map[string]interface{}, input interface{}, exchange map[string]interface{}) (interface{}, error) {

	// preparing model
	//-----------------
	cmdModel, err := it.model.New()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	modelAsObject := cmdModel.(models.InterfaceObject)
	modelAsStorable := cmdModel.(models.InterfaceStorable)

	if modelID, present := itemData[it.idKey]; present {

		// loading model by id
		//---------------------
		err = modelAsStorable.Load(utils.InterfaceToString(modelID))
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		// filling model attributes
		//--------------------------
		for attribute, value := range itemData {
			if attribute == it.idKey {
				continue
			}

			if useAttribute, wasMentioned := it.attributes[attribute]; !wasMentioned || useAttribute {
				modelAsObject.Set(attribute, value)
			}
		}

		// storing model
		//---------------
		err = modelAsStorable.Save()
		if err != nil {
			return nil, err
		}
	}

	return cmdModel, nil
}

// Init is a DELETE command initialization routines
func (it *ImportCmdDelete) Init(args []string, exchange map[string]interface{}) error {
	workingModel, err := ArgsFindWorkingModel(args, []string{"InterfaceStorable"})
	if err != nil {
		return env.ErrorDispatch(err)
	}

	it.model = workingModel
	it.idKey = ArgsFindIDKey(args)

	if it.model == nil {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "049c8839f2414b9bb7369338552b8143", "DELETE command have no assigned model to work on")
	}

	if it.idKey == "" {
		it.idKey = "_id"
	}

	return nil
}

// Process is a DELETE command processor
func (it *ImportCmdDelete) Process(itemData map[string]interface{}, input interface{}, exchange map[string]interface{}) (interface{}, error) {
	// preparing model
	//-----------------
	cmdModel, err := it.model.New()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	modelAsStorable := cmdModel.(models.InterfaceStorable)

	if modelID, present := itemData[it.idKey]; present {

		// setting id to model
		//---------------------
		err = modelAsStorable.SetID(utils.InterfaceToString(modelID))
		if err != nil {
			return nil, env.ErrorDispatch(err)
		}

		// deleting model
		//----------------
		err = modelAsStorable.Delete()
		if err != nil {
			return nil, err
		}
	}

	return cmdModel, nil
}

// Init is a STORE command initialization routines
func (it *ImportCmdStore) Init(args []string, exchange map[string]interface{}) error {
	namedArgs := ArgsGetAsNamed(args, false)
	if len(args) > 1 && len(namedArgs) != len(args)-1 {
		it.storeObjectAs = args[1]
	}

	for argName, argValue := range namedArgs {
		if strings.HasPrefix(argValue, "-") {
			switch strings.TrimPrefix(argName, "-") {
			case "prefix":
				it.prefix = argValue
			case "prefixKey":
				it.prefixKey = argValue
			}
			continue
		}

		it.storeValueAs[argValue] = argName
	}

	return nil
}

// Process is a STORE command processor
func (it *ImportCmdStore) Process(itemData map[string]interface{}, input interface{}, exchange map[string]interface{}) (interface{}, error) {
	if it.storeObjectAs != "" {
		exchange[it.storeObjectAs] = input
	}

	prefix := ""
	if it.prefix != "" {
		prefix = it.prefix
	}

	if it.prefixKey != "" {
		if _, present := itemData[it.prefixKey]; present {
			prefix = utils.InterfaceToString(itemData[it.prefixKey])
		}
	}

	for itemKey, storeAs := range it.storeValueAs {
		if _, present := itemData[itemKey]; present {
			exchange[prefix+storeAs] = itemData[itemKey]
		}
	}

	return input, nil
}

// Init is a MEDIA command initialization routines
func (it *ImportCmdMedia) Init(args []string, exchange map[string]interface{}) error {

	if len(args) > 1 {
		it.mediaField = args[1]
	}

	if len(args) > 2 {
		it.mediaType = args[2]
	}

	if len(args) > 3 {
		it.mediaName = args[3]
	}

	if it.mediaField == "" {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "34994aea87a04d6f935ee875d0708e65", "media field was not specified")
	}

	return nil
}

// Process is a MEDIA command processor
func (it *ImportCmdMedia) Process(itemData map[string]interface{}, input interface{}, exchange map[string]interface{}) (interface{}, error) {
	inputAsMedia, ok := input.(models.InterfaceMedia)
	if !ok {
		return nil, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "f0f9b0ae2852411996eee9712e448419", "object not implements InterfaceMedia interface")
	}

	// checking for media field in itemData
	if value, present := itemData[it.mediaField]; present {
		var mediaArray []string

		// checking media field type and making it uniform
		switch typedValue := value.(type) {
		case string:
			mediaArray = append(mediaArray, typedValue)
		case []string:
			mediaArray = typedValue
		case []interface{}:
			for _, value := range typedValue {
				mediaArray = append(mediaArray, utils.InterfaceToString(value))
			}
		default:
			mediaArray = append(mediaArray, utils.InterfaceToString(typedValue))
		}

		// adding found media value(s)
		for _, mediaValue := range mediaArray {
			mediaContents := []byte{}
			var err error

			// looking for media type
			mediaType := it.mediaType
			if nameValue, present := itemData[it.mediaType]; present {
				mediaType = utils.InterfaceToString(nameValue)
			}

			// looking for media name
			mediaName := it.mediaName
			if nameValue, present := itemData[it.mediaName]; present {
				mediaName = utils.InterfaceToString(nameValue)
			}

			// checking value type
			if strings.HasPrefix(mediaValue, "http") { // we have http link
				response, err := http.Get(mediaValue)
				if err != nil {
					return input, env.ErrorDispatch(err)
				}

				if response.StatusCode != 200 {
					return input, env.ErrorNew(ConstErrorModule, ConstErrorLevel, "8fe36863b82f479bba9673a5e4008f75", "can't get image "+mediaValue+" (Status: "+response.Status+")")
				}

				// updating media type if wasn't set
				if contentType := response.Header.Get("Content-Type"); mediaType == "" && contentType != "" {
					if value := strings.Split(contentType, "/"); len(value) == 2 {
						mediaType = value[0]
					}
				}

				// updating media name if wasn't set
				if mediaName == "" {
					mediaName = path.Base(response.Request.URL.Path)
				}

				// receiving media contents
				mediaContents, err = ioutil.ReadAll(response.Body)
				if err != nil {
					return input, env.ErrorDispatch(err)
				}
			} else { // we have regular file

				// updating media name if wasn't set
				if mediaName == "" {
					mediaName = path.Base(mediaValue)
				}

				// receiving media contents
				mediaContents, err = ioutil.ReadFile(mediaValue)
				if err != nil {
					return input, env.ErrorDispatch(err)
				}
			}

			// checking if media type and name still not set
			if mediaType == "" && mediaName != "" {
				for _, imageExt := range []string{".jpg", ".jpeg", ".png", ".gif", ".svg", ".ico", ".bmp", ".tif", ".tiff"} {
					if strings.Contains(mediaName, imageExt) {
						mediaType = "image"
						break
					}
				}
				if mediaType == "" {
					for _, imageExt := range []string{".txt", ".rtf", ".pdf", ".doc", "docx", ".xls", ".xlsx", ".ppt", ".pptx"} {
						if strings.Contains(mediaName, imageExt) {
							mediaType = "document"
							break
						}
					}
				}
			}

			if mediaType == "" {
				mediaType = "unknown"
			}

			if mediaName == "" {
				mediaName = "media"

				if object, ok := inputAsMedia.(models.InterfaceObject); ok {
					if objectID := utils.InterfaceToString(object.Get("_id")); objectID != "" {
						mediaName += "_" + objectID
					}
				}
			}

			// finally adding media to object
			err = inputAsMedia.AddMedia(mediaType, mediaName, mediaContents)
			if err != nil {
				return input, env.ErrorDispatch(err)
			}
		}
	}

	return input, nil
}

// Init is a ATTRIBUTE_ADD command initialization routines
func (it *ImportCmdAttributeAdd) Init(args []string, exchange map[string]interface{}) error {

	workingModel, err := ArgsFindWorkingModel(args, []string{"InterfaceCustomAttributes"})
	if err != nil {
		return env.ErrorDispatch(err)
	}
	modelAsCustomAttributesInterface := workingModel.(models.InterfaceCustomAttributes)

	attributeName := ""

	namedArgs := ArgsGetAsNamed(args, true)
	for _, checkingKey := range []string{"attribute", "attr", "2"} {
		if argValue, present := namedArgs[checkingKey]; present {
			attributeName = argValue
			break
		}
	}

	if attributeName == "" {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "86b28a7b10a641ec9d59ccd55bb63632", "attribute name was not specified, untill impex attribute add")
	}

	attribute := models.StructAttributeInfo{
		Model:      workingModel.GetModelName(),
		Collection: modelAsCustomAttributesInterface.GetCustomAttributeCollectionName(),
		Attribute:  attributeName,
		Type:       "text",
		IsRequired: false,
		IsStatic:   false,
		Label:      strings.Title(attributeName),
		Group:      "General",
		Editors:    "text",
		Options:    "",
		Default:    "",
		Validators: "",
		IsLayered:  false,
	}

	for key, value := range namedArgs {
		switch strings.ToLower(key) {
		case "type":
			attribute.Type = utils.InterfaceToString(value)
		case "label":
			attribute.Label = utils.InterfaceToString(value)
		case "group":
			attribute.Group = utils.InterfaceToString(value)
		case "editors":
			attribute.Editors = utils.InterfaceToString(value)
		case "options":
			attribute.Options = utils.InterfaceToString(value)
		case "default":
			attribute.Default = utils.InterfaceToString(value)
		case "validators":
			attribute.Validators = utils.InterfaceToString(value)
		case "isrequired", "required":
			attribute.IsRequired = utils.InterfaceToBool(value)
		case "islayered", "layered":
			attribute.IsLayered = utils.InterfaceToBool(value)
		}
	}

	it.model = workingModel
	it.attribute = attribute

	return nil
}

// Process is a ATTRIBUTE_ADD command processor
func (it *ImportCmdAttributeAdd) Process(itemData map[string]interface{}, input interface{}, exchange map[string]interface{}) (interface{}, error) {
	modelAsCustomAttributesInterface := it.model.(models.InterfaceCustomAttributes)
	err := modelAsCustomAttributesInterface.AddNewAttribute(it.attribute)
	if err != nil {
		env.ErrorDispatch(err)
	}

	return input, nil
}
