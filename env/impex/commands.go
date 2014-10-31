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

type ImpexImportCmdAttributeAdd struct {
	model     models.I_Model
	attribute models.T_AttributeInfo
}

type ImpexImportCmdInsert struct {
	model      models.I_Model
	attributes map[string]bool
	skipErrors bool
}

type ImpexImportCmdUpdate struct {
	model      models.I_Model
	attributes map[string]bool
	idKey      string
}

type ImpexImportCmdDelete struct {
	model models.I_Model
	idKey string
}

type ImpexImportCmdMedia struct {
	mediaField string
	mediaType  string
	mediaName  string
}

type ImpexImportCmdStore struct {
	storeObjectAs string
	storeValueAs  map[string]string

	prefix    string
	prefixKey string
}

// checks that model support I_Object and I_Storable interfaces
func CheckModelImplements(modelName string, neededInterfaces []string) (models.I_Model, error) {
	cmdModel, err := models.GetModel(modelName)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for _, interfaceName := range neededInterfaces {
		ok := true
		switch interfaceName {
		case "I_Storable":
			_, ok = cmdModel.(models.I_Storable)
		case "I_Object":
			_, ok = cmdModel.(models.I_Object)
		case "I_Listable":
			_, ok = cmdModel.(models.I_Listable)
		case "I_Collection":
			_, ok = cmdModel.(models.I_Collection)
		case "I_CustomAttributes":
			_, ok = cmdModel.(models.I_CustomAttributes)
		case "I_Media":
			_, ok = cmdModel.(models.I_Media)
		}

		if !ok {
			return nil, env.ErrorNew("model " + modelName + " not implements " + interfaceName)
		}
	}

	return cmdModel, nil
}

// TODO: make command parameters standardized parser to split required/optional parameters and get then in one function call

// function collects arguments into map, unnamed arguments will go as position index
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

