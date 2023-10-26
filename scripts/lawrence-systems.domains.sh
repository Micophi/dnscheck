#!/bin/bash

# From the forums https://forums.lawrencesystems.com/t/which-is-the-best-dns-for-secure-browsing-cloudflare-quad9-nextdns-and-adguard-dns-youtube-release/18910/5
# The list comes from zonefile.io and filters it

DOMAIN_LIST_FILE=domains-ls.txt
FEED_URL=https://zonefiles.io/f/compromised/domains/live/

wget "$FEED_URL" -O "$DOMAIN_LIST_FILE"
tail -n +2 "$DOMAIN_LIST_FILE" | sponge "$DOMAIN_LIST_FILE"
grep -E '^[a-zA-Z0-9-]+\.(com|net)$' "$DOMAIN_LIST_FILE" | sponge "$DOMAIN_LIST_FILE"