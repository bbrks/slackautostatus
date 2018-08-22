# slackautostatus

Automatically change your slack status based on WiFi SSID

## Install

```
go install github.com/bbrks/slackautostatus
```

## Running and configuration

The first run of slackautostatus will generate an example config at `~/.config/slackautostatus/config.json`

```
slackautostatus
```

After this, edit the config for your workspaces, [legacy tokens](https://api.slack.com/custom-integrations/legacy-tokens), SSIDs, and statuses, and re-run.
