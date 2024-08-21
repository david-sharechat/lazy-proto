# Lazy Proto serde example

This repo hold a simple example and benchmark for lazy proto serde.

##  Setup 

### Install protoc and go plugin. 

- For protoc, see https://grpc.io/docs/protoc-installation/
- For Go plugin:

    ```shell
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    ```

### Compile protos:
```shell
protoc --go_out=. protos/example.proto
```

### Run test
```shell
go test
```

### Run benchmark
```shell
go test -bench=. -benchmem -count=10 | tee stats.txt
benchstat -row .name -col .fullname stats.txt
```

Sample Output:
```shell
goos: darwin
goarch: arm64
pkg: lazy-proto
      │  */Naive-10  │              */Lazy-10              │
      │    sec/op    │   sec/op     vs base                │
Merge   16.821m ± 3%   2.677m ± 6%  -84.09% (p=0.000 n=10)

      │  */Naive-10  │              */Lazy-10               │
      │     B/op     │     B/op      vs base                │
Merge   8.479Mi ± 0%   5.993Mi ± 0%  -29.31% (p=0.000 n=10)

      │  */Naive-10  │              */Lazy-10              │
      │  allocs/op   │  allocs/op   vs base                │
Merge   141.16k ± 0%   40.08k ± 0%  -71.61% (p=0.000 n=10)
```
