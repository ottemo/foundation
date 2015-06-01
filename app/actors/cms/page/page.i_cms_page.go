package page

import (
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// GetEnabled returns page enabled flag
func (it *DefaultCMSPage) GetEnabled() bool {
	return it.Enabled
}

// SetEnabled returns page enabled flag
func (it *DefaultCMSPage) SetEnabled(newValue bool) error {
	it.Enabled = newValue
	return nil
}

// GetIdentifier returns page identifier
func (it *DefaultCMSPage) GetIdentifier() string {
	return it.Identifier
}

// SetIdentifier sets page identifier value
func (it *DefaultCMSPage) SetIdentifier(newValue string) error {
	it.Identifier = newValue
	return nil
}

// GetTitle returns page title
func (it *DefaultCMSPage) GetTitle() string {
	return it.Title
}

// SetTitle sets page title value
func (it *DefaultCMSPage) SetTitle(newValue string) error {
	it.Title = newValue
	return nil
}

// GetContent returns page content
func (it *DefaultCMSPage) GetContent() string {
	return it.Content
}

// SetContent sets page content value
func (it *DefaultCMSPage) SetContent(newValue string) error {
	it.Content = newValue
	return nil
}

// LoadByIdentifier loads data of CMSBlock by its identifier
func (it *DefaultCMSPage) LoadByIdentifier(identifier string) error {
	collection, err := db.GetCollection(ConstCmsPageCollectionName)
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

	it.Identifier = utils.InterfaceToString(record["identifier"])
	it.Enabled = utils.InterfaceToBool(record["enabled"])

	it.Title = utils.InterfaceToString(record["title"])
	it.Content = utils.InterfaceToString(record["content"])

	it.CreatedAt = utils.InterfaceToTime(record["created_at"])
	it.UpdatedAt = utils.InterfaceToTime(record["updated_at"])

	return nil
}

// EvaluateContent applying GO text template to content value
func (it *DefaultCMSPage) EvaluateContent() string {
	evaluatedContent, err := utils.TextTemplate(it.GetContent(), it.ToHashMap())
	if err == nil {
		return evaluatedContent
	}

	return it.GetContent()
}
