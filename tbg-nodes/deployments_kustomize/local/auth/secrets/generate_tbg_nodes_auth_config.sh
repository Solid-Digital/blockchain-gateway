#!/bin/bash
set -e
sed_replace="s/{{host}}/${POSTGRES_HOST}/g"

cat tbg-nodes-auth-dev.toml.template | sed "${sed_replace}" > tbg-nodes-auth-dev.toml
cat tbg-nodes-auth-test.toml.template | sed "${sed_replace}" > tbg-nodes-auth-test.toml
