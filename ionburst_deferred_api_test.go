package ionburst

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func TestDeferredUpload(t *testing.T) {
	cli, err := NewClient()
	if err != nil {
		t.Error(err)
		return
	}

	a, err := cli.getClassifications()
	if err != nil {
		t.Error(err)
		return
	}

	//upload a random payload, get the token and use the token to check on the status (get response)

	fmt.Printf("Classifications: %d\n", len(a))

	name, ba := makeRandomPayload(1024)

	r := bytes.NewReader(ba)

	tk, err := cli.putDeferred(name, r, "")
	if err != nil {
		t.Error(err)
		return
	}

	time.Sleep(5 * time.Second)

	fmt.Printf("Uploaded Deferred: %s (%s)\n", name, tk)

	res, err := cli.checkDeferred(tk)
	if err != nil {
		t.Error(err)
		return
	}

	if !res.Success {
		t.Error(fmt.Sprintf("ERR: %s - %d", res.Message, res.Status))
		return
	} else {
		fmt.Printf("Uploaded Deferred Success: %s\n", res.ActivityToken)

		tk, err := cli.getDeferred(name)
		if err != nil {
			t.Error(err)
			return
		}

		time.Sleep(5 * time.Second)

		fmt.Printf("Download Deferred: %s (%s)\n", name, tk)

		res, err = cli.checkDeferred(tk)
		if err != nil {
			t.Error(err)
			return
		} else if !res.Success {
			t.Error(fmt.Sprintf("ERR: %s - %d", res.Message, res.Status))
			return
		} else {
			_, err := cli.fetchDeferred(tk)
			if err != nil {
				t.Error(err)
				return
			}

			err = cli.delete(name)
			if err != nil {
				t.Error(err)
				return
			}

			fmt.Printf("Deleted: %s\n", name)

		}

	}

}
