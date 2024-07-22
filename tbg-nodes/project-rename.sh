#!/bin/bash

set -e

CURRENT_NAME="BaasPhx"
CURRENT_OTP="baas_phx"

NEW_NAME="TbgNodes"
NEW_OTP="tbg_nodes"

ack -l $CURRENT_NAME --ignore-file=is:project-rename.sh | xargs sed -i '' -e "s/$CURRENT_NAME/$NEW_NAME/g"
ack -l $CURRENT_OTP --ignore-file=is:project-rename.sh | xargs sed -i '' -e "s/$CURRENT_OTP/$NEW_OTP/g"

mv lib/$CURRENT_OTP lib/$NEW_OTP
mv lib/$CURRENT_OTP.ex lib/$NEW_OTP.ex
mv lib/${CURRENT_OTP}_web lib/${NEW_OTP}_web
mv lib/${CURRENT_OTP}_web.ex lib/${NEW_OTP}_web.ex
mv test/${CURRENT_OTP}_web test/${NEW_OTP}_web
