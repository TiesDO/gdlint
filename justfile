test target="": (_run_test "" target)

testv target="": (_run_test "-v" target)

build:
    go build -o ./bin/gdlint

_run_test flags target="":
    #!/usr/bin/env bash
    set -eo pipefail
    go test ./... {{flags}} {{ if target == "" { "" } else { "-run " + target } }} \
      | grep -v -e '\[no test files\]' -e '\[no tests to run\]'
