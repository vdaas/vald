module github.com/vdaas/vald/example/client

go 1.18

replace (
	github.com/envoyproxy/protoc-gen-validate => github.com/envoyproxy/protoc-gen-validate v0.6.12
	github.com/goccy/go-json => github.com/goccy/go-json v0.9.11
	github.com/golang/protobuf => github.com/golang/protobuf v1.5.2
	github.com/kpango/glg => github.com/kpango/glg v1.6.13
	github.com/pkg/sftp => github.com/pkg/sftp v1.13.5
	golang.org/x/crypto => golang.org/x/crypto v0.0.0-20220926161630-eccd6366d1be
	golang.org/x/net => golang.org/x/net v0.0.0-20221002022538-bcab6841153b
	golang.org/x/text => golang.org/x/text v0.3.7
	google.golang.org/genproto => google.golang.org/genproto v0.0.0-20220930163606-c98284e70a91
	google.golang.org/grpc => google.golang.org/grpc v1.49.0
	google.golang.org/protobuf => google.golang.org/protobuf v1.28.1
	gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 => gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/kpango/fuid v0.0.0-20220209050620-e5987ba1ea5e
	github.com/kpango/glg v1.6.10
	github.com/vdaas/vald-client-go v1.5.6
	gonum.org/v1/hdf5 v0.0.0-20210714002203-8c5d23bc6946
	google.golang.org/grpc v1.49.0
)

require (
	github.com/envoyproxy/protoc-gen-validate v0.1.0 // indirect
	github.com/goccy/go-json v0.9.11 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/kpango/fastime v1.1.4 // indirect
	golang.org/x/net v0.0.0-20220909164309-bea034e7d591 // indirect
	golang.org/x/sys v0.0.0-20220908150016-7ac13a9a928d // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220926220553-6981cbe3cfce // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)
