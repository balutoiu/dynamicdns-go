# Dynamic DNS Go

[![Build DynamicDNS Go](https://github.com/balutoiu/dynamicdns-go/actions/workflows/build-dynamicdns-go.yaml/badge.svg)](https://github.com/balutoiu/dynamicdns-go/actions/workflows/build-dynamicdns-go.yaml)

Simple Dynamic DNS client written in Go.

## How to use

1. Clone the repository.

1. Build the container image:

    ```bash
    docker build -t balutoiu/dynamicdns-go .
    ```

1. Create the config file ([sample here](/testdata/config.yaml)):

    ```yaml
    googledomains:
      username: user
      password: pass
      domain: test.example.com
    ```

1. Run the container:

    ```bash
    docker run -d \
        -v ./config.yaml:/config.yaml \
        -e DNS_PROVIDER=googledomains \
        balutoiu/dynamicdns-go
    ```

Prebuilt docker images are available [here](https://github.com/balutoiu/dynamicdns-go/pkgs/container/dynamicdns-go).

## Environment variables

- `DNS_PROVIDER` - the DNS provider to be used (defaults to `googledomains`)

- `SLEEP_INTERVAL` - the amount of time to wait between updates,
example: `1s`, `1m`, `1h` etc. (defaults to `1h0m0s`)

## Supported Dynamic DNS providers

### Google domains

Example of `config.yaml` file:

```yaml
googledomains:
  username: user
  password: pass
  domain: test.example.com
```

Start the container with:

```bash
docker run -d --restart=unless-stopped \
    -v ./config.yaml:/config.yaml \
    -e DNS_PROVIDER=googledomains \
    alinbalutoiu/dynamicdns-go
```

### Mail-in-a-box

Example of `config.yaml` file:

```yaml
mailinabox:
  username: user
  password: pass
  domain: www.example.com
  api_url: https://box.mailinabox.com
```

Start the container with:

```bash
docker run -d --restart=unless-stopped \
    -v ./config.yaml:/config.yaml \
    -e DNS_PROVIDER=mailinabox \
    alinbalutoiu/dynamicdns-go
```

### OVH domains

Example of `config.yaml` file:

```yaml
ovhdomains:
  application_key: app-key
  application_secret: app-secret
  consumer_key: consumer-key
  zone_name: example
  sub_domain: www
```

Start the container with:

```bash
docker run -d --restart=unless-stopped \
    -v ./config.yaml:/config.yaml \
    -e DNS_PROVIDER=ovhdomains \
    alinbalutoiu/dynamicdns-go
```
