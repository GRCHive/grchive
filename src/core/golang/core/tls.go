package core

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"os"
)

func (t TLSConfig) Config() *tls.Config {
	tlsCfg := &tls.Config{}
	tlsCfg.RootCAs = x509.NewCertPool()

	ca, err := ioutil.ReadFile(t.TLSRootCaCert)
	if err != nil {
		Warning("Failed to read CA cert: " + err.Error())
		return nil
	}
	tlsCfg.RootCAs.AppendCertsFromPEM(ca)

	if _, ok := os.LookupEnv("ALLOW_TLS_INSECURE"); ok {
		tlsCfg.InsecureSkipVerify = true
	}

	return tlsCfg
}
