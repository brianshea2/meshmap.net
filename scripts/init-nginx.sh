#!/bin/sh
envsubst '${SITE_TITLE} ${SITE_URL}' \
  < /data/meshmap.net/website/index.html.template \
  > /data/meshmap.net/website/index.html
