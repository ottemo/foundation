package default_address

func (it *DefaultVisitorAddress) GetStreet() string { return it.Street }
func (it *DefaultVisitorAddress) GetCity() string { return it.City }
func (it *DefaultVisitorAddress) GetState() string { return it.State }
func (it *DefaultVisitorAddress) GetPhone() string { return it.Phone }
func (it *DefaultVisitorAddress) GetZipCode() string { return it.ZipCode }
