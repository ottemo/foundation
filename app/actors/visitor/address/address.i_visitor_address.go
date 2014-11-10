package address

// GetVisitorID returns the Visitor ID for the address as a string
func (it *DefaultVisitorAddress) GetVisitorID() string { return it.visitorID }

// GetFirstName will return the Visitor first name for the address as a string
func (it *DefaultVisitorAddress) GetFirstName() string { return it.FirstName }

// GetLastName will return the Visitor last name for the address
func (it *DefaultVisitorAddress) GetLastName() string { return it.LastName }

// GetCompany will return the Company attribute for the address as a string
func (it *DefaultVisitorAddress) GetCompany() string { return it.Company }

// GetCountry will return the Country on the current address as a string
func (it *DefaultVisitorAddress) GetCountry() string { return it.Country }

// GetState will return the State on the current address as a string
func (it *DefaultVisitorAddress) GetState() string { return it.State }

// GetCity will return the Citty on the currenty address as a string
func (it *DefaultVisitorAddress) GetCity() string { return it.City }

// GetAddress will return the Address as a string
func (it *DefaultVisitorAddress) GetAddress() string { return it.AddressLine1 + " " + it.AddressLine2 }

// GetAddressLine1 will return Line 1 of the address as a string
func (it *DefaultVisitorAddress) GetAddressLine1() string { return it.AddressLine1 }

// GetAddressLine2 will return Line 2 of the address as a string
func (it *DefaultVisitorAddress) GetAddressLine2() string { return it.AddressLine2 }

// GetPhone will return the Phone as a string
func (it *DefaultVisitorAddress) GetPhone() string { return it.Phone }

// GetZipCode will return the Zip as a string
func (it *DefaultVisitorAddress) GetZipCode() string { return it.ZipCode }
