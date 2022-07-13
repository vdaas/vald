module github.com/vdaas/vald/example/client

go 1.18

replace (
	github.com/envoyproxy/protoc-gen-validate => github.com/envoyproxy/protoc-gen-validate v0.6.7
	github.com/goccy/go-json => github.com/goccy/go-json v0.9.8
	github.com/golang/protobuf => github.com/golang/protobuf v1.5.2
	github.com/kpango/glg => github.com/kpango/glg v1.6.11
	github.com/pkg/sftp => github.com/pkg/sftp v1.13.5
	golang.org/x/crypto => golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d
	golang.org/x/net => golang.org/x/net v0.0.0-20220630215102-69896b714898
	golang.org/x/text => golang.org/x/text v0.3.7
	google.golang.org/genproto => google.golang.org/genproto v0.0.0-20220630174209-ad1d48641aa7
	google.golang.org/grpc => google.golang.org/grpc v1.47.0
	google.golang.org/protobuf => google.golang.org/protobuf v1.28.0
	gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 => gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/kpango/fuid v0.0.0-20220209050620-e5987ba1ea5e
	github.com/kpango/glg v1.6.10
	github.com/vdaas/vald-client-go v1.5.5
	gonum.org/v1/hdf5 v0.0.0-20210714002203-8c5d23bc6946
	google.golang.org/grpc v1.47.0
)

require (
	github.com/envoyproxy/protoc-gen-validate v0.1.0 // indirect
	github.com/goccy/go-json v0.9.7 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/kpango/fastime v1.1.4 // indirect
	golang.org/x/net v0.0.0-20220624214902-1bab6f366d9e // indirect
	golang.org/x/sys v0.0.0-20220610221304-9f5ed59c137d // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220616135557-88e70c0c3a90 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)
