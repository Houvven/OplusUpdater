package updater

import (
	"crypto/aes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/deatil/go-cryptobin/cryptobin/rsa"
	"golang.org/x/net/context"
	"golang.org/x/net/proxy"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func RandomIv() ([]byte, error) {
	iv := make([]byte, aes.BlockSize)
	_, err := io.ReadFull(rand.Reader, iv)
	return iv, err
}

func RandomKey() ([]byte, error) {
	key := make([]byte, aes.BlockSize*2)
	_, err := io.ReadFull(rand.Reader, key)
	return key, err
}

func GenerateProtectedVersion() string {
	return strconv.FormatInt(time.Now().Add(time.Hour*24).UnixNano(), 10)
}

func GenerateProtectedKey(key []byte, pubKey []byte) (string, error) {
	encrypt := rsa.New().
		FromString(base64.StdEncoding.EncodeToString(key)).
		FromPublicKey(pubKey).
		EncryptOAEP()
	return encrypt.ToBase64String(), encrypt.Error()
}

func GetDefaultDeviceId() string {
	return strings.Repeat("0", 64)
}

func ParseTransportFromProxyStr(p string) (*http.Transport, error) {
	if p == "" {
		return &http.Transport{}, nil
	}

	parsedURL, err := url.Parse(p)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return nil, fmt.Errorf("invalid proxy URL: %s", p)
	}

	switch parsedURL.Scheme {
	case "socks":
		auth := func() *proxy.Auth {
			if parsedURL.User != nil {
				pass, _ := parsedURL.User.Password()
				return &proxy.Auth{User: parsedURL.User.Username(), Password: pass}
			}
			return nil
		}()
		dialer, err := proxy.SOCKS5("tcp", parsedURL.Host, auth, proxy.Direct)
		if err != nil {
			return nil, err
		}
		return &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			},
		}, nil
	case "http", "https":
		return &http.Transport{Proxy: http.ProxyURL(parsedURL)}, nil
	default:
		return nil, fmt.Errorf("unsupported proxy scheme: %s", parsedURL.Scheme)
	}
}
