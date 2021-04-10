package main

import (
	"fmt"
	"github.com/lestrrat-go/libxml2"
	"github.com/lestrrat-go/libxml2/xsd"
	"io/ioutil"
	"log"
	"os"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	xsdFile := "test.xsd"
	f, err := os.Open(xsdFile)
	check(err)
	defer f.Close()

	buf, err := ioutil.ReadAll(f)
	check(err)

	s, err := xsd.Parse(buf)
	check(err)
	defer s.Free()

	xmlFile := "test.xml"
	f2, err := os.Open(xmlFile)
	check(err)
	defer f2.Close()

	buf2, err := ioutil.ReadAll(f2)
	check(err)

	d, err := libxml2.Parse(buf2)
	check(err)

	if err := s.Validate(d); err != nil {
		for _, e := range err.(xsd.SchemaValidationError).Errors() {
			log.Printf("error: %s", e.Error())
		}
		return
	}

	fmt.Println("Validation of", xmlFile, "against", xsdFile, "successful!")
}
