package ionburst

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"gitlab.com/ionburst/ionburst-sdk-go/models"
	"io"
	"os"
)

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (cli *Client) makeClientFromCreds() (string, *resty.Request, error) {

	cli.logger.Debug("Creating new ionburst client")

	creds, err := cli.ionConfig.GetDefaultCredsProfile()
	if err != nil {
		return "", nil, err
	}

	if creds == nil || creds.IonburstID == "" || creds.IonburstKey == "" {
		return "", nil, errors.New("No credentials provided to perform operation")
	}

	if cli.auth == nil {
		cli.auth, err = cli.doSignin(creds.IonburstUri, creds.IonburstID, creds.IonburstKey, "")
		if err != nil {
			return "", nil, err
		}
	}

	if cli.auth == nil || cli.auth.IdToken == "" {
		return "", nil, errors.New("Cannot proceed with request as sign-in hasn't taken place yet")
	}

	client := resty.New()

	r := client.R().SetHeader("Authorization", "Bearer "+cli.auth.IdToken)

	return creds.IonburstUri, r, nil

}

func (cli *Client) doSignin(uri string, id string, key string, refreshToken string) (*models.AuthResponse, error) {
	var auth models.Auth
	if refreshToken == "" {
		auth = models.Auth{
			Username: id,
			Password: key,
		}
	} else {
		auth = models.Auth{
			Username:     id,
			RefreshToken: refreshToken,
		}
	}
	client := resty.New()
	r, err := client.R().SetBody(auth).SetHeader("Content-Type", "application/json").Post(uri + "api/signin")
	if err != nil {
		return nil, err
	}

	statusCode := r.StatusCode()
	cli.logger.Debug("Starting authentication for: ", id, " uri: ", uri)
	if statusCode != 200 {
		return nil, errors.New("Authentication has failed: " + r.String())
	} else {
		var authResp models.AuthResponse
		err = json.Unmarshal(r.Body(), &authResp)
		if err != nil {
			return nil, err
		}
		return &authResp, nil
	}
}

func (cli *Client) checkResponse(res *resty.Response) (bool, error) {
	if !res.IsSuccess() {
		if res.StatusCode() == 401 {
			creds, err := cli.ionConfig.GetDefaultCredsProfile()
			if err != nil {
				return false, err
			}
			rtoken := ""
			if cli.auth != nil && cli.auth.RefreshToken != "" {
				rtoken = cli.auth.RefreshToken
			}
			cli.logger.Debug("Signing in to Ionburst using profile", cli.ionConfig.DefaultProfile)

			cli.auth, err = cli.doSignin(creds.IonburstUri, creds.IonburstID, creds.IonburstKey, rtoken)
			if err != nil {
				return false, err
			}
			return true, nil
		}
		return false, errors.New("Error performing Ionburst API Operation: [" + res.Request.Method + "] " + res.Request.URL + " :: " + fmt.Sprintf("%d", res.StatusCode()) + " -> " + res.String())
	}
	return false, nil
}

func (cli *Client) doGet(url string, params map[string]string) (*resty.Response, error) {
	if uri, req, err := cli.makeClientFromCreds(); err != nil {
		return nil, err
	} else {
		if uri[len(uri)-1] != '/' {
			uri += "/"
		}
		if params != nil {
			req.SetQueryParams(params)
		}
		res, err := req.Get(uri + url)
		if err != nil {
			return nil, err
		}
		if redo, err := cli.checkResponse(res); err != nil {
			return nil, err
		} else if redo {
			return cli.doGet(url, params)
		} else {
			return res, nil
		}
	}
}

func (cli *Client) doGetBinary(url string, params map[string]string) (io.Reader, error) {
	if uri, req, err := cli.makeClientFromCreds(); err != nil {
		return nil, err
	} else {
		if uri[len(uri)-1] != '/' {
			uri += "/"
		}
		if params != nil {
			req.SetQueryParams(params)
		}
		res, err := req.SetDoNotParseResponse(true).Get(uri + url)
		if err != nil {
			return nil, err
		}
		if redo, err := cli.checkResponse(res); err != nil {
			return nil, err
		} else if redo {
			return cli.doGetBinary(url, params)
		} else {
			return res.RawResponse.Body, nil
		}
	}
}

func (cli *Client) doGetBinaryLen(url string, params map[string]string) (io.Reader, int64, error) {
	if uri, req, err := cli.makeClientFromCreds(); err != nil {
		return nil, 0, err
	} else {
		if uri[len(uri)-1] != '/' {
			uri += "/"
		}
		if params != nil {
			req.SetQueryParams(params)
		}
		res, err := req.SetDoNotParseResponse(true).Get(uri + url)
		if err != nil {
			return nil, 0, err
		}
		if redo, err := cli.checkResponse(res); err != nil {
			return nil, 0, err
		} else if redo {
			return cli.doGetBinaryLen(url, params)
		} else {
			return res.RawResponse.Body, res.Size(), nil
		}
	}
}

