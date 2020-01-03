/*
 * Copyright The Titan Project Contributors.
 */
package s3web

import (
	"errors"
	"fmt"
	"github.com/titan-data/remote-sdk-go/pkg/remote"
	"net/url"
	"reflect"
	"strings"
)

type s3webRemote struct {
}

func (s s3webRemote) Type() string {
	return "s3web"
}

func (s s3webRemote) FromURL(url url.URL, additionalProperties map[string]string) (map[string]interface{}, error) {
	if url.Scheme != "s3web" {
		return nil, errors.New("invalid remote scheme")
	}

	if url.User != nil {
		return nil, errors.New("username and password cannot be specified")
	}

	if url.Hostname() == "" {
		return nil, errors.New("missing host name in remote")
	}

	if len(additionalProperties) != 0 {
		return nil, errors.New(fmt.Sprintf("invalid property '%s'", reflect.ValueOf(additionalProperties).MapKeys()[0].String()))
	}

	u := fmt.Sprintf("http://%s%s", url.Host, url.Path)
	return map[string]interface{}{"url": u}, nil
}

func (s s3webRemote) ToURL(properties map[string]interface{}) (string, map[string]string, error) {
	u := properties["url"].(string)
	return strings.Replace(u, "http", "s3web", 1), map[string]string{}, nil
}

func (s s3webRemote) GetParameters(remoteProperties map[string]interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

func init() {
	remote.Register(s3webRemote{})
}
