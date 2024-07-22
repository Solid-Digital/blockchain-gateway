#!/bin/sh
set -e

# Create default buckets
mkdir -p /data/s3dev
mkdir -p /data/s3test

cp -r /testdata/* /data/s3dev/
cp -r /testdata/* /data/s3test/

# Proceed to the default entrypoint
/usr/bin/docker-entrypoint.sh "$@"
