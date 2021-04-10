package main

import (
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
		return ValidateErr
	}
}
