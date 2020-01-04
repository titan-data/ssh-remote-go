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
	"net/url"
	"os"
	"path/filepath"
	"testing"
)

func TestRegistered(t *testing.T) {
	r := remote.Get("ssh")
	assert.Equal(t, "ssh", r.Type())
}

func TestFromURL(t *testing.T) {
	r := remote.Get("ssh")
	u, _ := url.Parse("ssh://user:pass@host:8022/path")
	props, _ := r.FromURL(u, map[string]string{})
	assert.Equal(t, "user", props["username"])
	assert.Equal(t, "pass", props["password"])
	assert.Equal(t, "host", props["address"])
	assert.Equal(t, 8022, props["port"])
	assert.Equal(t, "/path", props["path"])
	assert.Nil(t, props["keyFile"])
}

func TestSimple(t *testing.T) {
	r := remote.Get("ssh")
	u, _ := url.Parse("ssh://user@host/path")
	props, _ := r.FromURL(u, map[string]string{})
	assert.Equal(t, "user", props["username"])
	assert.Nil(t, props["password"])
	assert.Equal(t, "host", props["address"])
	assert.Nil(t, props["port"])
	assert.Equal(t, "/path", props["path"])
	assert.Nil(t, props["keyFile"])
}

func TestKeyFile(t *testing.T) {
	r := remote.Get("ssh")
	u, _ := url.Parse("ssh://user@host/path")
	props, _ := r.FromURL(u, map[string]string{"keyFile": "~/.ssh/id_dsa"})
	assert.Equal(t, "~/.ssh/id_dsa", props["keyFile"])
}

func TestRelativePath(t *testing.T) {
	r := remote.Get("ssh")
	u, _ := url.Parse("ssh://user@host/~/relative/path")
	props, _ := r.FromURL(u, map[string]string{})
	assert.Equal(t, "relative/path", props["path"])
}

func TestBadScheme(t *testing.T) {
	r := remote.Get("ssh")
	u, _ := url.Parse("foo://user:pass@host:8022/path")
	_, err := r.FromURL(u, map[string]string{})
	assert.NotNil(t, err)
}

func TestBadPasswordAndKeyFile(t *testing.T) {
	r := remote.Get("ssh")
	u, _ := url.Parse("ssh://user:password@host/path")
	_, err := r.FromURL(u, map[string]string{"keyFile": "~/.ssh/id_dsa"})
	assert.NotNil(t, err)
}

func TestBadProperty(t *testing.T) {
	r := remote.Get("ssh")
	u, _ := url.Parse("ssh://user@host/path")
	_, err := r.FromURL(u, map[string]string{"foo": "bar"})
	assert.NotNil(t, err)
}

func TestBadMissingHost(t *testing.T) {
	r := remote.Get("ssh")
	u, _ := url.Parse("ssh:///path")
	_, err := r.FromURL(u, map[string]string{})
	assert.NotNil(t, err)
}

func TestBadSchemeOnly(t *testing.T) {
	r := remote.Get("ssh")
	u, _ := url.Parse("ssh")
	_, err := r.FromURL(u, map[string]string{})
	assert.NotNil(t, err)
}

func TestBadMissingUsername(t *testing.T) {
	r := remote.Get("ssh")
	u, _ := url.Parse("ssh://host/path")
	_, err := r.FromURL(u, map[string]string{})
	assert.NotNil(t, err)
}

func TestBadPort(t *testing.T) {
	r := remote.Get("ssh")
	u, _ := url.Parse("ssh://user@host:29348529384572398457932847539/path")
	_, err := r.FromURL(u, map[string]string{})
	assert.NotNil(t, err)
}

func TestBadMissingPath(t *testing.T) {
	r := remote.Get("ssh")
	u, _ := url.Parse("ssh://user@host")
	_, err := r.FromURL(u, map[string]string{})
	assert.NotNil(t, err)
}

func TestBadMissingHostWithUser(t *testing.T) {
	r := remote.Get("ssh")
	u, _ := url.Parse("ssh://user@/path")
	_, err := r.FromURL(u, map[string]string{})
	assert.NotNil(t, err)
}

