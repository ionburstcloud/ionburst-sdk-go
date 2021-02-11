package ionburst

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gobuffalo/envy"
	"gopkg.in/ini.v1"
)

const DefaultIonburstConfigPath = ".ionburst/credentials"
const DefaultIonburstCredentialsProfileName = "default"

func GetDefaultIonburstConfigPath() string {
	cwd := envy.Get("HOME", "./")

	if cwd[:2] == "./" {
		cwd, _ = os.Getwd()
	}

	configPath := filepath.Join(cwd, DefaultIonburstConfigPath)
	return configPath
}

func NewEmptyIonConfig(cli *Client) *IonConfig {
	cli.logger.Debug("Creating new empty ionburst config")
	cwd := envy.Get("HOME", "./")

	if cwd[:2] == "./" {
		cwd, _ = os.Getwd()
	}

	configPath := filepath.Join(cwd, DefaultIonburstConfigPath)

	config := &IonConfig{
		cli:            cli,
		isNew:          true,
		file:           configPath,
		DefaultProfile: DefaultIonburstCredentialsProfileName,
		Profiles:       map[string]*CredentialsProfile{},
	}

	return config.init()

}

func NewIonConfig(cli *Client, uri string, ionburstID string, ionburstKey string) *IonConfig {

	cwd := envy.Get("HOME", "./")

	if cwd[:2] == "./" {
		cwd, _ = os.Getwd()
	}

	configPath := filepath.Join(cwd, DefaultIonburstConfigPath)

	return NewIonConfigWithFilePaths(cli, uri, ionburstID, ionburstKey, configPath)

}

func NewIonConfigWithFilePaths(cli *Client, uri string, ionburstID string, ionburstKey string, configFile string) *IonConfig {
	cli.logger.Debug("Creating new ionburst config with supplied credentials")
	config := &IonConfig{
		cli:            cli,
		isNew:          true,
		file:           configFile,
		DefaultProfile: DefaultIonburstCredentialsProfileName,
		Profiles:       map[string]*CredentialsProfile{},
	}

	config.UpsertCredsProfile(DefaultIonburstCredentialsProfileName, uri, ionburstID, ionburstKey)

	return config.init()

}

func LoadIonConfig(cli *Client, configFile string) (*IonConfig, error) {
	if !FileExists(configFile) {
		return nil, os.ErrNotExist
	}
	//load the config file...

	cli.logger.Debug("Loading ionburst config from", configFile)

	ba, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var out IonConfig

	if err := json.Unmarshal(ba, &out); err != nil {
		return nil, err
	}

	out.cli = cli

	return (&out).init(), nil

}

type IonConfig struct {
	cli            *Client
	file           string
	isNew          bool
	DefaultProfile string                         `json:"DefaultProfile,omitempty"`
	Profiles       map[string]*CredentialsProfile `json:"Profiles"`
}

func (conf *IonConfig) init() *IonConfig {
	if err := conf.LoadIniCreds(); err != nil {
		conf.cli.logger.Warn("No credentials file found, environment variables only")
	}
	return conf
}

func (conf *IonConfig) GetDefaultCredsProfile() (*CredentialsProfile, error) {
	profName := conf.DefaultProfile
	if profName == "" {
		profName = DefaultIonburstCredentialsProfileName
	}
	return conf.GetCredsProfile(profName)
}

func (conf *IonConfig) SetDefaultProfile(profileName string) error {
	p, err := conf.GetCredsProfile(profileName)
	if err != nil {
		return err
	} else if p == nil {
		return errors.New("Cannot find profile to set default")
	}
	conf.DefaultProfile = profileName
	return nil
}

func (conf *IonConfig) SaveConfig() error {
	return conf.SaveConfigToFile(conf.file)
}

func (conf *IonConfig) SaveConfigToFile(file string) error {
	return nil
}

func (conf *IonConfig) LoadIniCreds() error {

	//using the path for current config file - look and see if the ini is there a "credentials" file

	opath := filepath.Dir(conf.file)
	credsPath := filepath.Join(opath, "credentials")

	if FileExists(credsPath) {
		//lets load the credentials from the ini..

		creds, err := ini.Load(credsPath)
		if err != nil {
			return err
		}
		for _, sec := range creds.Sections() {
			conf.UpsertCredsProfile(sec.Name(), sec.Key("ionburst_uri").String(), sec.Key("ionburst_id").String(), sec.Key("ionburst_key").String())
		}
	}

	return nil
}

func (conf *IonConfig) GetCredsProfile(name string) (*CredentialsProfile, error) {
	if conf.Profiles == nil {
		return nil, errors.New("Unable to find credentials profile " + name + " in config")
	} else if v, ok := conf.Profiles[name]; !ok || v == nil {
		return nil, errors.New("Unable to find credentials profile " + name + " in config")
	} else {
		return v, nil
	}
}

func (conf *IonConfig) UpsertCredsProfile(name string, uri string, ionburstID string, ionburstKey string) *CredentialsProfile {
	if len(uri) > 0 && uri[len(uri)-1] != '/' {
		uri += "/"
	}
	if conf.Profiles == nil {
		conf.Profiles = map[string]*CredentialsProfile{}
	} else if v, ok := conf.Profiles[name]; !ok || v == nil {
		conf.Profiles[name] = NewCredentialsProfile(uri, ionburstID, ionburstKey)
	} else {
		if uri != "" {
			conf.Profiles[name].IonburstURI = uri
		}
		if ionburstID != "" {
			conf.Profiles[name].IonburstID = ionburstID
		}
		if ionburstKey != "" {
			conf.Profiles[name].IonburstKey = ionburstKey
		}
	}
	return conf.Profiles[name]
}

func (conf *IonConfig) RemoveCredsProfile(name string) error {
	if conf.Profiles == nil {
		return errors.New("Unable to find credentials profile " + name + " in config")
	} else if v, ok := conf.Profiles[name]; !ok || v == nil {
		return errors.New("Unable to find credentials profile " + name + " in config")
	} else {
		delete(conf.Profiles, name)
		return nil
	}
}

func NewCredentialsProfile(uri string, ionburstID string, ionburstKey string) *CredentialsProfile {

	return &CredentialsProfile{
		IonburstID:  ionburstID,
		IonburstKey: ionburstKey,
		IonburstURI: uri,
	}
}

type CredentialsProfile struct {
	IonburstURI string `json:"URI"`
	IonburstID  string `json:"ID"`
	IonburstKey string `json:"KEY"`
}
