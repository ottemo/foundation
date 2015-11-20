package events

const (
	ConstKGEventConfigPath = "kg.event"
	ConstKGMailSubscriptionEvent = "kg.event.subscribe_to_list"
)

type MailSubscriptionConfig struct {
	ListId string `json:"listId"`
	Sku    string `json:"sku"`
}
