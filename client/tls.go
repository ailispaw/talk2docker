package client

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"path/filepath"

	"github.com/codegangsta/cli"
)

func getTLSConfig(ctx *cli.Context) (*tls.Config, error) {
	var tlsConfig tls.Config

	certPath := ctx.GlobalString("tls")
	if certPath == "" {
		return nil, nil
	}

	tlscacert := filepath.Join(certPath, "ca.pem")
	tlscert := filepath.Join(certPath, "cert.pem")
	tlskey := filepath.Join(certPath, "key.pem")

	tlsConfig.InsecureSkipVerify = ctx.GlobalBool("insecure-tls")

	if !tlsConfig.InsecureSkipVerify {
		certPool := x509.NewCertPool()
		file, err := ioutil.ReadFile(tlscacert)
		if err != nil {
			return nil, err
		}
		certPool.AppendCertsFromPEM(file)
		tlsConfig.RootCAs = certPool
	}

	cert, err := tls.LoadX509KeyPair(tlscert, tlskey)
	if err != nil {
		return nil, err
	}
	tlsConfig.Certificates = []tls.Certificate{cert}
	tlsConfig.MinVersion = tls.VersionTLS10

	return &tlsConfig, nil
}
