[![Go Report Card](https://goreportcard.com/badge/kubeops.dev/fsloader)](https://goreportcard.com/report/kubeops.dev/fsloader)

[Website](https://appscode.com) • [Slack](https://slack.appscode.com) • [Twitter](https://twitter.com/AppsCodeHQ)

# fsloader
Runs commands when a file changes.

## Why Fsloader?
Many applications require configuration via some combination of config files. These configuration artifacts
should be decoupled from image content in order to keep containerized applications portable.

`Fsloader` watches a specified file. In case of any update in file data `Fsloader` updates the mounted file and run an additional bash script.
