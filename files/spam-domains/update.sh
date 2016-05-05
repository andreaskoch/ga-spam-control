#!/bin/sh

cat filters.json | jq ".items[].excludeDetails.expressionValue" | sed 's/[\"]//g' | sed $'s/|/\\\n/g' | sort | uniq > referrer-spam-domains.txt
