package client

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
)

func getTLSConfig(host *Host) (*tls.Config, error) {
	var tlsConfig tls.Config

	if !host.TLS {
		return nil, nil
	}

	tlsConfig.InsecureSkipVerify = !host.TLSVerify

	if host.TLSVerify {
		certPool := x509.NewCertPool()
		file, err := ioutil.ReadFile(host.TLSCaCert)
		if err != nil {
			return nil, err
		}
		certPool.AppendCertsFromPEM(file)
		tlsConfig.RootCAs = certPool
	}

	cert, err := tls.LoadX509KeyPair(host.TLSCert, host.TLSKey)
	if err != nil {
		return nil, err
	}
	tlsConfig.Certificates = []tls.Certificate{cert}
	tlsConfig.MinVersion = tls.VersionTLS10

	return &tlsConfig, nil
}
