package client

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
)

func getTLSConfig(hostConfig *HostConfig) (*tls.Config, error) {
	var tlsConfig tls.Config

	if !hostConfig.TLS {
		return nil, nil
	}

	tlsConfig.InsecureSkipVerify = !hostConfig.TLSVerufy

	if hostConfig.TLSVerufy {
		certPool := x509.NewCertPool()
		file, err := ioutil.ReadFile(hostConfig.TLSCaCert)
		if err != nil {
			return nil, err
		}
		certPool.AppendCertsFromPEM(file)
		tlsConfig.RootCAs = certPool
	}

	cert, err := tls.LoadX509KeyPair(hostConfig.TLSCert, hostConfig.TLSKey)
	if err != nil {
		return nil, err
	}
	tlsConfig.Certificates = []tls.Certificate{cert}
	tlsConfig.MinVersion = tls.VersionTLS10

	return &tlsConfig, nil
}
