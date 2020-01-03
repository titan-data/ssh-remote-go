/*
 * Copyright The Titan Project Contributors.
 */
package ssh

import (
	"github.com/stretchr/testify/assert"
	"github.com/titan-data/remote-sdk-go/remote"
	"net/url"
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

/*

   "basic SSH remote to URI succeeds" {
       val (uri, parameters) = client.toUri(mapOf("username" to "username", "address" to "host",
               "path" to "/path"))
       uri shouldBe "ssh://username@host/path"
       parameters.size shouldBe 0
   }

   "SSH remote with password to URI succeeds" {
       val (uri, parameters) = client.toUri(mapOf("username" to "username", "address" to "host",
               "path" to "/path", "password" to "pass"))
       uri shouldBe "ssh://username:*****@host/path"
       parameters.size shouldBe 0
   }

   "SSH remote with port to URI succeeds" {
       val (uri, parameters) = client.toUri(mapOf("username" to "username", "address" to "host",
               "path" to "/path", "port" to 812))
       uri shouldBe "ssh://username@host:812/path"
       parameters.size shouldBe 0
   }

   "SSH remote with relative path to URI succeeds" {
       val (uri, parameters) = client.toUri(mapOf("username" to "username", "address" to "host",
               "path" to "path"))
       uri shouldBe "ssh://username@host/~/path"
       parameters.size shouldBe 0
   }

   "SSH remote with keyfile to URI succeeds" {
       val (uri, parameters) = client.toUri(mapOf("username" to "username", "address" to "host",
               "path" to "/path", "keyFile" to "keyfile"))
       uri shouldBe "ssh://username@host/path"
       parameters.size shouldBe 1
       parameters["keyFile"] shouldBe "keyfile"
   }

   "SSH remote with port as double succeeds" {
       val (uri, parameters) = client.toUri(mapOf("username" to "username", "address" to "host",
               "path" to "/path", "port" to 812.0))
       uri shouldBe "ssh://username@host:812/path"
       parameters.size shouldBe 0
   }

   "get basic SSH get parameters succeeds" {
       val params = client.getParameters(mapOf("username" to "username", "address" to "host",
               "path" to "/path", "password" to "pass"))
       params["password"] shouldBe null
       params["key"] shouldBe null
   }

   "get SSH parameters with keyfile succeeds" {
       val keyFile = createTempFile()
       try {
           keyFile.writeText("KEY")
           val params = client.getParameters(mapOf("username" to "username", "address" to "host",
                   "path" to "/path", "keyFile" to keyFile.absolutePath))
           params["password"] shouldBe null
           params["key"] shouldBe "KEY"
       } finally {
           keyFile.delete()
       }
   }

   "prompt for SSH password succeeds" {
       every { console.readPassword(any()) } returns "pass".toCharArray()
       val params = client.getParameters(mapOf("username" to "username", "address" to "host",
               "path" to "/path"))
       params["password"] shouldBe "pass"
       params["key"] shouldBe null
   }
*/
