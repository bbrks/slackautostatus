package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/nlopes/slack"
)

func main() {
	config, err := loadConfig()
	if os.IsNotExist(err) {
		if err = createExampleConfig(); err != nil {
			fatalf("error creating example config: %s", err)
		}
		fatalf("example config created at: %s\n", filepath.Join(defaultConfigFilePath, defaultConfigFileName))
	} else if err != nil {
		fatalf("error reading config: %s", err)
	}

	ticker := time.NewTicker(config.updateFrequency)
	for {
		ssid, changed, err := getSSID()
		if err != nil {
			errorf("error getting SSID: %s\n", err)
			goto wait
		} else if !changed {
			logf("SSID not changed... Skipping.")
			goto wait
		}

		for _, w := range config.Workspaces {
			if w.api == nil {
				api := slack.New(w.Token)
				if _, err = api.AuthTest(); err != nil {
					errorf("%s: error authorising Slack token: %s", w.Name, err)
					continue
				}
				w.api = api
			}

			s, ok := w.SSIDMap[ssid]
			if !ok {
				logf("%s: SSID %q not configured. Removing custom status.\n", w.Name, ssid)
			} else {
				logf("%s: setting status: '%s %s'\n", w.Name, s.Text, s.Emoji)
			}
			err = w.api.SetUserCustomStatus(s.Text, s.Emoji, 0)
			if err != nil {
				errorf("%s: error setting status: %s\n", w.Name, err)
				continue
			}
		}

	wait:
		<-ticker.C
	}

}