func (cli *Client) doGetBinaryDeferred(url string, params map[string]string) (string, error) {
	if uri, req, err := cli.makeClientFromCreds(); err != nil {
		return "", err
	} else {
		if uri[len(uri)-1] != '/' {
			uri += "/"
		}
		if params != nil {
			req.SetQueryParams(params)
		}
		res, err := req.SetDoNotParseResponse(true).Get(uri + url)
		if err != nil {
			return "", err
		}
		if redo, err := cli.checkResponse(res); err != nil {
			return "", err
		} else if redo {
			return cli.doGetBinaryDeferred(url, params)
		} else {
			return res.String(), nil
		}
	}
}

func (cli *Client) doGetBinaryLenDeferred(url string, params map[string]string) (string, int64, error) {
	if uri, req, err := cli.makeClientFromCreds(); err != nil {
		return "", 0, err
	} else {
		if uri[len(uri)-1] != '/' {
			uri += "/"
		}
		if params != nil {
			req.SetQueryParams(params)
		}
		res, err := req.Get(uri + url)
		if err != nil {
			return "", 0, err
		}
		if redo, err := cli.checkResponse(res); err != nil {
			return "", 0, err
		} else if redo {
			return cli.doGetBinaryLenDeferred(url, params)
		} else {
			return res.String(), res.Size(), nil
		}
	}
}

func (cli *Client) doPost(url string, payload interface{}, params map[string]string) (*resty.Response, error) {
	if uri, req, err := cli.makeClientFromCreds(); err != nil {
		return nil, err
	} else {
		if uri[len(uri)-1] != '/' {
			uri += "/"
		}
		if params != nil {
			req.SetQueryParams(params)
		}
		res, err := req.Post(uri + url)
		if err != nil {
			return nil, err
		}
		if redo, err := cli.checkResponse(res); err != nil {
			return nil, err
		} else if redo {
			return cli.doPost(url, payload, params)
		} else {
			return res, nil
		}
	}

}

func (cli *Client) doPut(url string, payload interface{}, params map[string]string) (*resty.Response, error) {
	if uri, req, err := cli.makeClientFromCreds(); err != nil {
		return nil, err
	} else {
		if uri[len(uri)-1] != '/' {
			uri += "/"
		}
		if params != nil {
			req.SetQueryParams(params)
		}
		res, err := req.Put(uri + url)
		if err != nil {
			return nil, err
		}
		if redo, err := cli.checkResponse(res); err != nil {
			return nil, err
		} else if redo {
			return cli.doPut(url, payload, params)
		} else {
			return res, nil
		}
	}

}

func (cli *Client) doPostBinary(url string, reader io.Reader, classification string) (*resty.Response, error) {
	if uri, req, err := cli.makeClientFromCreds(); err != nil {
		return nil, err
	} else {
		if uri[len(uri)-1] != '/' {
			uri += "/"
		}
		if classification != "" {
			req.SetQueryParam("classstr", classification)
		}
		res, err := req.SetBody(reader).Post(uri + url)
		if err != nil {
			return nil, err
		}
		if redo, err := cli.checkResponse(res); err != nil {
			return nil, err
		} else if redo {
			return cli.doPostBinary(url, reader, classification)
		} else {
			return res, nil
		}
	}

}

func (cli *Client) doPostBinaryDeferred(url string, reader io.Reader, classification string) (string, error) {
	if uri, req, err := cli.makeClientFromCreds(); err != nil {
		return "", err
	} else {
		if uri[len(uri)-1] != '/' {
			uri += "/"
		}
		if classification != "" {
			req.SetQueryParam("classstr", classification)
		}
		res, err := req.SetBody(reader).Post(uri + url)
		if err != nil {
			return "", err
		}
		if redo, err := cli.checkResponse(res); err != nil {
			return "", err
		} else if redo {
			return cli.doPostBinaryDeferred(url, reader, classification)
		} else {
			var defResp models.DeferredTokenResponse
			err = json.Unmarshal(res.Body(), &defResp)
			if err != nil {
				return "", err
			}
			cli.logger.Debug("Got Deferred status")
			return defResp.DeferredToken, nil
		}
	}

}

func (cli *Client) doPutBinary(url string, reader io.Reader, classification string) (*resty.Response, error) {
	if uri, req, err := cli.makeClientFromCreds(); err != nil {
		return nil, err
	} else {
		if uri[len(uri)-1] != '/' {
			uri += "/"
		}
		if classification != "" {
			req.SetQueryParam("classstr", classification)
		}
		res, err := req.SetBody(reader).Put(uri + url)
		if err != nil {
			return nil, err
		}
		if redo, err := cli.checkResponse(res); err != nil {
			return nil, err
		} else if redo {
			return cli.doPutBinary(url, reader, classification)
		} else {
			return res, nil
		}
	}

}

func (cli *Client) doDelete(url string, params map[string]string) (*resty.Response, error) {
	if uri, req, err := cli.makeClientFromCreds(); err != nil {
		return nil, err
	} else {
		if uri[len(uri)-1] != '/' {
			uri += "/"
		}
		if params != nil {
			req.SetQueryParams(params)
		}
		res, err := req.Delete(uri + url)
		if err != nil {
			return nil, err
		}
		if redo, err := cli.checkResponse(res); err != nil {
			return nil, err
		} else if redo {
			return cli.doDelete(url, params)
		} else {
			return res, nil
		}
	}

}
