[![Build Status](https://travis-ci.org/theodesp/binpack.svg)](https://travis-ci.org/theodesp/binpack)
[![GoDoc](https://godoc.org/github.com/theodesp/binpack?status.svg)](https://godoc.org/github.com/theodesp/binpack)
[![Go Report Card](https://goreportcard.com/badge/github.com/theodesp/binpack)](https://goreportcard.com/badge/github.com/theodesp/binpack)
[![codecov.io](https://codecov.io/github/theodesp/binpack/branch/master/graph/badge.svg)](https://codecov.io/github/theodesp/binpack)
[![GolangCI](https://golangci.com/badges/github.com/golangci/golangci-lint.svg)](https://golangci.com/r/github.com/theodesp/binpack)
[![CodeFactor](https://www.codefactor.io/repository/github/theodesp/binpack/badge)](https://www.codefactor.io/repository/github/theodesp/binpack)
[![AppVeyor](https://ci.appveyor.com/api/projects/status/0wdydpxh3o74vc9t?svg=true)](https://ci.appveyor.com/api/projects/status/0wdydpxh3o74vc9t?svg=true)

# Binpack Encoding for Golang

**binpack** is a library implementing [encoding/binpack](http://binpack.liaohuqiu.net/). Which is a 
binary serialize format, it is json like, but smaller and faster.

## Install

    go get github.com/theodesp/binpack

## Documentation and Examples

Visit [godoc](http://godoc.org/github.com/theodesp/binpack) for general examples and public api reference.
See **.travis.yml** for supported **go** versions.


See implementation examples:


## Supported Types

- [x] string
- [x] []byte
- [x] []uint8
- [x] float32
- [x] float64
- [x] bool
- [x] nil
- [x] basic slices ([]string, []int, ...)
- [x] basic arrays ([n]string, [n]int, ...)


## Run tests

    go test -race
    
## Contributions

See our [Contributing](./CONTRIBUTING.md) guide. This project uses [allcontributors](https://allcontributors.org/).

## License

The [Apache](https://www.apache.org/licenses/LICENSE-2.0)
