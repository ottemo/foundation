package page

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
