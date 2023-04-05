module github.com/vdaas/vald/example/client

go 1.20

replace (
	github.com/envoyproxy/protoc-gen-validate => github.com/envoyproxy/protoc-gen-validate v0.10.1
	github.com/goccy/go-json => github.com/goccy/go-json v0.10.2
	github.com/golang/protobuf => github.com/golang/protobuf v1.5.3
	github.com/kpango/glg => github.com/kpango/glg v1.6.15
	github.com/pkg/sftp => github.com/pkg/sftp v1.13.5
	golang.org/x/crypto => golang.org/x/crypto v0.7.0
	golang.org/x/net => golang.org/x/net v0.8.0
	golang.org/x/text => golang.org/x/text v0.8.0
	google.golang.org/genproto => google.golang.org/genproto v0.0.0-20230403163135-c38d8f061ccd
	google.golang.org/grpc => google.golang.org/grpc v1.54.0
	google.golang.org/protobuf => google.golang.org/protobuf v1.30.0
	gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 => gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/kpango/fuid v0.0.0-20221203053508-503b5ad89aa1
	github.com/kpango/glg v1.6.14
	github.com/vdaas/vald-client-go v1.7.4
	gonum.org/v1/hdf5 v0.0.0-20210714002203-8c5d23bc6946
	google.golang.org/grpc v1.54.0
)

require (
	github.com/envoyproxy/protoc-gen-validate v0.9.1 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/kpango/fastime v1.1.9 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/genproto v0.0.0-20230110181048-76db0878b65f // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)
