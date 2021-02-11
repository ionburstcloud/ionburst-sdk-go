package ionburst

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"
)

func makeRandomPayload(size int) (string, []byte) {
	token := make([]byte, size)
	rand.Read(token)
	name := randSeq(32)
	return name, token
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func TestClassifications(t *testing.T) {
	cli, err := NewClient()
	if err != nil {
		t.Error(err)
		return
	}

	a, err := cli.GetClassifications()
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Printf("Classifications: %d\n", len(a))

}

func TestPostData(t *testing.T) {
	cli, err := NewClient()
	if err != nil {
		t.Error(err)
		return
	}

	name, ba := makeRandomPayload(1024)

	r := bytes.NewReader(ba)

	err = cli.Put(name, r, "")
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Printf("Uploaded: %s\n", name)

	_, err = cli.Get(name)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Printf("Downloaded: %s\n", name)

	err = cli.Delete(name)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Printf("Deleted: %s\n", name)

}
