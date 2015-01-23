/*!
 * Copyright 2014 Sam Alba <sam.alba@gmail.com>
 * Licensed under the Apache License, Version 2.0
 * github.com/samalba/dockerclient/LICENSE
 *
 * github.com/samalba/dockerclient/dockerclient.go
 */

package api

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	API_VERSION = "1.15"
)

type DockerClient struct {
	URL           *url.URL
	HTTPClient    *http.Client
	TLSConfig     *tls.Config
	monitorEvents int32
	out           io.Writer
}

type Error struct {
	StatusCode int
	Status     string
	msg        string
}

func (e Error) Error() string {
	return fmt.Sprintf("Error response from daemon: %s", strings.TrimSpace(e.msg))
}

func newHTTPClient(u *url.URL, tlsConfig *tls.Config, timeout time.Duration) *http.Client {
	httpTransport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	switch u.Scheme {
	default:
		httpTransport.Dial = func(proto, addr string) (net.Conn, error) {
			return net.DialTimeout(proto, addr, timeout)
		}
	case "unix":
		socketPath := u.Path
		unixDial := func(proto, addr string) (net.Conn, error) {
			return net.DialTimeout("unix", socketPath, timeout)
		}
		httpTransport.Dial = unixDial
		// Override the main URL object so the HTTP lib won't complain
		u.Scheme = "http"
		u.Host = "unix.sock"
		u.Path = ""
	}
	return &http.Client{Transport: httpTransport}
}

func NewDockerClient(daemonUrl string, tlsConfig *tls.Config, timeout time.Duration, out io.Writer) (*DockerClient, error) {
	u, err := url.Parse(daemonUrl)
	if err != nil {
		return nil, err
	}
	if u.Scheme == "tcp" {
		if tlsConfig == nil {
			u.Scheme = "http"
		} else {
			u.Scheme = "https"
		}
	}
	httpClient := newHTTPClient(u, tlsConfig, timeout)
	return &DockerClient{u, httpClient, tlsConfig, 0, out}, nil
}

func (client *DockerClient) doRequest(method string, path string, body []byte, headers map[string]string) ([]byte, error) {
	b := bytes.NewBuffer(body)
	req, err := http.NewRequest(method, client.URL.String()+path, b)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	if headers != nil {
		for header, value := range headers {
			req.Header.Add(header, value)
		}
	}
	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		if !strings.Contains(err.Error(), "connection refused") && client.TLSConfig == nil {
			return nil, fmt.Errorf("%v. Are you trying to connect to a TLS-enabled daemon without TLS?", err)
		}
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return nil, Error{StatusCode: resp.StatusCode, Status: resp.Status, msg: string(data)}
	}
	return data, nil
}

func (client *DockerClient) doStreamRequest(method string, path string, in io.Reader, headers map[string]string, quiet bool) (string, error) {
	if (method == "POST" || method == "PUT") && in == nil {
		in = bytes.NewReader([]byte{})
	}

	req, err := http.NewRequest(method, client.URL.String()+path, in)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	if headers != nil {
		for header, value := range headers {
			req.Header.Add(header, value)
		}
	}
	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		if !strings.Contains(err.Error(), "connection refused") && client.TLSConfig == nil {
			return "", fmt.Errorf("%v. Are you trying to connect to a TLS-enabled daemon without TLS?", err)
		}
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		if len(body) == 0 {
			return "", fmt.Errorf("Error :%s", http.StatusText(resp.StatusCode))
		}
		return "", fmt.Errorf("Error: %s", bytes.TrimSpace(body))
	}

	mimetype, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err == nil && mimetype == "application/json" {
		out := client.out
		if quiet {
			f, err := os.Open(os.DevNull)
			if err == nil {
				out = f
			} else {
				quiet = false
			}
		}
		message, err := displayJSONMessagesStream(resp.Body, out)
		if quiet && (message != "") {
			fmt.Fprintf(client.out, "%s", message)
		}
		return message, err
	}

	_, err = io.Copy(client.out, resp.Body)
	return "", err
}
