package engine

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"

	"github.com/admpub/resty/v2"
	"github.com/webx-top/restyclient"
)

var defaultClient = resty.New().SetRetryCount(restyclient.DefaultMaxRetryCount).
	SetTimeout(restyclient.DefaultTimeout).
	SetRedirectPolicy(restyclient.DefaultRedirectPolicy)

func init() {
	restyclient.InitRestyHook(defaultClient)
}

// NewAPIClient("cert.pem", "key.pem")
func NewAPIClient(certPEMBlock, keyPEMBlock []byte) (*APIClient, error) {
	var rclient *resty.Client
	if len(certPEMBlock) > 0 && len(keyPEMBlock) > 0 {
		// create certificate
		cert, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
		if err != nil {
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(certPEMBlock)

		// Create a HTTPS client and supply the created CA pool and certificate
		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs:      caCertPool,
					Certificates: []tls.Certificate{cert},
				},
			},
		}
		rclient = resty.NewWithClient(client)
		rclient.SetRetryCount(restyclient.DefaultMaxRetryCount).
			SetTimeout(restyclient.DefaultTimeout).
			SetRedirectPolicy(restyclient.DefaultRedirectPolicy)
		restyclient.InitRestyHook(rclient)
	} else {
		rclient = defaultClient
	}
	return &APIClient{client: rclient}, nil
}

type APIClient struct {
	client *resty.Client
}

func (a *APIClient) Post(url string, data interface{}) error {
	resp, err := a.client.R().SetBody(data).Post(url)
	if err != nil {
		return err
	}
	if resp.IsError() {
		err = fmt.Errorf(`%s`, resp.Body())
	}
	return err
}

// documention: https://docs.docker.com/engine/api/v1.43/#tag/Exec/operation/ContainerExec
// API: /containers/{id}/exec
type RequestDockerExec struct {
	AttachStdin  bool     `json:",omitempty"`
	AttachStdout bool     `json:",omitempty"`
	AttachStderr bool     `json:",omitempty"`
	DetachKeys   string   `json:",omitempty"` //"ctrl-p,ctrl-q",
	Tty          bool     `json:",omitempty"`
	Cmd          []string `json:",omitempty"` // ["date"],
	Env          []string `json:",omitempty"` // ["FOO=bar","BAZ=quux"]
}
