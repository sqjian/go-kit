package id_test

import (
	"github.com/sqjian/toolkit/id"
	"log"
	"testing"
)

func TestGenerateUuidV1(t *testing.T) {
	log.Println(id.GenerateUuidV1())
	log.Println(id.GenerateUuidV1())
	log.Println(id.GenerateUuidV1())
}
