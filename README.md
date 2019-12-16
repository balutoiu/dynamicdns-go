# Dynamic DNS Go

A simple Dynamic DNS client written in Go.

## How to use

1. Clone the repository.

2. Build the docker container:
```
docker build -t dynamicdns-go .
```

3. Run the container:
```
docker run -d -e USERNAME=<username> -e PASSWORD=<password> -e DOMAIN=<domain> dynamicdns-go
```

Prebuilt docker images are available [here](https://hub.docker.com/r/alinbalutoiu/dynamicdns-go).

## Supported Dynamic DNS providers

Currently only Google Domains is supported.
