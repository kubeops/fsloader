[![Build Status](https://github.com/kubeops/fsloader/workflows/CI/badge.svg)](https://github.com/kubeops/fsloader/actions?workflow=CI)
[![Slack](https://shields.io/badge/Join_Slack-salck?color=4A154B&logo=slack)](https://slack.appscode.com)
[![Twitter](https://img.shields.io/twitter/follow/kubeops.svg?style=social&logo=twitter&label=Follow)](https://twitter.com/intent/follow?screen_name=Kubeops)

# fsloader
Runs commands when a file changes.

## Why Fsloader?
Many applications require configuration via some combination of config files. These configuration artifacts
should be decoupled from image content in order to keep containerized applications portable.

`Fsloader` watches a specified file. In case of any update in file data `Fsloader` updates the mounted file and run an additional bash script.
