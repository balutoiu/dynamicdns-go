# Dynamic DNS Go

A simple Dynamic DNS client written in Go.

## How to use

1. Clone the repository.

2. Build the docker container:
```
docker build -t alinbalutoiu/dynamicdns-go .
```

3. Create the config file ([sample here](/config.yaml)):
```
googledomains:
  username: user
  password: pass
  domain: test.example.com
```

4. Run the container:
```
docker run -d -v `pwd`/config.yaml:/config.yaml -e DNS_PROVIDER=googledomains alinbalutoiu/dynamicdns-go
```

Prebuilt docker images are available [here](https://hub.docker.com/r/alinbalutoiu/dynamicdns-go).

## Environment variables

- `DNS_PROVIDER` - the DNS provider to be used (defaults to `googledomains`)

- `SLEEP_INTERVAL` - the amount of time to wait between updates,
example: `1s`, `1m`, `1h` etc. (defaults to `1h0m0s`)

## Supported Dynamic DNS providers

### Google domains

Example of `config.yaml` file:
```
googledomains:
  username: user
  password: pass
  domain: test.example.com
```

Start the container with:
```
docker run -d --restart=unless-stopped \
    -v `pwd`/config.yaml:/config.yaml \
    -e DNS_PROVIDER=googledomains \
    alinbalutoiu/dynamicdns-go
```

### Mail-in-a-box

Example of `config.yaml` file:
```
mailinabox:
  username: user
  password: pass
  domain: www.example.com
  api_url: https://box.mailinabox.com
```

Start the container with:
```
docker run -d --restart=unless-stopped \
    -v `pwd`/config.yaml:/config.yaml \
    -e DNS_PROVIDER=mailinabox \
    alinbalutoiu/dynamicdns-go
```
