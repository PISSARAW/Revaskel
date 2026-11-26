.PHONY: demo test

demo:
	@set -eu; tmp=$$(mktemp -d); trap 'rm -rf "$$tmp"' EXIT; go build -o "$$tmp/revaskel" .; "$$tmp/revaskel" examples/complete; code=0; "$$tmp/revaskel" examples/partial || code=$$?; test "$$code" -eq 1

test:
	go test ./...

.PHONY: setup quality

setup:
	git config --local core.hooksPath .githooks

quality:
	sh quality/gate.sh staged
