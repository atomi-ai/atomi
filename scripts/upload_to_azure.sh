#!/bin/sh

echo "SAS token: \$AZURE_SAS_TOKEN"
azcopy cp "${DRONE_REPO_OWNER}-${DRONE_REPO_NAME}-${DRONE_COMMIT_SHA:0:7}" "https://atomidrone.blob.core.windows.net/atomi-backend-drone-report${AZURE_SAS_TOKEN}" --recursive
