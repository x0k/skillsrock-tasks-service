package validator_adapter

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func ValidateStruct(data any) error {
	err := validate.Struct(data)
	if err != nil {
		sb := strings.Builder{}
		for i, err := range err.(validator.ValidationErrors) {
			if i > 0 {
				sb.WriteString(" and ")
			}
			sb.WriteByte('[')
			sb.WriteString(err.Field())
			sb.WriteString("]: '")
			fmt.Fprint(&sb, err.Value())
			sb.WriteString("' | Needs to implement '")
			sb.WriteString(err.Tag())
			sb.WriteByte('\'')
		}
		return errors.New(sb.String())
	}
	return nil
}
