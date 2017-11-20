# Toxy (work in progress)

A thrift microservice proxy.

It can receive thrift api from multiple clients through `multiplexed protocol`, and transform (`binary protocol`, `buffered socket transport`) to **any** combinations of `protocol`, `transport`, `transport wrapper`.

## Other Design Purpose

1. Gracefully shutdown. (by return a application exception)
2. Metric api statistic data. (statsd protocol)
3. Gracefully downgrade specified apis while backend service is down. (return empty value)
4. Provide a socket pool to backend service which uses `socket transport`

Nice to have:

1. tracking.
2. parallel.

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
