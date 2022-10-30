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

By default, `toc` will print out the `depth`, `url`, & `title` of any sites that it sucessfully visits.

```sh
$ git clone https://github.com/aedenmurray/toc && cd toc
$ go run main.go 

the onion cralwer - https://github.com/aedenmurray/toc
------------------------------------------------------

0    https://github.com/fastfire/deepdarkCTI/blob/main/forum.md (deepdarkCTI/forum.md at main · fastfire/deepdarkCTI · GitHub)
1    http://xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx.onion (Threat Actors | Onion Forums)
1    http://xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx.onion/login (Forum)
1    http://xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx.onion (CryptBB)

...etc
```

## Usage

- `-url` - Initial URL to crawl: _`https://github.com/fastfire/deepdarkCTI/blob/main/forum.md`_
- `-shost` - The SOCKS5 host: _`localhost`_
- `-sport` - The SOCKS5 port: _`9050`_
