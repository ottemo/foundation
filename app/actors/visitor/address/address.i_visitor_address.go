package address

func (it *DefaultVisitorAddress) GetVisitorId() string { return it.visitor_id }

func (it *DefaultVisitorAddress) GetFirstName() string { return it.FirstName }
func (it *DefaultVisitorAddress) GetLastName() string  { return it.LastName }

func (it *DefaultVisitorAddress) GetCompany() string { return it.Company }

func (it *DefaultVisitorAddress) GetCountry() string { return it.Country }
func (it *DefaultVisitorAddress) GetState() string   { return it.State }
func (it *DefaultVisitorAddress) GetCity() string    { return it.City }

func (it *DefaultVisitorAddress) GetAddress() string      { return it.AddressLine1 + " " + it.AddressLine2 }
func (it *DefaultVisitorAddress) GetAddressLine1() string { return it.AddressLine1 }
func (it *DefaultVisitorAddress) GetAddressLine2() string { return it.AddressLine2 }

func (it *DefaultVisitorAddress) GetPhone() string   { return it.Phone }
func (it *DefaultVisitorAddress) GetZipCode() string { return it.ZipCode }
