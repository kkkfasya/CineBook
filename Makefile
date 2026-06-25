# i only add this to remember the flag of go compiler

.PHONY: test test_no_cache test_race

test:
	go test ./... -v

test_no_cache:
	go test ./... -v -count=1

test_race:
	go test ./... -v -count=1 -race
