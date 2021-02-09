package xml

import (
	"fmt"
	"github.com/lestrrat-go/libxml2"
	"github.com/lestrrat-go/libxml2/xsd"
)

func ValidateXmlWithXsd(xmlData, xsdData []byte) error {
	xsdInstance, xsdInstanceErr := xsd.Parse(xsdData)
	if xsdInstanceErr != nil {
		return xsdInstanceErr
	}
	defer xsdInstance.Free()

	xmlInstance, xmlInstanceErr := libxml2.Parse(xmlData)
	if xmlInstanceErr != nil {
		return xmlInstanceErr
	}
	defer xmlInstance.Free()

	if ValidateErr := xsdInstance.Validate(xmlInstance); ValidateErr != nil {
		var errDesc []string
		for _, e := range ValidateErr.(xsd.SchemaValidationError).Errors() {
			errDesc = append(errDesc, e.Error())
		}
		return fmt.Errorf("validateFailed:%v,errDesc:%v", ValidateErr, errDesc)
	}
	return nil
}
