package address

// GetVisitorID returns the Visitor ID for the Visitor Address
func (it *DefaultVisitorAddress) GetVisitorID() string { return it.visitorID }

// GetFirstName returns the First Name of the Visitor Address
func (it *DefaultVisitorAddress) GetFirstName() string { return it.FirstName }

// GetLastName returns the Last Name of the Visitor Address
func (it *DefaultVisitorAddress) GetLastName() string { return it.LastName }

// GetCompany will return the Company attribute of the Visitor Address
func (it *DefaultVisitorAddress) GetCompany() string { return it.Company }

// GetCountry will return the Country attribute of the Visitor Address
func (it *DefaultVisitorAddress) GetCountry() string { return it.Country }

// GetState will return the State attribute of the Visitor Address
func (it *DefaultVisitorAddress) GetState() string { return it.State }

// GetCity will return the City attribute of the Visitor Address
func (it *DefaultVisitorAddress) GetCity() string { return it.City }

// GetAddress will return the full Address of the current Visitor Address
func (it *DefaultVisitorAddress) GetAddress() string { return it.AddressLine1 + " " + it.AddressLine2 }

// GetAddressLine1 will return the Line 1 attribute of the Visitor Address
func (it *DefaultVisitorAddress) GetAddressLine1() string { return it.AddressLine1 }

// GetAddressLine2 will return the Line 2 attribute of the Visitor Address
func (it *DefaultVisitorAddress) GetAddressLine2() string { return it.AddressLine2 }

// GetPhone will return the phone attribute of the Visitor Address
func (it *DefaultVisitorAddress) GetPhone() string { return it.Phone }

// GetZipCode will return the zip code attribute of the Visitor Address
func (it *DefaultVisitorAddress) GetZipCode() string { return it.ZipCode }
