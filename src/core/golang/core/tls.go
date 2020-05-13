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
	} else {
		tlsCfg.RootCAs.AppendCertsFromPEM(ca)
	}

	_, ok := os.LookupEnv("ALLOW_TLS_INSECURE")
	if ok {
		Warning("Skipping TLS verification.")
		tlsCfg.InsecureSkipVerify = true
	}

	return tlsCfg
}
