package hdsclient

import (
	"crypto/ecdsa"
	"net/http"

	httpcmn "github.com/aukilabs/hagall-common/http"
)

type ClientOpts func(*Client)

func WithHagallEndpoint(v string) ClientOpts {
	return func(c *Client) {
		c.HagallEndpoint = httpcmn.NormalizeEndpoint(v)
	}
}

func WithHDSEndpoint(v string) ClientOpts {
	return func(c *Client) {
		c.HDSEndpoint = httpcmn.NormalizeEndpoint(v)
	}
}

func WithEncoder(v Encoder) ClientOpts {
	return func(c *Client) {
		c.Encode = v
	}
}

func WithDecoder(v Decoder) ClientOpts {
	return func(c *Client) {
		c.Decode = v
	}
}

func WithTransport(v http.RoundTripper) ClientOpts {
	return func(c *Client) {
		c.Transport = v
	}
}

func WithSecret(v string) ClientOpts {
	return func(c *Client) {
		c.secret = v
	}
}

func WithRemoteAddr(v string) ClientOpts {
	return func(c *Client) {
		c.RemoteAddr = v
	}
}

func WithServerID(v string) ClientOpts {
	return func(c *Client) {
		c.serverID = v
	}
}

func WithClientID(v string) ClientOpts {
	return func(c *Client) {
		c.clientID = v
	}
}

func WithPrivateKey(v *ecdsa.PrivateKey) ClientOpts {
	return func(c *Client) {
		c.privateKey = v
	}
}
