#!/bin/sh

echo "${HOSTS}" | jq -c -R 'split(" ")'
