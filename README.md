# parvati-api-client
Golang command line client for Parvati's API. Alpha quality. There be bugs.

## Overview

This should repo contains an implementation of the Parvati interface in
golang, as well as a command line client which uses it. They should
really be split out into separate repositories.

Configuration for the interface can be done from the command line and
from a `gitconfig` file local to your home directory (or specified).

Todo: Allow `.netrc` style password loading rather than relying on the
config file. Perhaps find modules to integrate with OS specific password
managers.

For more info run the resultant binary with `--help`.
