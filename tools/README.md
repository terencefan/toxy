# Tools

## Ting

Simply send the ping reqeust to the server, used for health-checking.

```
# -t Transport, default: TSocket
# -w TransportWrapper, default: TTransport
# -p Protocol, default: TBinaryProtocol
# -s Service, optional
ting -h 127.0.0.1 -p 2020 -t socket -w buffered -s Revenue
```

## Tenchmark

A benchmark tool / framework based on thrift protocol.

```
# -t Transport, default: TSocket
# -w TransportWrapper, default: TTransport
# -p Protocol, default: TBinaryProtocol
tenchmark -h 127.0.0.1 -p 2020 -t socket -w buffered
```
