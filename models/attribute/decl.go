package attribute

import "github.com/ottemo/foundation/models"

const (
	CUSTOM_ATTRIBUTES_COLLECTION = "custom_attributes"
)

var globalCustomAttributes = map[string]map[string]models.T_AttributeInfo{}

type CustomAttributes struct {
	model      string
	attributes map[string]models.T_AttributeInfo

	values map[string]interface{}
}
