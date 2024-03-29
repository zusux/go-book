# Booklist Service

This is the Booklist service

Generated with

```
micro new --namespace=zusux.book --type=service booklist
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: zusux.book.service.booklist
- Type: service
- Alias: booklist

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
./booklist-service
```

Build a docker image
```
make docker
```