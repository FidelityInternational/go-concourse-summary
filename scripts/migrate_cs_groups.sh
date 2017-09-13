#!/bin/bash

sanitised_cs_grous=$(sed 's/null/{}/g' <<< "${CS_GROUPS}")
echo "${sanitised_cs_grous}" | jq 'to_entries | map({group: .key, hosts: (.value | to_entries | map({fqdn: .key, pipelines: (.value | to_entries | map({name: .key, groups: (.value | to_entries | map(.value))}))}))})'
