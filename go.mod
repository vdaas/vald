module github.com/vdaas/vald

go 1.15

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v14.2.0+incompatible
	github.com/aws/aws-sdk-go => github.com/aws/aws-sdk-go v1.34.7
	github.com/boltdb/bolt => github.com/boltdb/bolt v1.3.1
	github.com/cockroachdb/errors => github.com/cockroachdb/errors v1.6.1
	github.com/coreos/etcd => go.etcd.io/etcd v0.5.0-alpha.5.0.20200425165423-262c93980547
	github.com/docker/docker => github.com/moby/moby v17.12.0-ce-rc1.0.20200618181300-9dc6525e6118+incompatible
	github.com/envoyproxy/protoc-gen-validate => github.com/envoyproxy/protoc-gen-validate v0.4.1
	github.com/go-sql-driver/mysql => github.com/go-sql-driver/mysql v1.5.1-0.20200818111213-46351a889297
	github.com/gocql/gocql => github.com/gocql/gocql v0.0.0-20200815110948-5378c8f664e9
	github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.2-0.20200807193113-deb6fe8ca7c6
	github.com/gophercloud/gophercloud => github.com/gophercloud/gophercloud v0.12.0
	github.com/gorilla/mux => github.com/gorilla/mux v1.7.5-0.20200711200521-98cb6bf42e08
	github.com/gorilla/websocket => github.com/gorilla/websocket v1.4.2
	github.com/tensorflow/tensorflow => github.com/tensorflow/tensorflow v2.1.0+incompatible
	golang.org/x/crypto => golang.org/x/crypto v0.0.0-20200728195943-123391ffb6de
	k8s.io/api => k8s.io/api v0.18.8
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.18.8
	k8s.io/apimachinery => k8s.io/apimachinery v0.18.8
	k8s.io/client-go => k8s.io/client-go v0.18.8
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.6.2
)

require (
	cloud.google.com/go v0.64.1-0.20200818234218-e455ed5db86b
	code.cloudfoundry.org/bytefmt v0.0.0-20200131002437-cf55d5288a48
	contrib.go.opencensus.io/exporter/jaeger v0.2.1
	contrib.go.opencensus.io/exporter/prometheus v0.2.1-0.20200609204449-6bcf6f8577f0
	contrib.go.opencensus.io/exporter/stackdriver v0.13.3
	github.com/99designs/gqlgen v0.12.2 // indirect
	github.com/Masterminds/goutils v1.1.0 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Masterminds/sprig v2.22.0+incompatible // indirect
	github.com/agnivade/levenshtein v1.1.0 // indirect
	github.com/aokoli/goutils v1.1.0 // indirect
	github.com/asaskevich/govalidator v0.0.0-20200819183940-29e1ff8eb0bb // indirect
	github.com/aws/aws-sdk-go v1.33.17
	github.com/cespare/xxhash/v2 v2.1.1
	github.com/cockroachdb/errors v1.6.1
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/danielvladco/go-proto-gql/pb v0.6.1
	github.com/envoyproxy/protoc-gen-validate v0.4.2-0.20200814213351-f4107c7d75f2
	github.com/fsnotify/fsnotify v1.4.10-0.20200417215612-7f4cf4dd2b52
	github.com/go-openapi/spec v0.19.9 // indirect
	github.com/go-redis/redis/v7 v7.2.1-0.20200519055202-64bb0b7f3af4
	github.com/go-sql-driver/mysql v1.5.1-0.20200818111213-46351a889297
	github.com/go-swagger/go-swagger v0.25.0 // indirect
	github.com/gocql/gocql v0.0.0-20200131111108-92af2e088537
	github.com/gocraft/dbr/v2 v2.7.1-0.20200218045517-f487ccffc6d0
	github.com/gogo/protobuf v1.3.2-0.20200807193113-deb6fe8ca7c6
	github.com/google/go-cmp v0.5.2-0.20200818193711-d2fcc899bdc2
	github.com/google/gofuzz v1.1.0
	github.com/gorilla/handlers v1.5.0 // indirect
	github.com/gorilla/mux v1.7.5-0.20200711200521-98cb6bf42e08
	github.com/grpc-ecosystem/grpc-gateway v1.14.7 // indirect
	github.com/hashicorp/go-version v1.2.1
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/iancoleman/strcase v0.1.0 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/json-iterator/go v1.1.11-0.20200806011408-6821bec9fa5c
	github.com/klauspost/compress v1.10.12-0.20200818095508-f5ee0f4fc064
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/kpango/fastime v1.0.17-0.20200818143642-9b95d47eeba9
	github.com/kpango/fuid v0.0.0-20190507064958-80435564606b
	github.com/kpango/gache v1.2.2-0.20200709224359-34beea72198c
	github.com/kpango/glg v1.5.2-0.20200818134832-7fe3ee4e76c0
	github.com/kr/pretty v0.2.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.0.3
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/matryer/moq v0.0.0-20200816112511-720d53e65d2f // indirect
	github.com/mitchellh/copystructure v1.0.0 // indirect
	github.com/mitchellh/mapstructure v1.3.3 // indirect
	github.com/mitchellh/reflectwalk v1.0.1 // indirect
	github.com/mwitkow/go-proto-validators v0.3.2 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pierrec/lz4/v3 v3.3.2
	github.com/pseudomuto/protoc-gen-doc v1.3.2 // indirect
	github.com/scylladb/gocqlx v1.5.1-0.20200423154401-507391a34cf0
	github.com/spf13/viper v1.7.1 // indirect
	github.com/tensorflow/tensorflow v0.0.0-00010101000000-000000000000
	github.com/urfave/cli/v2 v2.2.0 // indirect
	github.com/vektah/dataloaden v0.3.0 // indirect
	github.com/yahoojapan/gongt v0.0.0-20190517050727-966dcc7aa5e8
	github.com/yahoojapan/ngtd v0.0.0-20200424071638-9872bbae3700
	go.mongodb.org/mongo-driver v1.4.0 // indirect
	go.opencensus.io v0.22.5-0.20200719225510-d7677d6af595
	go.uber.org/automaxprocs v1.3.1-0.20200415073007-b685be8c1c23
	go.uber.org/goleak v1.1.10
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a // indirect
	golang.org/x/net v0.0.0-20200822124328-c89045814202
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
	golang.org/x/sys v0.0.0-20200821140526-fda516888d29
	golang.org/x/tools v0.0.0-20200823205832-c024452afbcd // indirect
	gonum.org/v1/hdf5 v0.0.0-20200504100616-496fefe91614
	gonum.org/v1/plot v0.7.1-0.20200803120916-6a037fda5e90
	google.golang.org/api v0.30.1-0.20200818171802-feaa1c6611ae
	google.golang.org/genproto v0.0.0-20200815001618-f69a88009b70
	google.golang.org/grpc v1.32.0-dev.0.20200818224027-0f73133e3aa3
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/ini.v1 v1.60.1 // indirect
	gopkg.in/yaml.v2 v2.3.0
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
	k8s.io/api v0.18.8
	k8s.io/apimachinery v0.18.8
	k8s.io/client-go v0.18.8
	k8s.io/metrics v0.18.8
	sigs.k8s.io/controller-runtime v0.0.0-00010101000000-000000000000
)
