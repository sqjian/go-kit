package uuid_test

import (
	"github.com/sqjian/toolkit/uuid"
	"log"
	"testing"
)

func TestGenerateUuidV1(t *testing.T) {
	log.Println(uuid.GenerateUuidV1())
	log.Println(uuid.GenerateUuidV1())
	log.Println(uuid.GenerateUuidV1())
}
