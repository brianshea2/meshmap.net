#!/bin/bash

docker stop meshobserv
docker rm meshobserv

docker run --name meshobserv \
    --restart unless-stopped \
    -v /data:/data \
    -d meshobserv \
    -f /data/meshmap.net/website/nodes.json \
    -b /data/meshmap.net/blocklist.txt
