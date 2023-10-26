#!/bin/bash
DOMAIN_LIST_FILE=domains-ls.txt
FEED_URL=https://zonefiles.io/f/compromised/domains/live/

wget "$FEED_URL" -O "$DOMAIN_LIST_FILE"
tail -n +2 "$DOMAIN_LIST_FILE" | sponge "$DOMAIN_LIST_FILE"
