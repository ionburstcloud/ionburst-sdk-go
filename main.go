package ionburst

import (
	"errors"
	"github.com/gobuffalo/envy"
	"gitlab.com/ionburst/ionburst-sdk-go/models"
	"go.uber.org/zap"
)

var FileList = []string{
	"",
	".ion/config.json",
}

func NewClient() (*Client, error) {
	return NewClientDebug(false)
}

func NewClientDebug(debug bool) (*Client, error) {
	cli := &Client{}
	if debug {
		logger, _ := zap.NewDevelopment()
		sugar := logger.Sugar()
		cli.logger = sugar
	} else {
		logger, _ := zap.NewProduction()
		sugar := logger.Sugar()
		cli.logger = sugar
	}
	return cli.init()
}

func NewClientPathAndProfile(configPath string, credentialsProfile string, debug bool) (*Client, error) {
	cli := &Client{}
	if debug {
		logger, _ := zap.NewDevelopment()
		sugar := logger.Sugar()
		cli.logger = sugar
	} else {
		logger, _ := zap.NewProduction()
		sugar := logger.Sugar()
		cli.logger = sugar
	}
	return cli.initWithPathAndProfile(configPath, credentialsProfile)
}

func NewClientWithCredentials(uri string, ionburstID string, ionburstKey string, debug bool) (*Client, error) {
	cli := &Client{}
	if debug {
		logger, _ := zap.NewDevelopment()
		sugar := logger.Sugar()
		cli.logger = sugar
	} else {
		logger, _ := zap.NewProduction()
		sugar := logger.Sugar()
		cli.logger = sugar
	}
	return cli.initWithCreds(uri, ionburstID, ionburstKey)
}

type Client struct {
	ionConfig *IonConfig
	auth      *models.AuthResponse
	logger    *zap.SugaredLogger
}

func (cli *Client) init() (*Client, error) {

	//load what needs to be loaded int he correct order
	//errList := []error{}

	// 1. Try Load config, do we have environment variables?

	ionburstConfigPath := envy.Get("IONBURST_CONFIG", GetDefaultIonburstConfigPath())
	ionburstCredsProfileName := envy.Get("IONBURST_CREDENTIALS_PROFILE", "")
	conf, err := LoadIonConfig(cli, ionburstConfigPath)
	if err != nil {
		//we couldnt load a config in either the default or specified locations, lets create an empty one
		cli.ionConfig = NewEmptyIonConfig(cli)

		//now we check if the environment variables are there...

		cli.loadEnvironmentVariables()

	} else {
		cli.ionConfig = conf
	}

	if ionburstCredsProfileName != "" {
		if err := cli.ionConfig.SetDefaultProfile(ionburstCredsProfileName); err != nil {
			return nil, err
		}
	}

	return cli.initAuth()
}

func (cli *Client) initWithPathAndProfile(ionburstConfigPath string, ionburstCredsProfileName string) (*Client, error) {

	//load what needs to be loaded int he correct order
	//errList := []error{}

	// 1. Try Load config, do we have environment variables?

	conf, err := LoadIonConfig(cli, ionburstConfigPath)
	if err != nil {
		//we couldnt load a config in either the default or specified locations, lets create an empty one
		cli.ionConfig = NewEmptyIonConfig(cli)

		//now we check if the environment variables are there...

		cli.loadEnvironmentVariables()

	} else {
		cli.ionConfig = conf
	}

	if ionburstCredsProfileName != "" {
		if err := cli.ionConfig.SetDefaultProfile(ionburstCredsProfileName); err != nil {
			return nil, err
		}
	}

	return cli.initAuth()
}

func (cli *Client) initWithCreds(uri string, ionburstID string, ionburstKey string) (*Client, error) {
	if ionburstID == "" || ionburstKey == "" {
		return nil, errors.New("Please supply an Ionburst ID and Ionburst Key credential set")
	}
	//since we have credentials these get put into a new CredentialsProfile
	conf := NewIonConfig(cli, uri, ionburstID, ionburstKey)

	cli.ionConfig = conf

	return cli.initAuth()
}

func (cli *Client) initAuth() (*Client, error) {
	_, _, err := cli.makeClientFromCreds()
	if err != nil {
		return nil, err
	}
	return cli, nil
}

func (cli *Client) loadEnvironmentVariables() {
	ionburstID := envy.Get("IONBURST_ID", "")
	ionburstKey := envy.Get("IONBURST_KEY", "")
	ionburstUri := envy.Get("IONBURST_URI", "")
	cli.ionConfig.UpsertCredsProfile(cli.ionConfig.DefaultProfile, ionburstUri, ionburstID, ionburstKey)
}