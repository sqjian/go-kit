package url

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"net/url"
)

func IsValidUrl(u1 string) error {
	_, err := url.ParseRequestURI(u1)
	if err != nil {
		return err
	}

	u, err := url.Parse(u1)
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
