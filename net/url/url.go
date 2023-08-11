package url

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"net/url"
)

func IsValidUrl(target string) error {
	_, err := url.ParseRequestURI(target)
	if err != nil {
		return err
	}

	u, err := url.Parse(target)
	if err != nil {
		return err
	}
	if u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("illegal schema:%v or host:%v", spew.Sdump(u.Scheme), spew.Sdump(u.Host))
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("schema:%v must be one of http or https", spew.Sdump(u.Scheme))
	}

	return nil
}
func CheckUrl[T []byte | string](target T) error {
	return IsValidUrl(string(target))
}
