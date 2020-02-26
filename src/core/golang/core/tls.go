package core

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
)

func (t TLSConfig) Config() *tls.Config {
	tlsCfg := &tls.Config{}
	tlsCfg.RootCAs = x509.NewCertPool()

	ca, err := ioutil.ReadFile(t.TLSRootCaCert)
	if err != nil {
		Error("Failed to read CA cert: " + err.Error())
	}
	tlsCfg.RootCAs.AppendCertsFromPEM(ca)
	return tlsCfg
}
