package block

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
