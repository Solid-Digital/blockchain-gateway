#!/bin/sh


get_latest_release() {
  curl --silent "https://api.github.com/repos/$1/releases/latest" | # Get latest release from GitHub api
    grep '"tag_name":' |                                            # Get tag line
    sed -E 's/.*"([^"]+)".*/\1/'                                    # Pluck JSON value
}


if mountpoint -q /target; then
    echo "Installing docker-compose to /target"
    release=$(get_latest_release docker/compose)

    curl -L https://github.com/docker/compose/releases/download/"$release"/docker-compose-`uname -s`-`uname -m` > /target/docker-compose
    chmod +x /target/docker-compose
else
    echo "/target is not a mountpoint."
    echo "- re-run this with -v /opt/bin:/target"
fi
