#!/bin/bash
set -e
REDIS_URL=$(echo "${REDIS_URL}" | sed 's/\//\\\//g')
sed_replace="s/{{POSTGRES_HOST}}/${POSTGRES_HOST}/g;s/{{REDIS_URL}}/${REDIS_URL}/g"

cat local-tbg-nodes-config.template | sed "${sed_replace}" > local-tbg-nodes-config
