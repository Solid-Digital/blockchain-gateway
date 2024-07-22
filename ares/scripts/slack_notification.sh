#!/bin/sh

# Buildchannel is #build on slack
BUILDCHANNEL=https://hooks.slack.com/services/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
ARES_BRANCH_URL=https://ares.branch.unchain.io

# Allow override for other channels
if [ -z $SLACK_WEBHOOK ]; then SLACK_WEBHOOK=$BUILDCHANNEL; fi
# Allow override of Ares url
if [ -z $ARES_URL ]; then ARES_URL=$ARES_BRANCH_URL; fi

# The message in JSON format
MESSAGE='{"text": "Deployment of ares was completed to {{site}}"}'

# inject site
PARSED_DATA=$(echo $MESSAGE | sed -e "s|{{site}}|$ARES_URL|g")

curl -H 'Content-type: application/json' --data "$PARSED_DATA" $SLACK_WEBHOOK
