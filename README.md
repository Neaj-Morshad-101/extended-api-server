# extended-api-server

### Kubernetes Extended API Server using net/http library


-------------------------------------

```console
## Terminal 1
`$ go run apiserver/main.go`

listening on 127.0.0.1:8443

## Terminal 2
`export APISERVER_ADDR=127.0.0.1:8443`

`$ curl -k https://${APISERVER_ADDR}/core/pods`

Resource: pods

```


-------------------------------------

## Terminal 3
`$ go run database-apiserver/main.go`

listening on 127.0.0.2:8443

## Terminal 2
`export EAS_ADDR=127.0.0.2:8443`

`$ curl -k https://${EAS_ADDR}/database/postgres`

Resource: postgres


-------------------------------------

## Terminal 1
`$ go run apiserver/main.go --send-proxy-request=true`

listening on 127.0.0.1:8443
forwarding request to https://127.0.0.2:8443/database/postgres

## Terminal 3
`$ go run database-apiserver/main.go --receive-proxy-request=true`

listening on 127.0.0.2:8443

## Terminal 2
`$ cd cd /tmp/extended-api-server/`
`$ curl -k https://${APISERVER_ADDR}/core/pods`

Resource: pods

`$ curl -k https://${APISERVER_ADDR}/database/postgres`

Resource: postgres requested by user[X-Remote-User]=

`$ curl https://${APISERVER_ADDR}/database/postgres \
--cacert ./apiserver-ca.crt \
--cert ./apiserver-john.crt \
--key ./apiserver-john.key`

Resource: postgres requested by user[X-Remote-User]=john

`$ curl https://${EAS_ADDR}/database/postgres \
--cacert ./database-ca.crt \
--cert ./apiserver-john.crt \
--key ./apiserver-john.key`

Resource: postgres requested by user[Client-Cert-CN]=john

`$ curl https://${EAS_ADDR}/database/postgres \
--cacert ./database-ca.crt \
--cert ./database-jane.crt \
--key ./database-jane.key`

curl: (35) error:14094412:SSL routines:ssl3_read_bytes:sslv3 alert bad certificate

`$ curl -k https://${EAS_ADDR}/database/postgres`

Resource: postgres requested by user[-]=system:anonymous

-------------------------------------

## Resources: 
https://www.youtube.com/watch?v=pTIwy6fpxwY&t=2s&ab_channel=AppsCodeInc.
https://github.com/tamalsaha/DIY-k8s-extended-apiserver

## Kubernetes External API Server

https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/

https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/apiserver-aggregation/

https://kubernetes.io/docs/tasks/extend-kubernetes/configure-aggregation-layer/

https://medium.com/@vanSadhu/kubernetes-api-aggregation-setup-nuts-bolts-733fef22a504

-------------------------------------


-------------------------------------