// looking for model mention among command attributes
func ArgsFindWorkingModel(args []string, neededInterfaces []string) (models.I_Model, error) {
	var result models.I_Model = nil
	var err error = nil

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

// looking for _id mention among command attributes
func ArgsFindIdKey(args []string) string {
	namedArgs := ArgsGetAsNamed(args, false)
	for _, checkingKey := range []string{"idKey", "id", "_id"} {
		if argValue, present := namedArgs[checkingKey]; present {
			return argValue
		}
	}
	return ""
}

// looking for attributes inclusion/exclusion among args
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

// INSERT command initialization
func (it *ImpexImportCmdInsert) Init(args []string, exchange map[string]interface{}) error {

	workingModel, err := ArgsFindWorkingModel(args, []string{"I_Storable", "I_Object"})
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

// INSERT command processing
func (it *ImpexImportCmdInsert) Process(itemData map[string]interface{}, input interface{}, exchange map[string]interface{}) (interface{}, error) {

	// preparing model
	//-----------------
	if it.model == nil {
		return nil, env.ErrorNew("INSERT command have no assigned model to work on")
	}
	cmdModel, err := it.model.New()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	modelAsObject := cmdModel.(models.I_Object)
	modelAsStorable := cmdModel.(models.I_Storable)

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

// UPDATE command initialization
func (it *ImpexImportCmdUpdate) Init(args []string, exchange map[string]interface{}) error {
	workingModel, err := ArgsFindWorkingModel(args, []string{"I_Storable", "I_Object"})
	if err != nil {
		return env.ErrorDispatch(err)
	}

	it.model = workingModel
	it.attributes = ArgsFindWorkingAttributes(args)
	it.idKey = ArgsFindIdKey(args)

	if it.model == nil {
		return env.ErrorNew("INSERT command have no assigned model to work on")
	}

	if it.idKey == "" {
		it.idKey = "_id"
	}

	return nil
}

// UPDATE command processing
func (it *ImpexImportCmdUpdate) Process(itemData map[string]interface{}, input interface{}, exchange map[string]interface{}) (interface{}, error) {

	// preparing model
	//-----------------
	cmdModel, err := it.model.New()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	modelAsObject := cmdModel.(models.I_Object)
	modelAsStorable := cmdModel.(models.I_Storable)

	if modelId, present := itemData[it.idKey]; present {

		// loading model by id
		//---------------------
		err = modelAsStorable.Load(utils.InterfaceToString(modelId))
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

// DELETE command initialization
func (it *ImpexImportCmdDelete) Init(args []string, exchange map[string]interface{}) error {
	workingModel, err := ArgsFindWorkingModel(args, []string{"I_Storable"})
	if err != nil {
		return env.ErrorDispatch(err)
	}

	it.model = workingModel
	it.idKey = ArgsFindIdKey(args)

	if it.model == nil {
		return env.ErrorNew("DELETE command have no assigned model to work on")
	}

	if it.idKey == "" {
		it.idKey = "_id"
	}

	return nil
}

// DELETE command processing
func (it *ImpexImportCmdDelete) Process(itemData map[string]interface{}, input interface{}, exchange map[string]interface{}) (interface{}, error) {
	// preparing model
	//-----------------
	cmdModel, err := it.model.New()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	modelAsStorable := cmdModel.(models.I_Storable)

	if modelId, present := itemData[it.idKey]; present {

		// setting id to model
		//---------------------
		err = modelAsStorable.SetId(utils.InterfaceToString(modelId))
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

// STORE command initialization
func (it *ImpexImportCmdStore) Init(args []string, exchange map[string]interface{}) error {
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

// STORE command processing
func (it *ImpexImportCmdStore) Process(itemData map[string]interface{}, input interface{}, exchange map[string]interface{}) (interface{}, error) {
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

// MEDIA command initialization
func (it *ImpexImportCmdMedia) Init(args []string, exchange map[string]interface{}) error {

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
		return env.ErrorNew("media field was not specified")
	}

	return nil
}

// MEDIA command processing
func (it *ImpexImportCmdMedia) Process(itemData map[string]interface{}, input interface{}, exchange map[string]interface{}) (interface{}, error) {
	inputAsMedia, ok := input.(models.I_Media)
	if !ok {
		return nil, env.ErrorNew("object not implements I_Media interface")
	}

	// checking for media field in itemData
	if value, present := itemData[it.mediaField]; present {
		mediaArray := make([]string, 0)

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
			var mediaContents []byte = []byte{}
			var err error = nil

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
					return input, env.ErrorNew("can't get image " + mediaValue + " (Status: " + response.Status + ")")
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

				if object, ok := inputAsMedia.(models.I_Object); ok {
					if objectId := utils.InterfaceToString(object.Get("_id")); objectId != "" {
						mediaName += "_" + objectId
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

// ATTRIBUTE_ADD command initialization
func (it *ImpexImportCmdAttributeAdd) Init(args []string, exchange map[string]interface{}) error {

	workingModel, err := ArgsFindWorkingModel(args, []string{"I_CustomAttributes"})
	if err != nil {
		return env.ErrorDispatch(err)
	}
	modelAsCustomAttributesInterface := workingModel.(models.I_CustomAttributes)

	attributeName := ""

	namedArgs := ArgsGetAsNamed(args, true)
	for _, checkingKey := range []string{"attribute", "attr", "2"} {
		if argValue, present := namedArgs[checkingKey]; present {
			attributeName = argValue
			break
		}
	}

	if attributeName == "" {
		return env.ErrorNew("attribute name was not specified, untill impex attribute add")
	}

	attribute := models.T_AttributeInfo{
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

// ATTRIBUTE_ADD command processing
func (it *ImpexImportCmdAttributeAdd) Process(itemData map[string]interface{}, input interface{}, exchange map[string]interface{}) (interface{}, error) {
	modelAsCustomAttributesInterface := it.model.(models.I_CustomAttributes)
	err := modelAsCustomAttributesInterface.AddNewAttribute(it.attribute)
	if err != nil {
		env.ErrorDispatch(err)
	}

	return input, nil
}
