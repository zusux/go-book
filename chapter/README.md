# Chapter Service

This is the Chapter service

Generated with

```
micro new --namespace=zusux.book --type=service chapter
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: zusux.book.service.chapter
- Type: service
- Alias: chapter

## Dependencies

Micro services depend on service discovery. The default is multicast DNS, a zeroconf system.

In the event you need a resilient multi-host setup we recommend etcd.

```
# install etcd
brew install etcd

# run etcd
etcd
```

## Usage

A Makefile is included for convenience

Build the binary

```
make build
```

Run the service
```
./chapter-service
```

Build a docker image
```
make docker
```