/*
 * Copyright The Titan Project Contributors.
 */
package ssh

import (
	"errors"
	"fmt"
	"github.com/titan-data/remote-sdk-go/remote"
	"net/url"
	"strconv"
	"strings"
)

type sshRemote struct {
}

func (s sshRemote) Type() string {
	return "ssh"
}

func (s sshRemote) FromURL(url *url.URL, additionalProperties map[string]string) (map[string]interface{}, error) {
	if url.Scheme != "ssh" {
		return nil, errors.New("invalid remote scheme")
	}

	if url.Path == "" {
		return nil, errors.New("missing remote path")
	}

	if url.Hostname() == "" {
		return nil, errors.New("missing remote host")
	}

	if url.User == nil || url.User.Username() == "" {
		return nil, errors.New("missing remote username")
	}

	path := url.Path
	if strings.Index(path, "/~/") == 0 {
		path = path[3:]
	}

	keyFile := additionalProperties["keyFile"]
	password, passwordSet := url.User.Password()
	if keyFile != "" && passwordSet {
		return nil, errors.New("both remote password and key file cannot be specified")
	}

	for k := range additionalProperties {
		if k != "keyFile" {
			return nil, errors.New(fmt.Sprintf("invalid rmeote property '%s'", k))
		}
	}

	result := map[string]interface{}{
		"username": url.User.Username(),
		"address":  url.Hostname(),
		"path":     path,
	}

	if password != "" {
		result["password"] = password
	}
	if url.Port() != "" {
		port, err := strconv.Atoi(url.Port())
		if err != nil {
			return nil, fmt.Errorf("invalid port '%s': %w", url.Port(), err)
		}
		result["port"] = port
	}
	if keyFile != "" {
		result["keyFile"] = keyFile
	}

	return result, nil
}

func (s sshRemote) ToURL(properties map[string]interface{}) (string, map[string]string, error) {
	u := properties["url"].(string)
	return strings.Replace(u, "http", "s3web", 1), map[string]string{}, nil
}

func (s sshRemote) GetParameters(remoteProperties map[string]interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

func init() {
	remote.Register(sshRemote{})
}
