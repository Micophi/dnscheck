#!/bin/bash

# https://urlhaus.abuse.ch/downloads/hostfile/

DOMAIN_LIST_FILE=domains-urlhaus.txt
FEED_URL=https://urlhaus.abuse.ch/downloads/hostfile/

wget "$FEED_URL" -O "$DOMAIN_LIST_FILE"
tail -n +9 "$DOMAIN_LIST_FILE" | sponge "$DOMAIN_LIST_FILE"
head -n -1 "$DOMAIN_LIST_FILE" | sponge "$DOMAIN_LIST_FILE"
cut -f2- "$DOMAIN_LIST_FILE" | sponge "$DOMAIN_LIST_FILE"