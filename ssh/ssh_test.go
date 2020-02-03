/*
 * Copyright The Titan Project Contributors.
 */
package ssh

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/titan-data/remote-sdk-go/remote"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestRegistered(t *testing.T) {
	r := remote.Get("ssh")
	ret, err := r.Type()
	if assert.NoError(t, err) {
		assert.Equal(t, "ssh", ret)
	}
}

func TestFromURL(t *testing.T) {
	r := remote.Get("ssh")
	props, err := r.FromURL("ssh://user:pass@host:8022/path", map[string]string{})
	if assert.NoError(t, err) {
		assert.Equal(t, "user", props["username"])
		assert.Equal(t, "pass", props["password"])
		assert.Equal(t, "host", props["address"])
		assert.Equal(t, 8022, props["port"])
		assert.Equal(t, "/path", props["path"])
		assert.Nil(t, props["keyFile"])
	}
}

func TestSimple(t *testing.T) {
	r := remote.Get("ssh")
	props, err := r.FromURL("ssh://user@host/path", map[string]string{})
	if assert.NoError(t, err) {
		assert.Equal(t, "user", props["username"])
		assert.Nil(t, props["password"])
		assert.Equal(t, "host", props["address"])
		assert.Nil(t, props["port"])
		assert.Equal(t, "/path", props["path"])
		assert.Nil(t, props["keyFile"])
	}
}

func TestKeyFile(t *testing.T) {
	r := remote.Get("ssh")
	props, err := r.FromURL("ssh://user@host/path", map[string]string{"keyFile": "~/.ssh/id_dsa"})
	if assert.NoError(t, err) {
		assert.Equal(t, "~/.ssh/id_dsa", props["keyFile"])
	}
}

func TestRelativePath(t *testing.T) {
	r := remote.Get("ssh")
	props, err := r.FromURL("ssh://user@host/~/relative/path", map[string]string{})
	if assert.NoError(t, err) {
		assert.Equal(t, "relative/path", props["path"])
	}
}

func TestBadScheme(t *testing.T) {
	r := remote.Get("ssh")
	_, err := r.FromURL("foo://user:pass@host:8022/path", map[string]string{})
	assert.Error(t, err)
}

func TestBadPasswordAndKeyFile(t *testing.T) {
	r := remote.Get("ssh")
	_, err := r.FromURL("ssh://user:password@host/path", map[string]string{"keyFile": "~/.ssh/id_dsa"})
	assert.Error(t, err)
}

func TestBadProperty(t *testing.T) {
	r := remote.Get("ssh")
	_, err := r.FromURL("ssh://user@host/path", map[string]string{"foo": "bar"})
	assert.Error(t, err)
}

func TestBadMissingHost(t *testing.T) {
	r := remote.Get("ssh")
	_, err := r.FromURL("ssh:///path", map[string]string{})
	assert.Error(t, err)
}

func TestBadSchemeOnly(t *testing.T) {
	r := remote.Get("ssh")
	_, err := r.FromURL("ssh", map[string]string{})
	assert.Error(t, err)
}

func TestBadMissingUsername(t *testing.T) {
	r := remote.Get("ssh")
	_, err := r.FromURL("ssh://host/path", map[string]string{})
	assert.Error(t, err)
}

func TestBadPort(t *testing.T) {
	r := remote.Get("ssh")
	_, err := r.FromURL("ssh://user@host:29348529384572398457932847539/path", map[string]string{})
	assert.Error(t, err)
}

func TestBadMissingPath(t *testing.T) {
	r := remote.Get("ssh")
	_, err := r.FromURL("ssh://user@host", map[string]string{})
	assert.Error(t, err)
}

func TestBadMissingHostWithUser(t *testing.T) {
	r := remote.Get("ssh")
	_, err := r.FromURL("ssh://user@/path", map[string]string{})
	assert.Error(t, err)
}

func TestToURL(t *testing.T) {
	r := remote.Get("ssh")
	u, props, err := r.ToURL(map[string]interface{}{"username": "username", "address": "host",
		"path": "/path"})
	if assert.NoError(t, err) {
		assert.Equal(t, "ssh://username@host/path", u)
		assert.Empty(t, props)
	}
}

