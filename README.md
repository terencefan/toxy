# Toxy (work in progress)

A thrift microservice proxy.

**It's also a easy-to-use thrift library written in go**

[documentation]()

It can receive thrift api from multiple clients through `multiplexed protocol`, and transform `(binary protocol, buffered socket transport)` to **any** combinations of `(protocol, transport, transport wrapper)`.

## Current achieved goals

1. Gracefully shutdown. (by return a application exception)
2. Metric api statistic data. (buffered statsd)

### Protocols

* binary protocol

### Transports

* socket
* http (over tcp)
* memory

### Transport Wrappers

* buffered
* framed

## Todo list

1. Gracefully downgrade specified apis while backend service is down. (return empty value)
2. Provide a socket pool to backend service which uses `socket transport`

### Protocols

* compact protocol
* json protocol

### Transports

* tls socket
* https (over tls)

## Config

```ini
[metric]
handler=buffered_statsd
addr=0.0.0.0:8125

[socketserver]
addr=0.0.0.0:6000
;processor=default
processor=multiplexed

[service.Revenue]
transport=http
addr=0.0.0.0:10010
path=/Revenue

[service.RevenueCoupon]
transport=http
addr=0.0.0.0:10010
path=/RevenueCoupon

[service.RevenueOrder]
transport=http
addr=0.0.0.0:10010
path=/RevenueOrder

[downgrade]
apilist=/var/config/downgrade.apilist
```
