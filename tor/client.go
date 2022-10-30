package tor

import (
	"fmt"
	"net/http"
	"net/url"
)

type Client struct {
	SocksHost string
	SocksPort string
	http.Client
}

func CreateClient(socksHost string, socksPort string) (*Client, error) {
	socksProxy := fmt.Sprintf("socks5://%s:%s", socksHost, socksPort)
	proxyURL, err := url.Parse(socksProxy)
	if err != nil {
		return nil, err
	}

	return &Client{
		SocksHost: socksHost,
		SocksPort: socksPort,
		Client: http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
		},
	}, nil
}
