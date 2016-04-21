#!/bin/sh

cat filters.json | jq ".items[].excludeDetails.expressionValue" | sed 's/[\"]//g' | sed $'s/|/\\\n/g' | sort | uniq > referer-spam-domains.txt
