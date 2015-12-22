# deque

[![GoDoc](https://godoc.org/github.com/tebeka/deque?status.svg)](https://godoc.org/github.com/tebeka/deque)
[![Travis](https://travis-ci.org/tebeka/deque.svg?branch=master)](https://travis-ci.org/tebeka/deque)

Implementation of [Python deque][src] in Go.

## Speed

`deque` is pretty fast, run `make compare` to see comparison against some other
methods. On my machine I get:

```
$ make compare
Git head is 765f6b0
cd compare && go test -run NONE -bench . -v
testing: warning: no tests to run
PASS
BenchmarkHistAppend-4	 3000000	       517 ns/op
BenchmarkHistList-4  	 2000000	       702 ns/op
BenchmarkHistQueue-4 	 3000000	       576 ns/op
BenchmarkHistDeque-4 	 3000000	       423 ns/op
ok  	_/home/miki/Projects/go/src/github.com/tebeka/deque/compare	8.505s
```

[src]: https://hg.python.org/cpython/file/tip/Modules/_collectionsmodule.c
