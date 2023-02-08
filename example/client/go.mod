module github.com/vdaas/vald/example/client

go 1.19

replace (
	github.com/envoyproxy/protoc-gen-validate => github.com/envoyproxy/protoc-gen-validate v0.9.1
	github.com/goccy/go-json => github.com/goccy/go-json v0.10.0
	github.com/golang/protobuf => github.com/golang/protobuf v1.5.2
	github.com/kpango/glg => github.com/kpango/glg v1.6.14
	github.com/pkg/sftp => github.com/pkg/sftp v1.13.5
<<<<<<< HEAD
	golang.org/x/crypto => golang.org/x/crypto v0.5.0
	golang.org/x/net => golang.org/x/net v0.5.0
	golang.org/x/text => golang.org/x/text v0.6.0
	google.golang.org/genproto => google.golang.org/genproto v0.0.0-20230202175211-008b39050e57
	google.golang.org/grpc => google.golang.org/grpc v1.52.3
=======
	golang.org/x/crypto => golang.org/x/crypto v0.0.0-20220829220503-c86fa9a7ed90
	golang.org/x/net => golang.org/x/net v0.0.0-20220906165146-f3363e06e74c
	golang.org/x/text => golang.org/x/text v0.3.7
	google.golang.org/genproto => google.golang.org/genproto v0.0.0-20220902135211-223410557253
	google.golang.org/grpc => google.golang.org/grpc v1.49.0
>>>>>>> feature/internal-net/support-quic
	google.golang.org/protobuf => google.golang.org/protobuf v1.28.1
	gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 => gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/kpango/fuid v0.0.0-20221203053508-503b5ad89aa1
	github.com/kpango/glg v1.6.14
	github.com/vdaas/vald-client-go v1.7.0
	gonum.org/v1/hdf5 v0.0.0-20210714002203-8c5d23bc6946
	google.golang.org/grpc v1.51.0
)

require (
	github.com/envoyproxy/protoc-gen-validate v0.1.0 // indirect
	github.com/goccy/go-json v0.9.11 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/kpango/fastime v1.1.6 // indirect
	golang.org/x/net v0.4.0 // indirect
	golang.org/x/sys v0.4.0 // indirect
	golang.org/x/text v0.6.0 // indirect
	google.golang.org/genproto v0.0.0-20221118155620-16455021b5e6 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)
