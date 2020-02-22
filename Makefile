test:
	go test -v

bench:
	go test -run NONE -bench . -v

compare:
	@echo Git head is $(shell git head)
	cd _compare && go test -run NONE -bench . -v

.PHONY: test bench compare