func TestToURL(t *testing.T) {
	r := remote.Get("ssh")
	u, props, _ := r.ToURL(map[string]interface{}{"username": "username", "address": "host",
		"path": "/path"})
	assert.Equal(t, "ssh://username@host/path", u)
	assert.Empty(t, props)
}

func TestToPassword(t *testing.T) {
	r := remote.Get("ssh")
	u, props, _ := r.ToURL(map[string]interface{}{"username": "username", "address": "host",
		"path": "/path", "password": "pass"})
	assert.Equal(t, "ssh://username:*****@host/path", u)
	assert.Empty(t, props)
}

func TestToPort(t *testing.T) {
	r := remote.Get("ssh")
	u, props, _ := r.ToURL(map[string]interface{}{"username": "username", "address": "host",
		"path": "/path", "port": 812})
	assert.Equal(t, "ssh://username@host:812/path", u)
	assert.Empty(t, props)
}

func TestToRelativePath(t *testing.T) {
	r := remote.Get("ssh")
	u, props, _ := r.ToURL(map[string]interface{}{"username": "username", "address": "host",
		"path": "path"})
	assert.Equal(t, "ssh://username@host/~/path", u)
	assert.Empty(t, props)
}

func TestToKeyFile(t *testing.T) {
	r := remote.Get("ssh")
	u, props, _ := r.ToURL(map[string]interface{}{"username": "username", "address": "host",
		"path": "/path", "keyFile": "keyfile"})
	assert.Equal(t, "ssh://username@host/path", u)
	assert.Len(t, props, 1)
	assert.Equal(t, "keyfile", props["keyFile"])
}

func TestToPortFloat(t *testing.T) {
	p := float32(812)
	r := remote.Get("ssh")
	u, props, _ := r.ToURL(map[string]interface{}{"username": "username", "address": "host",
		"path": "/path", "port": p})
	assert.Equal(t, "ssh://username@host:812/path", u)
	assert.Empty(t, props)
}

func TestToPortDouble(t *testing.T) {
	r := remote.Get("ssh")
	u, props, _ := r.ToURL(map[string]interface{}{"username": "username", "address": "host",
		"path": "/path", "port": 812.0})
	assert.Equal(t, "ssh://username@host:812/path", u)
	assert.Empty(t, props)
}

func TestGetParameters(t *testing.T) {
	r := remote.Get("ssh")
	props, _ := r.GetParameters(map[string]interface{}{"username": "username", "address": "host",
		"path": "/path", "password": "pass"})
	assert.Empty(t, props)
}

func TestKeyFileParameters(t *testing.T) {
	r := remote.Get("ssh")
	file, err := ioutil.TempFile("", "ssh.test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	path, err := filepath.Abs(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile(path, []byte("KEY"), 0600)
	if err != nil {
		t.Fatal(err)
	}

	props, _ := r.GetParameters(map[string]interface{}{"username": "username", "address": "host",
		"path": "/path", "keyFile": path})
	assert.Nil(t, props["password"])
	assert.Equal(t, "KEY", props["key"])
}

func TestBadKeyFileParameters(t *testing.T) {
	r := remote.Get("ssh")
	file, err := ioutil.TempFile("", "ssh.test")
	if err != nil {
		t.Fatal(err)
	}
	path, err := filepath.Abs(file.Name())
	if err != nil {
		t.Fatal(err)
	}
	os.Remove(path)

	_, err = r.GetParameters(map[string]interface{}{"username": "username", "address": "host",
		"path": "/path", "keyFile": path})
	assert.NotNil(t, err)
}

func TestPasswordPrompt(t *testing.T) {
	r := remote.Get("ssh")
	readPassword = func(fd int) (bytes []byte, err error) {
		return []byte("pass"), nil
	}
	fmtPrintf = func(format string, a ...interface{}) (n int, err error) {
		return 0, nil
	}
	props, _ := r.GetParameters(map[string]interface{}{"username": "username", "address": "host",
		"path": "/path"})
	readPassword = terminal.ReadPassword
	fmtPrintf = fmt.Printf

	assert.Nil(t, props["key"])
	assert.Equal(t, "pass", props["password"])
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

	assert.NotNil(t, err)
}
