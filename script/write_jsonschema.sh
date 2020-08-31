#!/bin/bash
set -eu

tmpl=$(cat << EOM
// Code generated by script/write_jsonschema.sh. DO NOT EDIT.
// This file was generated by script/write_jsonschema.sh.
// using data from: jsonschema/$1.json
package schema

var %s = []byte(\`%s
\`)
EOM
)

json=$(< "jsonschema/$1.json")

# shellcheck disable=SC2059
printf "$tmpl\n" "$2" "$json" > "$3"