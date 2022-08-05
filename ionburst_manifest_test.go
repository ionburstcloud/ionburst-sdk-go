package ionburst

import (
	"bytes"
	"fmt"
	"testing"
)

func TestManifestData(t *testing.T) {
	cli, err := NewClientPathAndProfile("", "dhcp-sucks", true)
	if err != nil {
		t.Error(err)
		return
	}

	name := "go-sdk-manifest"
	file := "/tmp/STScI-01G7ETPF7DVBJAC42JR5N6EQRH.png"

	err = cli.PutManifestFromFile(name, file, "")
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Printf("Uploaded: %s\n", name)

	err = cli.Head(name)
	if err != nil {
		t.Error(err)
	} else {
		fmt.Printf("Checked: %s\n", name)
	}

	size, err := cli.HeadWithLen(name)
	if err != nil {
		t.Error(err)
	} else {
		fmt.Printf("Size: %d\n", size)
	}

	manifest, err := cli.Get(name)
	if err != nil {
		t.Error(err)
		return
	} else {
		buf := new(bytes.Buffer)
		buf.ReadFrom(manifest)
		fmt.Println(buf.String())
	}

	err = cli.DeleteManifest(name)
	if err != nil {
		t.Error(err)
		return
	} else {
		fmt.Println("Manifest deleted.")
	}
}