func TestToPassword(t *testing.T) {
	r := remote.Get("ssh")
	u, props, err := r.ToURL(map[string]interface{}{"username": "username", "address": "host",
		"path": "/path", "password": "pass"})
	if assert.NoError(t, err) {
		assert.Equal(t, "ssh://username:*****@host/path", u)
		assert.Empty(t, props)
	}
}

func TestToPort(t *testing.T) {
	r := remote.Get("ssh")
	u, props, err := r.ToURL(map[string]interface{}{"username": "username", "address": "host",
		"path": "/path", "port": 812})
	if assert.NoError(t, err) {
		assert.Equal(t, "ssh://username@host:812/path", u)
		assert.Empty(t, props)
	}
}

func TestToRelativePath(t *testing.T) {
	r := remote.Get("ssh")
	u, props, err := r.ToURL(map[string]interface{}{"username": "username", "address": "host",
		"path": "path"})
	if assert.NoError(t, err) {
		assert.Equal(t, "ssh://username@host/~/path", u)
		assert.Empty(t, props)
	}
}

func TestToKeyFile(t *testing.T) {
	r := remote.Get("ssh")
	u, props, err := r.ToURL(map[string]interface{}{"username": "username", "address": "host",
		"path": "/path", "keyFile": "keyfile"})
	if assert.NoError(t, err) {
		assert.Equal(t, "ssh://username@host/path", u)
		assert.Len(t, props, 1)
		assert.Equal(t, "keyfile", props["keyFile"])
	}
}

func TestToPortFloat(t *testing.T) {
	p := float32(812)
	r := remote.Get("ssh")
	u, props, err := r.ToURL(map[string]interface{}{"username": "username", "address": "host",
		"path": "/path", "port": p})
	if assert.NoError(t, err) {
		assert.Equal(t, "ssh://username@host:812/path", u)
		assert.Empty(t, props)
	}
}

func TestToPortDouble(t *testing.T) {
	r := remote.Get("ssh")
	u, props, err := r.ToURL(map[string]interface{}{"username": "username", "address": "host",
		"path": "/path", "port": 812.0})
	if assert.NoError(t, err) {
		assert.Equal(t, "ssh://username@host:812/path", u)
		assert.Empty(t, props)
	}
}

func TestGetParameters(t *testing.T) {
	r := remote.Get("ssh")
	props, err := r.GetParameters(map[string]interface{}{"username": "username", "address": "host",
		"path": "/path", "password": "pass"})
	if assert.NoError(t, err) {
		assert.Empty(t, props)
	}
}

func TestKeyFileParameters(t *testing.T) {
	r := remote.Get("ssh")
	file, err := ioutil.TempFile("", "ssh.test")
	if !assert.NoError(t, err) {
		return
	}
	defer os.Remove(file.Name())
	path, err := filepath.Abs(file.Name())
	if !assert.NoError(t, err) {
		return
	}

	err = ioutil.WriteFile(path, []byte("KEY"), 0600)
	if !assert.NoError(t, err) {
		return
	}

	props, err := r.GetParameters(map[string]interface{}{"username": "username", "address": "host",
		"path": "/path", "keyFile": path})
	if assert.NoError(t, err) {
		assert.Nil(t, props["password"])
		assert.Equal(t, "KEY", props["key"])
	}
}

func TestBadKeyFileParameters(t *testing.T) {
	r := remote.Get("ssh")
	file, err := ioutil.TempFile("", "ssh.test")
	if !assert.NoError(t, err) {
		return
	}
	path, err := filepath.Abs(file.Name())
	if !assert.NoError(t, err) {
		return
	}
	err = file.Close()
	if !assert.NoError(t, err) {
		return
	}
	err = os.Remove(path)
	if assert.NoError(t, err) {
		_, err = r.GetParameters(map[string]interface{}{"username": "username", "address": "host",
			"path": "/path", "keyFile": path})
		assert.Error(t, err)
	}
}

