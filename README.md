# :onion: `toc` - _The Onion Crawler_

The Onion Crawler is a simple, straightforward web crawler designed to traverse `.onion` sites.

## The TOR Proxy

To access the TOR network, you must utilize the `tor` daemon with a SOCKS5 proxy.

This will allow `toc` to programatically proxy all traffic through the TOR network.

```sh
$ brew install tor
$ brew services start tor
```

## Getting Started

```sh
$ git clone https://github.com/aedenmurray/toc && cd toc
$ go run main.go 
```

## Usage

- `-url` - The initial URL to start the crawler. (default: `https://raw.githubusercontent.com/fastfire/deepdarkCTI/main/forum.md`)
- `-shost` - The SOCKS5 host (default: `localhost`)
- `-sport` - The SOCKS5 port (default: `9050`)
