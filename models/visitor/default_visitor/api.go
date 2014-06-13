package default_visitor

import (
	"errors"
	"net/http"

	"github.com/ottemo/foundation/models"
	"github.com/ottemo/foundation/models/visitor"
)

func jsonError(err error) map[string]interface{} {
	return map[string]interface{} { "error": err.Error() }
}


// sample: http://127.0.0.1:9000/visitor/load?id=5
func (it *DefaultVisitor) LoadVisitorAPI(req *http.Request) map[string]interface{} {
	queryParams := req.URL.Query()

	productId := queryParams.Get("id")
	if productId == "" {
		return jsonError( errors.New("visitor 'id' was not specified") )
	}

	if model, err := models.GetModel("Visitor"); err == nil {
		if model, ok := model.(visitor.I_Visitor); ok {

			err = model.Load( productId )
			if err != nil { return jsonError(err) }

			return model.ToHashMap()
		}
	}

	return jsonError( errors.New("Something went wrong...") )
}


// usage: http://127.0.0.1:9000/visitor/update?id=10&email=bad_guy@gmail.com&first_name=Bad&last_name=Guy
func (it *DefaultVisitor) UpdateVisitorAPI(req *http.Request) map[string]interface{} {
	queryParams := req.URL.Query()

	if queryParams.Get("id") == "" {
		return jsonError( errors.New("visitor 'id' was not specified") )
	}

	if model, err := models.GetModel("Visitor"); err == nil {
		if model, ok := model.(visitor.I_Visitor); ok {

			err := model.Load( queryParams.Get("id") )
			if err != nil { return jsonError( err ) }

			for attribute, value := range queryParams {
				err := model.Set(attribute, value[0])
				if err != nil { return jsonError(err) }
			}

			err = model.Save()
			if err != nil { return jsonError(err) }

			return model.ToHashMap()
		}
	}

	return jsonError( errors.New("Something went wrong...") )
}

// usage: http://127.0.0.1:9000/visitor/create?email=bad_guy@gmail.com&first_name=Bad&last_name=Guy
func (it *DefaultVisitor) CreateVisitorAPI(req *http.Request) map[string]interface{} {
	queryParams := req.URL.Query()

	if queryParams.Get("email") == "" {
		return jsonError( errors.New("'email' was not specified") )
	}

	if model, err := models.GetModel("Visitor"); err == nil {
		if model, ok := model.(visitor.I_Visitor); ok {

			for attribute, value := range queryParams {
				err := model.Set(attribute, value[0])
				if err != nil { return jsonError(err) }
			}

			err := model.Save()
			if err != nil { return jsonError(err) }

			return model.ToHashMap()
		}
	}

	return jsonError( errors.New("Something went wrong...") )
}
