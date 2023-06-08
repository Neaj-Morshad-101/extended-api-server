package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"k8s.io/klog/v2"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/Neaj-Morshad-101/extended-api-server/lib/certstore"
	"github.com/Neaj-Morshad-101/extended-api-server/lib/server"
	"github.com/gorilla/mux"
	"github.com/spf13/afero"
	"k8s.io/client-go/util/cert"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}

func main() {
	var proxy = false
	flag.BoolVar(&proxy, "send-proxy-request", proxy, "forward requests to database extended apiserver")
	flag.Parse()

	fs := afero.NewOsFs()
	store, err := certstore.NewCertStore(fs, "/tmp/extended-api-server")
	if err != nil {
		klog.Info()
		log.Fatalln(err)
		log.Println("Line 34")
	}
	err = store.InitCA("apiserver")
	if err != nil {
		klog.Info()
		log.Fatalln(err)
		log.Println("Line")
	}
	serverCert, serverKey, err := store.NewServerCertPair(cert.AltNames{
		IPs: []net.IP{net.ParseIP("127.0.0.1")},
	})
	if err != nil {
		klog.Info()
		log.Fatalln(err)
	}
	err = store.Write("tls", serverCert, serverKey)
	if err != nil {
		klog.Info()
		log.Fatalln(err)
	}
	clientCert, clientKey, err := store.NewClientCertPair(cert.AltNames{
		DNSNames: []string{"john"},
	})
	if err != nil {
		klog.Info()
		log.Fatalln(err)
	}
	err = store.Write("john", clientCert, clientKey)
	if err != nil {
		log.Fatalln(err)
	}

	// ---------------------------------------------------------
	rhStore, err := certstore.NewCertStore(fs, "/tmp/extended-api-server")
	if err != nil {
		klog.Info()
		log.Fatalln(err)
	}
	err = rhStore.InitCA("requestheader")
	if err != nil {
		klog.Info()
		log.Fatalln(err)
	}

	rhClientCert, rhClientKey, err := rhStore.NewClientCertPair(cert.AltNames{
		DNSNames: []string{"apiserver"}, // because apiserver is making the calls to database eas
	})

	if err != nil {
		klog.Info()
		log.Fatalln(err)
	}

	err = rhStore.Write("apiserver", rhClientCert, rhClientKey)

	if err != nil {
		klog.Info()
		log.Fatalln(err)
	}

	rhCert, err := tls.LoadX509KeyPair(rhStore.CertFile("apiserver"), rhStore.KeyFile("apiserver"))

	if err != nil {
		klog.Info()
		log.Fatalln(err)
	}
	// ---------------------------------

	// -----------------------------
	easCACertPool := x509.NewCertPool()

	if proxy {
		easStore, err := certstore.NewCertStore(fs, "/tmp/extended-api-server")
		if err != nil {
			klog.Info(err)
		}
		err = easStore.LoadCA("database")
		if err != nil {
			klog.Info(err)
		}

		easCACertPool.AppendCertsFromPEM(easStore.CACertBytes())
	}

	// -----------------------------

	cfg := server.Config{
		Address: "127.0.0.1:8443",
		CACertFiles: []string{
			store.CertFile("ca"),
		},
		CertFile: store.CertFile("tls"),
		KeyFile:  store.KeyFile("tls"),
	}
	srv := server.NewGenericServer(cfg)

	r := mux.NewRouter()
	r.HandleFunc("/core/{resource}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Resource: %v\n", vars["resource"])
	})
	if proxy {
		r.HandleFunc("/database/{resource}", func(w http.ResponseWriter, r *http.Request) {
			tr := &http.Transport{
				MaxIdleConnsPerHost: 10,
				TLSClientConfig: &tls.Config{
					Certificates: []tls.Certificate{rhCert},
					RootCAs:      easCACertPool,
				},
			}
			client := http.Client{
				Transport: tr,
				Timeout:   time.Duration(30 * time.Second),
			}

			u := *r.URL
			u.Scheme = "https"
			u.Host = "127.0.0.2:8443"
			fmt.Printf("forwarding request to %v\n", u.String())

			req, _ := http.NewRequest(r.Method, u.String(), nil)
			if len(r.TLS.PeerCertificates) > 0 {
				req.Header.Set("X-Remote-User", r.TLS.PeerCertificates[0].Subject.CommonName)
			}

			resp, err := client.Do(req)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "error: %v\n", err.Error())
				return
			}
			defer resp.Body.Close()

			w.WriteHeader(http.StatusOK)
			io.Copy(w, resp.Body)
		})
	}
	r.HandleFunc("/", handler)
	srv.ListenAndServe(r)
}
