package uuid

import (
	"github.com/google/uuid"
)

func GenerateUuidV1() (string, error) {
	return generateUuidV1()
}
func generateUuidV1() (string, error) {
	i, e := uuid.NewUUID()
	if e != nil {
		return "", e
	}
	return i.String(), nil
}
