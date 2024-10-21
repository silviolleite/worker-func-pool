# Worker Func Pool

![coverage](https://raw.githubusercontent.com/silviolleite/worker-func-pool/badges/.badges/main/coverage.svg)

---

This package provides a worker pool to functions.

Worker func pools can be used to efficiently parallelize the processing of functions with different signatures
and take advantage of multicore processors for better performance.

### Usage

See [example](/example/main.go).

### Install

---

Manual install:

```bash
go get -u github.com/silviolleite/worker-func-pool
```

Golang import:

```go
import "github.com/silviolleite/worker-func-pool"
```

### Notes

While worker func pools help control the number of concurrent tasks,
they may still result in more context switches between goroutines.
This can result in overhead, which can cancel out some of the performance benefits gained from using worker pools.
To reduce context switching overhead,
it’s critical to strike the right balance between the number of workers and the workload
