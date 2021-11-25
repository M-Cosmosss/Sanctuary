package form

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/flamego/validator"

	"github.com/flamego/binding"
	"github.com/flamego/flamego"
)

func Bind(model interface{}) flamego.Handler {
	return binding.JSON(model, binding.Options{
		ErrorHandler: errorHandler,
	})
}

func errorHandler(c flamego.Context, errs binding.Errors) {
	c.ResponseWriter().WriteHeader(http.StatusBadRequest)
	c.ResponseWriter().Header().Set("Content-Type", "application/json")

	var msg string
	if errs[0].Category == binding.ErrorCategoryDeserialization {
		msg = "validator error"
	} else {
		if v, ok := errs[0].Err.(validator.ValidationErrors); ok {
			n := v[0].Namespace()
			switch v[0].Tag() {
			case "required":
				msg = fmt.Sprintf("%s must be required", n)
			case "len":
				msg = fmt.Sprintf("len of %s error", n)
			default:
				msg = "unknown form error"
			}
		}
	}
	body := map[string]interface{}{
		"errMsg": msg,
	}
	if err := json.NewEncoder(c.ResponseWriter()).Encode(body); err != nil {
		log.Println("validator errorHandler: json encode error")
	}
}
