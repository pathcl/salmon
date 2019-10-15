package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Cert struct {
	CommonName       string    `json:"cn"`
	NotAfter         time.Time `json:"expires"`
	IssuerCommonName string    `json:"issuer"`
}

func main() {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// we will list every ingress using tls. Why? to check for expiration date and warn
	ingress, err := clientset.ExtensionsV1beta1().Ingresses("").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	var (
		c *Cert
	)
	// there must be a better way!
	// TODO: prometheus exporter
	for _, s := range ingress.Items {
		for p := range s.Spec.TLS {
			for _, h := range s.Spec.TLS[p].Hosts {
				port := h + ":443"
				c, _ = ParseRemoteCertificate(port, 10)
				log.Println(h, c.Jsonify())
			}
		}
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func getVerifiedCertificateChains(addr string, timeoutSecond time.Duration) ([][]*x509.Certificate, error) {
	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: timeoutSecond * time.Second}, "tcp", addr, nil)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	chains := conn.ConnectionState().VerifiedChains
	return chains, nil
}

func ParseRemoteCertificate(addr string, timeoutSecond int) (*Cert, error) {
	chains, err := getVerifiedCertificateChains(addr, time.Duration(timeoutSecond))
	if err != nil {
		return nil, err
	}

	var cert *Cert
	for _, chain := range chains {
		for _, crt := range chain {
			if !crt.IsCA {
				cert = &Cert{
					CommonName:       crt.Subject.CommonName,
					NotAfter:         crt.NotAfter,
					IssuerCommonName: crt.Issuer.CommonName,
				}
			}
		}
	}
	return cert, err
}

func ParseCertificateFile(certFile string) (*Cert, error) {
	b, err := ioutil.ReadFile(certFile)
	if err != nil {
		return nil, err
	}
	p, _ := pem.Decode(b)
	crt, err := x509.ParseCertificate(p.Bytes)
	if err != nil {
		return nil, err
	}
	return &Cert{
		CommonName:       crt.Subject.CommonName,
		NotAfter:         crt.NotAfter,
		IssuerCommonName: crt.Issuer.CommonName,
	}, err
}

func (cert *Cert) Jsonify() string {
	b, _ := json.Marshal(cert)
	return string(b)
}
