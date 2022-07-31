package ionburst

import (
	"fmt"
	"testing"
)

func TestManifestData(t *testing.T) {
	cli, err := NewClientPathAndProfile("", "send_dev", true)
	if err != nil {
		t.Error(err)
		return
	}

	name := "go-sdk-manifest"
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

	_, err = cli.Get(name)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Printf("Downloaded: %s\n", name)
}
