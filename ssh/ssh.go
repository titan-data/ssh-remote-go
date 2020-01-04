/*
 * Copyright The Titan Project Contributors.
 */
package ssh

import (
	"errors"
	"fmt"
	"github.com/titan-data/remote-sdk-go/remote"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
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
	u := fmt.Sprintf("ssh://%s", properties["username"])
	if properties["password"] != nil {
		u += ":*****"
	}
	u += fmt.Sprintf("@%s", properties["address"])
	if properties["port"] != nil {
		var port = 0
		if flt, ok := properties["port"].(float32); ok {
			port = int(flt)
		} else if flt, ok := properties["port"].(float64); ok {
			port = int(flt)
		} else {
			port = properties["port"].(int)
		}
		u += fmt.Sprintf(":%d", port)
	}
	if properties["path"].(string)[0:1] != "/" {
		u += "/~/"
	}
	u += properties["path"].(string)

	retProps := map[string]string{}
	if properties["keyFile"] != nil {
		retProps["keyFile"] = properties["keyFile"].(string)
	}

	return u, retProps, nil
}

var readPassword = terminal.ReadPassword
var fmtPrintf = fmt.Printf

func (s sshRemote) GetParameters(remoteProperties map[string]interface{}) (map[string]interface{}, error) {
	result := map[string]interface{}{}

	if remoteProperties["keyFile"] != nil {
		content, err := ioutil.ReadFile(remoteProperties["keyFile"].(string))
		if err != nil {
			return nil, fmt.Errorf("failed to read key file %s: %w", remoteProperties["keyFile"], err)
		}
		result["key"] = string(content)
	}

	if remoteProperties["password"] == nil && remoteProperties["keyFile"] == nil {
		fmtPrintf("password: ")
		pw, err := readPassword(0)
		if err != nil {
			return nil, fmt.Errorf("failed to read password: %w", err)
		}
		result["password"] = string(pw)
	}

	return result, nil
}

func init() {
	remote.Register(sshRemote{})
}
