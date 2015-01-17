package block

import (
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// GetIdentifier returns cms block identifier
func (it *DefaultCMSBlock) GetIdentifier() string {
	return it.Identifier
}

// SetIdentifier sets csm block identifier value
func (it *DefaultCMSBlock) SetIdentifier(newValue string) error {
	it.Identifier = newValue
	return nil
}

// GetContent returns cms block content
func (it *DefaultCMSBlock) GetContent() string {
	return it.Content
}

// SetContent sets cms block content value
func (it *DefaultCMSBlock) SetContent(newValue string) error {
	it.Content = newValue
	return nil
}

// LoadByIdentifier loads data of CMSBlock by its identifier
func (it *DefaultCMSBlock) LoadByIdentifier(identifier string) error {
	collection, err := db.GetCollection(ConstCmsBlockCollectionName)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	collection.AddFilter("identifier", "=", identifier)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	records, err := collection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if len(records) == 0 {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "8890c17e-56cb-4a54-b37c-9ee787e15067", "not found")
	}
	record := records[0]

	it.SetID(utils.InterfaceToString(record["_id"]))

	it.Content = utils.InterfaceToString(record["content"])
	it.Identifier = utils.InterfaceToString(record["identifier"])
	it.CreatedAt = utils.InterfaceToTime(record["created_at"])
	it.UpdatedAt = utils.InterfaceToTime(record["updated_at"])

	return nil
}

// EvaluateContent applying GO text template to content value
func (it *DefaultCMSBlock) EvaluateContent() string {
	evaluatedContent, err := utils.TextTemplate(it.GetContent(), it.ToHashMap())
	if err == nil {
		return evaluatedContent
	}

	return it.GetContent()
}
