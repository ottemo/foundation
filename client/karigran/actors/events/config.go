package events

import "github.com/ottemo/foundation/env"

func setupConfig() error {
	config := env.GetConfig()
	if config == nil {
		return env.ErrorDispatch("Unable configure mailchimp subscription event handler")
	}

	if err := config.RegisterItem(env.StructConfigItem{
		Path:        ConstKGEventConfigPath,
		Value:       nil,
		Type:        env.ConstConfigTypeGroup,
		Editor:      "",
		Options:     nil,
		Label:       "Kari Gran Events",
		Description: "Event handlers specific to Kari Gran",
		Image:       "",
	}, nil); err != nil {
		return env.ErrorDispatch(err)
	}

	if err := config.RegisterItem(env.StructConfigItem{
		Path:        ConstKGMailSubscriptionEvent,
		Value:       nil,
		Type:        env.ConstConfigTypeJSON,
		Editor:      "multiline_text",
		Options:     nil,
		Label:       "Mail Subscription",
		Description: `Mail subscription configuration is done in json as an array of
					  { "listId" : "mc list id", "sku":"sku trigger" }.`,
		Image:       "",
	}, nil); err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
