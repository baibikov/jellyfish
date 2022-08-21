# Jellyfish in-memory broker message.
[![PkgGoDev](https://img.shields.io/badge/go.dev-docs-007d9c?style=flat-square&logo=go&logoColor=white)](https://pkg.go.dev/github.com/baibikov/jellyfish)
[![GitHub release](https://img.shields.io/github/release/baibikov/jellyfish.svg?style=flat-square)](https://https://github.com/baibikov/jellyfish/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/baibikov/jellyfish)](https://goreportcard.com/report/github.com/baibikov/jellyfish)

![Jellyfish](logo.png?raw=true "Jellyfish Logo")

### Philosophy:
-----------

A basic message broker that can solve the most pressing development problems. 
_The main idea:_ scale without pain. 

### Design:
-----------

**_[SENDING]_:** 
The messaging design based on grpc format. 

**_[TOPIC]_:** 
The topic is a route (place) where a message will be sent and where it will be stored.

**_[PARTITION]_**: 
The partition is a way to scale a broker.

### Quick Start:
-----------

#### Make config from template:

```bach
cat configs/broker.yaml.template >> configs/broker.yaml
```

#### Start broker:

```bach
go run cmd/broker/main.go
```

or:

```bach
make run-broker
```

build:

```bach
make build-broker
```

A Config by once file looks like this:

```yaml
addr: 'localhost:7654'
```

A Config by slaves file looks like this:

```yaml
addr: 'localhost:7654'
slaves:
    - 'localhost:7653'
    - 'localhost:7652'
    - 'localhost:7651'
````

#### If you want start with replicas

- run slave broker
- run master broker with config slaves

### Future:
-----------

- add the ability to save sent messages
- add smarter reallocation implementation
- better consider data replication properties