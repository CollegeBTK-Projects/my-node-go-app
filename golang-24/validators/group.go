package validators

import (
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var groupReg = regexp.MustCompile(`^ББМО-\d{2}-\d{2}$`)

func Register() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("group", func(fl validator.FieldLevel) bool {
			return groupReg.MatchString(fl.Field().String())
		})
	}
}
