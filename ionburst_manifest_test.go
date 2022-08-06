package ionburst

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func TestManifestData(t *testing.T) {
	cli, err := NewClientPathAndProfile("", "dhcp-sucks", true)
	if err != nil {
		t.Error(err)
		return
	}

	name := "go-sdk-manifest"
	name2 := "go-sdk-manifest2"
	file := "/tmp/STScI-01G7DB1FHPMJCCY59CQGZC1YJQ.png"

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

	rdr, _ := os.Open(file)

	err = cli.PutManifest(name2, rdr, "")
	if err != nil {
		t.Error(err)
		return
	} else {
		fmt.Printf("Uploaded: %s\n", name2)
	}

	err = cli.Head(name2)
	if err != nil {
		t.Error(err)
	} else {
		fmt.Printf("Checked: %s\n", name2)
	}

	size, err = cli.HeadWithLen(name2)
	if err != nil {
		t.Error(err)
	} else {
		fmt.Printf("Size: %d\n", size)
	}

	manifest, err = cli.Get(name2)
	if err != nil {
		t.Error(err)
		return
	} else {
		buf := new(bytes.Buffer)
		buf.ReadFrom(manifest)
		fmt.Println(buf.String())
	}

	err = cli.DeleteManifest(name2)
	if err != nil {
		t.Error(err)
		return
	} else {
		fmt.Println("Manifest deleted.")
	}

}
