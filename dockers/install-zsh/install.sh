#!/bin/sh

if mountpoint -q /target; then
    echo "Installing zsh to /target"
    cp /bin/zsh /target/zsh
else
    echo "/target is not a mountpoint."
    echo "- re-run this with -v /opt/bin:/target"
fi