func TestPasswordPrompt(t *testing.T) {
	r := remote.Get("ssh")
	readPassword = func(fd int) (bytes []byte, err error) {
		return []byte("pass"), nil
	}
	fmtPrintf = func(format string, a ...interface{}) (n int, err error) {
		return 0, nil
	}
	props, err := r.GetParameters(map[string]interface{}{"username": "username", "address": "host",
		"path": "/path"})
	if assert.NoError(t, err) {
		readPassword = terminal.ReadPassword
		fmtPrintf = fmt.Printf

		assert.Nil(t, props["key"])
		assert.Equal(t, "pass", props["password"])
	}
}

func TestBadPasswordPrompt(t *testing.T) {
	r := remote.Get("ssh")
	readPassword = func(fd int) (bytes []byte, err error) {
		return []byte{}, errors.New("error")
	}
	fmtPrintf = func(format string, a ...interface{}) (n int, err error) {
		return 0, nil
	}
	_, err := r.GetParameters(map[string]interface{}{"username": "username", "address": "host",
		"path": "/path"})
	readPassword = terminal.ReadPassword
	fmtPrintf = fmt.Printf

	assert.Error(t, err)
}

func TestValidateRemoteRequiredOnly(t *testing.T) {
	r := remote.Get("ssh")
	err := r.ValidateRemote(map[string]interface{}{"username": "username", "address": "host", "path": "/path"})
	assert.NoError(t, err)
}

func TestValidateRemoteAllOptional(t *testing.T) {
	r := remote.Get("ssh")
	err := r.ValidateRemote(map[string]interface{}{"username": "username", "address": "host", "path": "/path",
		"keyFile": "/keyfile", "password": "password", "port": 8022})
	assert.NoError(t, err)
}

func TestValidateRemoteBadPort(t *testing.T) {
	r := remote.Get("ssh")
	err := r.ValidateRemote(map[string]interface{}{"username": "username", "address": "host", "path": "/path",
		"keyFile": "/keyfile", "password": "password", "port": "foo"})
	assert.Error(t, err)
}

func TestValidateRemoteBadPortNegative(t *testing.T) {
	r := remote.Get("ssh")
	err := r.ValidateRemote(map[string]interface{}{"username": "username", "address": "host", "path": "/path",
		"keyFile": "/keyfile", "password": "password", "port": -1})
	assert.Error(t, err)
}

func TestValidateRemotePortFloat(t *testing.T) {
	r := remote.Get("ssh")
	err := r.ValidateRemote(map[string]interface{}{"username": "username", "address": "host", "path": "/path",
		"keyFile": "/keyfile", "password": "password", "port": 22.0})
	assert.NoError(t, err)
}

func TestValidateRemotePortFloat32(t *testing.T) {
	r := remote.Get("ssh")
	var p float32 = 22.0
	err := r.ValidateRemote(map[string]interface{}{"username": "username", "address": "host", "path": "/path",
		"keyFile": "/keyfile", "password": "password", "port": p})
	assert.NoError(t, err)
}

func TestValidateRemoteMissingRequired(t *testing.T) {
	r := remote.Get("ssh")
	err := r.ValidateRemote(map[string]interface{}{"username": "username", "address": "host"})
	assert.Error(t, err)
}

func TestValidateRemoteExtraProperty(t *testing.T) {
	r := remote.Get("ssh")
	err := r.ValidateRemote(map[string]interface{}{"username": "username", "address": "host", "path": "/path",
		"foo": "bar"})
	assert.Error(t, err)
}

func TestValidateParametersEmpty(t *testing.T) {
	r := remote.Get("ssh")
	err := r.ValidateParameters(map[string]interface{}{})
	assert.NoError(t, err)
}

func TestValidateParametersAllOptional(t *testing.T) {
	r := remote.Get("ssh")
	err := r.ValidateParameters(map[string]interface{}{"key": "key", "password": "password"})
	assert.NoError(t, err)
}

func TestValidateParametersUnknown(t *testing.T) {
	r := remote.Get("ssh")
	err := r.ValidateParameters(map[string]interface{}{"foo": "bar"})
	assert.Error(t, err)
}
