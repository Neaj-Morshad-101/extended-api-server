module github.com/Neaj-Morshad-101/extended-api-server

go 1.20

require (
	github.com/gorilla/mux v1.8.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/afero v1.9.5
	//k8s.io/client-go v0.27.2
	k8s.io/client-go v9.0.0+incompatible
)

require k8s.io/klog/v2 v2.100.1

require (
	github.com/go-logr/logr v1.2.0 // indirect
	golang.org/x/text v0.8.0 // indirect
)
