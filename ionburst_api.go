package ionburst

import (
	"encoding/json"
	"io"
	"os"

	"gitlab.com/ionburst/ionburst-sdk-go/models"
)

type DeferredToken string

func (cli *Client) getClassifications() ([]string, error) {
	cli.logger.Debug("Getting Classifications")
	if res, err := cli.doGet("api/classification", nil); err != nil {
		return nil, err
	} else {
		var classifications []string
		if err := json.Unmarshal(res.Body(), &classifications); err != nil {
			return nil, err
		}
		cli.logger.Debug("Retrieved Classifications")
		return classifications, nil
	}
}

func (cli *Client) put(id string, reader io.Reader, classification string) error {
	cli.logger.Debug("Uploading Ionburst object", id)
	_, err := cli.doPostBinary("api/data/"+id, reader, classification)
	if err == nil {
		cli.logger.Debug("Ionburst Upload Complete for", id)
	}
	return err
}

func (cli *Client) putFromFile(id string, file string, classification string) error {
	cli.logger.Debug("Uploading Ionburst object", id, "from file", file)
	rdr, err := os.Open(file)
	if err != nil {
		return err
	}
	_, err = cli.doPostBinary("api/data/"+id, rdr, classification)
	if err == nil {
		cli.logger.Debug("Ionburst Upload Complete for", id)
	}
	return err
}

func (cli *Client) get(id string) (io.Reader, error) {
	cli.logger.Debug("Downloading Ionburst object", id)
	return cli.doGetBinary("api/data/"+id, nil)
}

func (cli *Client) getWithLen(id string) (io.Reader, int64, error) {
	cli.logger.Debug("Downloading Ionburst object", id)
	return cli.doGetBinaryLen("api/data/"+id, nil)
}

func (cli *Client) getToFile(id string, file string) error {
	cli.logger.Debug("Downloading Ionburst object", id, "to file", file)
	rdr, err := cli.doGetBinary("api/data/"+id, nil)
	if err != nil {
		return err
	}
	wr, err := os.Create(file)
	if err != nil {
		return err
	}
	_, err = io.Copy(wr, rdr)
	cli.logger.Debug("Ionburst object download", id, "complete")
	return err
}

func (cli *Client) delete(id string) error {
	cli.logger.Debug("Starting Deletion of Ionburst object", id)
	_, err := cli.doDelete("api/data/"+id, nil)
	cli.logger.Debug("Deletion of Ionburst object", id, "complete")
	return err
}

func (cli *Client) putDeferred(id string, reader io.Reader, classification string) (DeferredToken, error) {
	cli.logger.Debug("Uploading Ionburst object deferred for: ", id)
	tk, err := cli.doPostBinaryDeferred("api/data/deferred/start/"+id, reader, classification)
	if err == nil {
		cli.logger.Debug("Ionburst Deferred Upload Started for: ", id, " token: ", tk)
	} else {
		return "", err
	}
	return DeferredToken(tk), nil
}

func (cli *Client) getDeferred(id string) (DeferredToken, error) {
	cli.logger.Debug("Downloading Ionburst object deferred for: ", id)
	tk, err := cli.doGet("api/data/deferred/start/"+id, nil)
	if err == nil {

		var defResp models.DeferredTokenResponse
		err = json.Unmarshal(tk.Body(), &defResp)
		if err != nil {
			return "", err
		}
		cli.logger.Debug("Ionburst Deferred Upload Started for: ", id, " token: ", defResp.DeferredToken)
		return DeferredToken(defResp.DeferredToken), nil
	} else {
		return "", err
	}
}

func (cli *Client) checkDeferred(token DeferredToken) (*models.WorkflowResult, error) {
	res, err := cli.doGet("api/data/deferred/check/"+string(token), nil)
	if err != nil {
		return nil, err
	}
	var wr models.WorkflowResult
	err = json.Unmarshal(res.Body(), &wr)
	if err != nil {
		return nil, err
	}
	return &wr, nil
}

func (cli *Client) fetchDeferred(token DeferredToken) (io.Reader, error) {
	cli.logger.Debug("Downloading Deferred Ionburst object with token", token)
	return cli.doGetBinary("api/data/deferred/fetch/"+string(token), nil)
}
