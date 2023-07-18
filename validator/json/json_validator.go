package json

import (
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"strings"
)

type Validator struct{}

func NewValidator(opts ...Option) (*Validator, error) {
	return &Validator{}, nil
}

func (v *Validator) Validate(schema []byte, data []byte) error {
	if len(data) == 0 {
		return fmt.Errorf("empty data")
	}
	if len(schema) == 0 {
		return fmt.Errorf("empty schema")
	}

	return v.validateJson(schema, data)
}
func (v *Validator) validateJson(schema []byte, data []byte) error {
	schemaInst, schemaInstErr := gojsonschema.NewSchema(gojsonschema.NewBytesLoader(schema))
	if schemaInstErr != nil {
		return schemaInstErr
	}
	return v.checkSchemaResult(schemaInst.Validate(gojsonschema.NewBytesLoader(data)))
}

func (v *Validator) checkSchemaResult(result *gojsonschema.Result, err error) error {
	if err != nil {
		return err
	}
	if !result.Valid() {
		return v.formatError(result.Errors())
	}
	return nil
}

func (v *Validator) formatError(errors []gojsonschema.ResultError) error {
	logWithNumber := make([]string, len(errors))
	for i, l := range errors {
		if l != nil {
			logWithNumber[i] = fmt.Sprintf("#%d: %s", i+1, l.String())
		}
	}
	return fmt.Errorf("err detail:\n%s", strings.Join(logWithNumber, "\n"))
}
