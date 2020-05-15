package internal

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// HttpsClient is a simple function loading CA certificate
// and create an HTTPS client from it.
// As is, it does not support mTLS connectivity.
func HttpsClient(caFile string) *http.Client {
	caCert := caFile
	if exists(caFile) {
		caCert, err := ioutil.ReadFile(caFile)
		if err != nil {
			log.Printf("unable to load certificate from file [%s]. got [%v]", caCert, err)
		}
	}
	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM([]byte(caCert)); !ok {
		log.Printf("unable to configure server certificate authority with certificate [%s]\n", caCert)
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}
	return client
}

func exists(file string) bool {
	info, err := os.Stat(file)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
