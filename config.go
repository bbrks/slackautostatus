package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/nlopes/slack"
)

const (
	defaultUpdateFrequency = time.Minute * 5
	defaultConfigFileName  = "config.json"
)

var (
	verbose               bool
	defaultConfigFilePath = filepath.Join(os.Getenv("HOME"), ".config/slackautostatus")
)

// loadConfig loads the config from a file
func loadConfig() (*config, error) {
	c := config{
		updateFrequency: defaultUpdateFrequency,
	}

	f, err := os.Open(filepath.Join(defaultConfigFilePath, defaultConfigFileName))
	defer f.Close()
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(f).Decode(&c)
	if err != nil {
		return nil, err
	}

	c.updateFrequency, err = time.ParseDuration(c.UpdateFrequency)
	if err != nil {
		return nil, err
	}

	verbose = c.Verbose

	return &c, err
}

type config struct {
	Verbose         bool        `json:"verbose"`
	UpdateFrequency string      `json:"update_frequency"`
	Workspaces      []workspace `json:"workspaces"`

	updateFrequency time.Duration
}

type workspace struct {
	Name    string            `json:"name"`
	Token   string            `json:"token"`
	SSIDMap map[string]status `json:"ssids"`

	api *slack.Client
}

type status struct {
	Text  string `json:"text"`
	Emoji string `json:"emoji"`
}

var exampleConfig = config{
	Verbose:         true,
	UpdateFrequency: "5m",
	Workspaces: []workspace{
		{
			Name:  "Example Workspace",
			Token: "YOUR_SLACK_LEGACY_WORKSPACE_TOKEN",
			SSIDMap: map[string]status{
				"YOUR_OFFICE_SSID": {
					Text:  "Working from office",
					Emoji: ":office:",
				},
				"YOUR_HOME_SSID": {
					Text:  "Working remotely",
					Emoji: ":house_with_garden:",
				},
			},
		},
		{
			Name:  "Another Workspace",
			Token: "DIFFERENT_SLACK_LEGACY_TOKEN",
			SSIDMap: map[string]status{
				"YOUR_OFFICE_SSID": {
					Text:  "Busy at work",
					Emoji: ":hourglass:",
				},
				"YOUR_COFFEE_SHOP_WIFI": {
					Text:  "Caffeinated",
					Emoji: ":coffee:",
				},
			},
		},
	},
}

func createExampleConfig() error {
	b, err := json.MarshalIndent(exampleConfig, "", "  ")
	if err != nil {
		return err
	}
	err = os.MkdirAll(defaultConfigFilePath, 0700)
	if err != nil {
		return err
	}
	newFile, err := os.OpenFile(filepath.Join(defaultConfigFilePath, defaultConfigFileName), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer newFile.Close()
	if err != nil {
		return err
	}
	_, err = newFile.Write(b)
	if err != nil {
		return err
	}
	return nil
}
