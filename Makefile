test:
	go test -v

bench:
	go test -bench . -v


.PHONY: test bench
