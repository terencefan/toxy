# Toxy

A microservice tool.

## Design Purpose

1. Simply transcode through protocols. (Binary <-> Compact)
2. Metric api data (count, response-time).
3. Work as a proxy, listening on a `port` and proxy rpc request to bankend server via `HTTP`, `Buffered` or `Framed`.
4. Gracefully downgrade specified apis while backend server is down.
5. Provide JSON http api at the same time.
6. Work with multiple backend servers through `MultiplexedProcessor`

## Config

```ini
[metric]
handler=buffered_statsd
addr=0.0.0.0:8125

[socketserver]
addr=0.0.0.0:6000
;processor=default
processor=multiplexed

;[httpserver]
;addr=0.0.0.0:8000

[service.Revenue]
thrift=/var/thrift/revenue.thrift
transport=http
addr=0.0.0.0:10010
path=/Revenue

[service.RevenueCoupon]
thrift=/var/thrift/revenue-coupon.thrift
transport=http
addr=0.0.0.0:10010
path=/RevenueCoupon

[service.RevenueOrder]
thrift=/var/thrift/revenue-order.thrift
transport=http
addr=0.0.0.0:10010
path=/RevenueOrder

;[service.Ping]
;thrift=/var/thrift/ping.thrift
;transport=http
;addr=0.0.0.0:10010
;path=/Ping

[downgrade]
apilist=/var/config/downgrade.apilist
```
