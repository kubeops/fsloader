[![Go Report Card](https://goreportcard.com/badge/github.com/appscode/fsloader)](https://goreportcard.com/report/github.com/appscode/fsloader)

[Website](https://appscode.com) • [Slack](https://slack.appscode.com) • [Twitter](https://twitter.com/AppsCodeHQ)

# fsloader
Runs commands when a file changes.

## Why Fsloader?
Many applications require configuration via some combination of config files. These configuration artifacts
should be decoupled from image content in order to keep containerized applications portable.

`Fsloader` watches a specified file. In case of any update in file data `Fsloader` updates the mounted file and run an additional bash script.

---

**Fsloader collects anonymous usage statistics to help us learn how the software is being used and how we can improve it.
To disable stats collection, run the binary with the flag** `--analytics=false`.

---
