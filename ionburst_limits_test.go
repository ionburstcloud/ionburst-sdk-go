package ionburst

import (
	"fmt"
	"strconv"
	"testing"
)

func TestLimits(t *testing.T) {
	cli, err := NewClient()
	if err != nil {
		t.Error(err)
		return
	}

	data, err := cli.GetDataUploadLimit()
	if err != nil {
		t.Error(err)
		return
	} else {
		fmt.Println("Configured data limit: " + strconv.Itoa(data))
	}

	secrets, err := cli.GetSecretsUploadLimit()
	if err != nil {
		t.Error(err)
		return
	} else {
		fmt.Println("Configured secrets limit: " + strconv.Itoa(secrets))
	}

}
