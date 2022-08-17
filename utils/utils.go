/*
Copyright © 2020 Logiq.ai <cli@logiq.ai>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import (
	"bytes"
	"crypto/tls"
	b64 "encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptrace"
	"net/textproto"
	"os"

	"github.com/spf13/viper"
)

const (
	KeyCluster    = "cluster"
	KeyPort       = "port"
	DefaultPort   = "8081"
	KeyNamespace  = "namespace"
	AuthToken     = "uitoken"
	KeyUiUser     = "uiuser"
	KeyUiPassword = "uipassword"
)

var trace = &httptrace.ClientTrace{
	/*
		GetConn: func(hostPort string) {
			fmt.Printf("GetConne(%s)\n", hostPort)
		},
		DNSStart:     func(info httptrace.DNSStartInfo) { fmt.Println("starting to look up dns", info) },
		DNSDone:      func(info httptrace.DNSDoneInfo) { fmt.Println("done looking up dns", info) },
		ConnectStart: func(network, addr string) { fmt.Println("starting tcp connection", network, addr) },
		ConnectDone:  func(network, addr string, err error) { fmt.Println("tcp connection created", network, addr, err) },
		GotConn:      func(info httptrace.GotConnInfo) { fmt.Println("connection established", info) },

	*/

	// GetConn is called before a connection is created or
	// retrieved from an idle pool. The hostPort is the
	// "host:port" of the target or proxy. GetConn is called even
	// if there's already an idle cached connection available.
	GetConn: func(hostPort string) {
		fmt.Printf("Get Conn: hostPort: %s\n", hostPort)
	},
	// GotConn is called after a successful connection is
	// obtained. There is no hook for failure to obtain a
	// connection; instead, use the error from
	// Transport.RoundTrip.
	GotConn: func(connInfo httptrace.GotConnInfo) {
		fmt.Printf("Got Conn: connInfo: %+v\n", connInfo)
	},
	// PutIdleConn is called when the connection is returned to
	// the idle pool. If err is nil, the connection was
	// successfully returned to the idle pool. If err is non-nil,
	// it describes why not. PutIdleConn is not called if
	// connection reuse is disabled via Transport.DisableKeepAlives.
	// PutIdleConn is called before the caller's Response.Body.Close
	// call returns.
	// For HTTP/2, this hook is not currently used.
	PutIdleConn: func(err error) {
		fmt.Printf("PutIdlConn: ERR: %s\n", err)
	},
	// GotFirstResponseByte is called when the first byte of the response
	// headers is available.
	GotFirstResponseByte: func() {
		fmt.Println("GotFirstResponseByte")
	},
	// Got100Continue is called if the server replies with a "100
	// Continue" response.
	Got100Continue: func() {
		fmt.Println("Got100Continue")
	},
	// Got1xxResponse is called for each 1xx informational response header
	// returned before the final non-1xx response. Got1xxResponse is called
	// for "100 Continue" responses, even if Got100Continue is also defined.
	// If it returns an error, the client request is aborted with that error value.
	Got1xxResponse: func(code int, header textproto.MIMEHeader) error {
		fmt.Printf("Got1xxResponse: code: %d header: %+v\n", code, header)
		return nil
	},
	// DNSStart is called when a DNS lookup begins.
	DNSStart: func(dnsInfo httptrace.DNSStartInfo) {
		fmt.Printf("DNS Start: dnsInfo: %+v\n", dnsInfo)
	},
	// DNSDone is called when a DNS lookup ends.
	DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
		fmt.Printf("DNS Done: dnsInfo: %+v\n", dnsInfo)
	},
	// ConnectStart is called when a new connection's Dial begins.
	// If net.Dialer.DualStack (IPv6 "Happy Eyeballs") support is
	// enabled, this may be called multiple times.
	ConnectStart: func(network, addr string) {
		fmt.Printf("Connect Start: Network Addr: %s %s\n", network, addr)
	},
	// ConnectDone is called when a new connection's Dial
	// completes. The provided err indicates whether the
	// connection completedly successfully.
	// If net.Dialer.DualStack ("Happy Eyeballs") support is
	// enabled, this may be called multiple times.
	ConnectDone: func(network, addr string, err error) {
		fmt.Printf("Connect Done: Network Addr: %s %s ERR: %s\n", network, addr, err)
	},
	// TLSHandshakeStart is called when the TLS handshake is started. When
	// connecting to an HTTPS site via an HTTP proxy, the handshake happens
	// after the CONNECT request is processed by the proxy.
	TLSHandshakeStart: func() {
		fmt.Println("TLSHandshakeStart")
	},
	// TLSHandshakeDone is called after the TLS handshake with either the
	// successful handshake's connection state, or a non-nil error on handshake
	// failure.
	TLSHandshakeDone: func(connState tls.ConnectionState, err error) {
		fmt.Printf("TLSHandshakeDone: connState: %+v ERR: %s\n", connState, err)
	},
	// WroteHeaderField is called after the Transport has written
	// each request header. At the time of this call the values
	// might be buffered and not yet written to the network.
	WroteHeaderField: func(key string, value []string) {
		fmt.Printf("WroteHeaderField: key: %s val: %s\n", key, value)
	},
	// WroteHeaders is called after the Transport has written
	// all request headers.
	WroteHeaders: func() {
		fmt.Println("WroteHeaders")
	},
	// Wait100Continue is called if the Request specified
	// "Expect: 100-continue" and the Transport has written the
	// request headers but is waiting for "100 Continue" from the
	// server before writing the request body.
	Wait100Continue: func() {
		fmt.Println("Wait100Continue")
	},
	// WroteRequest is called with the result of writing the
	// request and any body. It may be called multiple times
	// in the case of retried requests.
	WroteRequest: func(info httptrace.WroteRequestInfo) {
		fmt.Printf("WroteRequest: %+v\n", info)
	},
}

func GetClusterUrl() string {
	var cluster string
	if FlagCluster != "" {
		cluster = FlagCluster
	} else {
		cluster = viper.GetString(KeyCluster)
	}
	port := viper.GetString(KeyPort)
	return fmt.Sprintf("%s:%s", cluster, port)
}

func GetDefaultNamespace() string {
	if FlagNamespace != "" {
		return FlagNamespace
	}
	ns := viper.GetString(KeyNamespace)
	return ns
}

func GetUIUser() string {
	uiEncodedUser := viper.GetString(KeyUiUser)
	uiUser, _ := b64.StdEncoding.DecodeString(uiEncodedUser)
	return string(uiUser)
}

func GetUIPass() string {
	uiEncodedPass := viper.GetString(KeyUiPassword)
	uiPass, _ := b64.StdEncoding.DecodeString(uiEncodedPass)
	return string(uiPass)
}

func CheckMesgErr(resp map[string]interface{}, orgi string) {
	errMesg := "Internal Server Error"
	if mesg, ok := resp["message"]; ok {
		if mesg == errMesg {
			fmt.Println("ERR> From", orgi, ", mesg:", errMesg, ", exit()")
			os.Exit(-1)
		}
	}
}

func AddNetTrace(req *http.Request) *http.Request {

	if FlagNetTrace {
		ctx := httptrace.WithClientTrace(req.Context(), trace)
		req = req.WithContext(ctx)
	}
	return req
}

func CreateHttpRequest(method string, uri string, payload *bytes.Buffer) (*http.Request, error) {
	req, err := http.NewRequest(method, uri, payload)
	if err != nil {
		return nil, err
	}
	if api_key := viper.GetString(AuthToken); api_key != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Key %s", api_key))
		req.Header.Add("Content-Type", "application/json")
	}
	return req, err
}
