package ionburst

import (
	"encoding/json"
	"io"
	"os"

	"gitlab.com/ionburst/ionburst-sdk-go/models"
)

type DeferredToken string

func (cli *Client) GetClassifications() ([]string, error) {
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

func (cli *Client) Put(id string, reader io.Reader, classification string) error {
	cli.logger.Debug("Uploading Ionburst object", id)
	_, err := cli.doPostBinary("api/data/"+id, reader, classification)
	if err == nil {
		cli.logger.Debug("Ionburst Upload Complete for", id)
	}
	return err
}

func (cli *Client) PutSecrets(id string, reader io.Reader, classification string) error {
	cli.logger.Debug("Uploading Ionburst secret", id)
	_, err := cli.doPostBinary("api/secrets/"+id, reader, classification)
	if err == nil {
		cli.logger.Debug("Ionburst Put Complete for", id)
	}
	return err
}

func (cli *Client) PutFromFile(id string, file string, classification string) error {
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

func (cli *Client) PutSecretsFromFile(id string, file string, classification string) error {
	cli.logger.Debug("Uploading Ionburst secret", id, "from file", file)
	rdr, err := os.Open(file)
	if err != nil {
		return err
	}
	_, err = cli.doPostBinary("api/secrets/"+id, rdr, classification)
	if err == nil {
		cli.logger.Debug("Ionburst Upload Complete for", id)
	}
	return err
}

func (cli *Client) Get(id string) (io.Reader, error) {
	cli.logger.Debug("Downloading Ionburst object", id)
	return cli.doGetBinary("api/data/"+id, nil)
}

func (cli *Client) GetSecrets(id string) (io.Reader, error) {
	cli.logger.Debug("Downloading Ionburst secret", id)
	return cli.doGetBinary("api/secrets/"+id, nil)
}

func (cli *Client) GetWithLen(id string) (io.Reader, int64, error) {
	cli.logger.Debug("Downloading Ionburst object", id)
	return cli.doGetBinaryLen("api/data/"+id, nil)
}

func (cli *Client) GetSecretsWithLen(id string) (io.Reader, int64, error) {
	cli.logger.Debug("Downloading Ionburst secret", id)
	return cli.doGetBinaryLen("api/secrets/"+id, nil)
}

func (cli *Client) GetToFile(id string, file string) error {
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

func (cli *Client) GetSecretsToFile(id string, file string) error {
	cli.logger.Debug("Downloading Ionburst secret", id, "to file", file)
	rdr, err := cli.doGetBinary("api/secrets/"+id, nil)
	if err != nil {
		return err
	}
	wr, err := os.Create(file)
	if err != nil {
		return err
	}
	_, err = io.Copy(wr, rdr)
	cli.logger.Debug("Ionburst secret download", id, "complete")
	return err
}

func (cli *Client) Delete(id string) error {
	cli.logger.Debug("Starting Deletion of Ionburst object", id)
	_, err := cli.doDelete("api/data/"+id, nil)
	cli.logger.Debug("Deletion of Ionburst object", id, "complete")
	return err
}

func (cli *Client) DeleteSecrets(id string) error {
	cli.logger.Debug("Starting Deletion of Ionburst secret", id)
	_, err := cli.doDelete("api/secrets/"+id, nil)
	cli.logger.Debug("Deletion of Ionburst secret", id, "complete")
	return err
}

func (cli *Client) Head(id string) error {
	cli.logger.Debug("Querying Ionburst object", id)
	_, err := cli.doHead("api/data/"+id, nil)
	cli.logger.Debug("Query of Ionburst object", id, "complete")
	return err
}

func (cli *Client) HeadSecrets(id string) error {
	cli.logger.Debug("Querying Ionburst secret", id)
	_, err := cli.doHead("api/secrets/"+id, nil)
	cli.logger.Debug("Query of Ionburst secret", id, "complete")
	return err
}

func (cli *Client) HeadWithLen(id string) (int64, error) {
	cli.logger.Debug("Querying Ionburst object", id)
	return cli.doHeadLen("api/data/"+id, nil)
}

func (cli *Client) HeadSecretsWithLen(id string) (int64, error) {
	cli.logger.Debug("Querying Ionburst secret", id)
	return cli.doHeadLen("api/secrets/"+id, nil)
}

func (cli *Client) PutDeferred(id string, reader io.Reader, classification string) (DeferredToken, error) {
	cli.logger.Debug("Uploading Ionburst object deferred for: ", id)
	tk, err := cli.doPostBinaryDeferred("api/data/deferred/start/"+id, reader, classification)
	if err == nil {
		cli.logger.Debug("Ionburst Deferred Upload Started for: ", id, " token: ", tk)
	} else {
		return "", err
	}
	return DeferredToken(tk), nil
}

func (cli *Client) PutDeferredSecrets(id string, reader io.Reader, classification string) (DeferredToken, error) {
	cli.logger.Debug("Uploading Ionburst secret deferred for: ", id)
	tk, err := cli.doPostBinaryDeferred("api/secrets/deferred/start/"+id, reader, classification)
	if err == nil {
		cli.logger.Debug("Ionburst Deferred Upload Started for: ", id, " token: ", tk)
	} else {
		return "", err
	}
	return DeferredToken(tk), nil
}

func (cli *Client) GetDeferred(id string) (DeferredToken, error) {
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

func (cli *Client) GetDeferredSecrets(id string) (DeferredToken, error) {
	cli.logger.Debug("Downloading Ionburst secret deferred for: ", id)
	tk, err := cli.doGet("api/secrets/deferred/start/"+id, nil)
	if err == nil {

		var defResp models.DeferredTokenResponse
		err = json.Unmarshal(tk.Body(), &defResp)
		if err != nil {
			return "", err
		}
		cli.logger.Debug("Ionburst Deferred GET Started for: ", id, " token: ", defResp.DeferredToken)
		return DeferredToken(defResp.DeferredToken), nil
	} else {
		return "", err
	}
}

func (cli *Client) CheckDeferred(token DeferredToken) (*models.WorkflowResult, error) {
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

func (cli *Client) CheckDeferredSecrets(token DeferredToken) (*models.WorkflowResult, error) {
	res, err := cli.doGet("api/secrets/deferred/check/"+string(token), nil)
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

func (cli *Client) FetchDeferred(token DeferredToken) (io.Reader, error) {
	cli.logger.Debug("Downloading Deferred Ionburst object with token", token)
	return cli.doGetBinary("api/data/deferred/fetch/"+string(token), nil)
}

func (cli *Client) FetchDeferredSecrets(token DeferredToken) (io.Reader, error) {
	cli.logger.Debug("Downloading Deferred Ionburst secret with token", token)
	return cli.doGetBinary("api/secrets/deferred/fetch/"+string(token), nil)
}
