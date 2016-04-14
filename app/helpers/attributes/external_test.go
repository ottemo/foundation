package attributes

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/utils"
	"errors"
	"fmt"
)

// -----------------------
// SampleModel declaration
// -----------------------

// SampleModel type implements InterfaceModel only interface, so it does not
// have any attribute at the beginning, it even don't have InterfaceObject and
// InterfaceStorable methods implementation, however they are available via
// embed ModelExternalAttributes type
type SampleModel struct {
	i string
	l string
	d bool
	s bool

	*ModelExternalAttributes
}

// GetModelName returns model name
func (it *SampleModel) GetModelName() string {
	return "Example"
}

// GetImplementationName returns model implementation name
func (it *SampleModel) GetImplementationName() string {
	return "ExampleObject"
}

// New constructor for a new instance
func (it *SampleModel) New() (models.InterfaceModel, error) {
	var err error

	newInstance := new(SampleModel)
	newInstance.ModelExternalAttributes, err = ExternalAttributes(newInstance)

	return newInstance, err
}

// --------------------
// Delegate declaration
// --------------------

// SampleDelegate type implements InterfaceAttributesDelegate and have handles
// on InterfaceStorable methods which should have call-back on model method call
// in order to test it we are pushing the callback status to model instance
type SampleDelegate struct {
	instance interface{}
	a string
	b float64
}

// New instantiates delegate
func (it *SampleDelegate) New(instance interface{}) (models.InterfaceAttributesDelegate, error) {
	return &SampleDelegate{instance: instance}, nil
}

// Get is a getter for external attributes
func (it *SampleDelegate) Get(attribute string) interface{} {
	switch attribute {
	case "a":
		return it.a
	case "b":
		return it.b
	}
	return nil
}

// Set is a setter for external attributes
func (it *SampleDelegate) Set(attribute string, value interface{}) error {
	switch attribute {
	case "a":
		it.a = utils.InterfaceToString(value)
	case "b":
		it.b = utils.InterfaceToFloat64(value)
	}
	return nil
}

// Load is a modelInstance.Load() method handler for external attributes
func (it *SampleDelegate) Load(id string) error {
	it.instance.(*SampleModel).l = id
	return nil
}

// Delete is a modelInstance.Delete() method handler for external attributes
func (it *SampleDelegate) Delete() error {
	it.instance.(*SampleModel).d = true
	return nil
}

// Save is a modelInstance.Save() method handler for external attributes
func (it *SampleDelegate) Save() error {
	it.instance.(*SampleModel).s = true
	return nil
}

// SetID is a modelInstance.SetID() method handler for external attributes
func (it *SampleDelegate) SetID(newID string) error {
	it.instance.(*SampleModel).i = newID
	return nil
}

// GetAttributesInfo is a specification of external attributes
func (it *SampleDelegate) GetAttributesInfo() []models.StructAttributeInfo {
	return []models.StructAttributeInfo{
		models.StructAttributeInfo{
			Model:      "",
			Collection: "",
			Attribute:  "a",
			Type:       utils.ConstDataTypeText,
			Label:      "A",
			IsRequired: false,
			IsStatic:   false,
			Group:      "Sample",
			Editors:    "text",
		},
		models.StructAttributeInfo{
			Model:      "Example",
			Collection: "",
			Attribute:  "b",
			Type:       utils.ConstDataTypeFloat,
			Label:      "B",
			IsRequired: false,
			IsStatic:   false,
			Group:      "Sample",
			Editors:    "text",
		},
	}
}

// ExampleExternalAttributes creates 2 instances of SampleModel
func ExampleExternalAttributes() {
	// registering SampleDelegate for SampleModel on attributes "a" and "b"
	modelInstance, err := new(SampleModel).New()
	if err != nil {
		panic(err)
	}

	modelEA, ok := modelInstance.(models.InterfaceExternalAttributes)
	if !ok {
		panic(errors.New("InterfaceExternalAttributes not impelemented"))
	}

	delegate := new(SampleDelegate)
	modelEA.AddExternalAttributes(delegate)

	// testing result: setting "a", "b" attributes to SampleModel instances and getting them back
	var obj1, obj2 models.InterfaceObject
	if x, err := modelInstance.New(); err == nil {
		if obj1, ok = x.(models.InterfaceObject); !ok {
			panic(errors.New("InterfaceObject not impelemented"))
		}
	} else {
		panic(err)
	}

	if x, err := modelInstance.New(); err == nil {
		if obj2, ok = x.(models.InterfaceObject); !ok {
			panic(errors.New("InterfaceObject not impelemented"))
		}
	} else {
		panic(err)
	}


	if err = obj1.Set("a", "object1"); err != nil {
		panic(err)
	}
	if err = obj2.Set("a", "object2"); err != nil {
		panic(err)
	}
	if err = obj1.Set("b", 1.2); err != nil {
		panic(err)
	}
	if err = obj2.Set("b", 3.3); err != nil {
		panic(err)
	}

	if obj1.Get("a") != "object1" || obj1.Get("b") != 1.2 ||
	   obj2.Get("a") != "object2" || obj2.Get("b") != 3.3 {
		panic(errors.New(fmt.Sprint("incorrect get values: " +
			"obj1.a=", obj1.Get("a"), ", ",
			"obj1.b=", obj1.Get("b"), ", ",
			"obj2.a=", obj2.Get("a"), ", ",
			"obj2.b=", obj2.Get("b"),
		)))
	}

	if obj1, ok := obj1.(models.InterfaceStorable); ok {
		obj1.Load("1")
		obj1.Save()
		obj1.Delete()
		obj1.SetID("10")
	} else {
		panic(errors.New("models.InterfaceStorable not implemented"))
	}

	if obj1, ok := obj1.(*SampleModel); ok {
		if !obj1.d || !obj1.s || obj1.l != "1" || obj1.i != "10" {
			panic(errors.New(fmt.Sprint("incorrect get values: " +
				"obj1.l=", obj1.l, ", ",
				"obj1.s=", obj1.s, ", ",
				"obj1.d=", obj1.d, ", ",
				"obj1.i=", obj1.i,
			)))
		}
	} else {
		panic(errors.New("(*SampleModel) conversion error"))
	}
}