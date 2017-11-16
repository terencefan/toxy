# Toxy

A microservice tool.

## Design Purpose

1. Simply transcode through protocols. (Binary <-> Compact)
2. Metric api data (count, response-time).
3. Work as a proxy, listening on a `port` and proxy rpc request to bankend server via `HTTP`, `Buffered` or `Framed`.
4. Gracefully downgrade specified apis while backend server is down.
5. Provide JSON http api at the same time.
6. Work with multiple backend servers through `MultiplexedProcessor`
