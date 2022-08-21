# Jellyfish in-memory broker message.
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
go run cmd/broker.go
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