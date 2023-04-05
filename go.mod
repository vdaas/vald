module github.com/vdaas/vald

go 1.20

replace (
	cloud.google.com/go => cloud.google.com/go v0.110.0
	cloud.google.com/go/bigquery => cloud.google.com/go/bigquery v1.50.0
	cloud.google.com/go/compute => cloud.google.com/go/compute v1.19.0
	cloud.google.com/go/datastore => cloud.google.com/go/datastore v1.11.0
	cloud.google.com/go/firestore => cloud.google.com/go/firestore v1.9.0
	cloud.google.com/go/iam => cloud.google.com/go/iam v1.0.0
	cloud.google.com/go/kms => cloud.google.com/go/kms v1.10.1
	cloud.google.com/go/monitoring => cloud.google.com/go/monitoring v1.13.0
	cloud.google.com/go/pubsub => cloud.google.com/go/pubsub v1.30.0
	cloud.google.com/go/secretmanager => cloud.google.com/go/secretmanager v1.10.0
	cloud.google.com/go/storage => cloud.google.com/go/storage v1.30.1
	cloud.google.com/go/trace => cloud.google.com/go/trace v1.9.0
	code.cloudfoundry.org/bytefmt => code.cloudfoundry.org/bytefmt v0.0.0-20211005130812-5bb3c17173e5
	contrib.go.opencensus.io/exporter/aws => contrib.go.opencensus.io/exporter/aws v0.0.0-20200617204711-c478e41e60e9
	contrib.go.opencensus.io/exporter/prometheus => contrib.go.opencensus.io/exporter/prometheus v0.4.2
	contrib.go.opencensus.io/integrations/ocsql => contrib.go.opencensus.io/integrations/ocsql v0.1.7
	git.sr.ht/~sbinet/gg => git.sr.ht/~sbinet/gg v0.3.1
	github.com/AdaLogics/go-fuzz-headers => github.com/AdaLogics/go-fuzz-headers v0.0.0-20230106234847-43070de90fa1
	github.com/Azure/azure-amqp-common-go/v3 => github.com/Azure/azure-amqp-common-go/v3 v3.2.3
	github.com/Azure/azure-sdk-for-go => github.com/Azure/azure-sdk-for-go v68.0.0+incompatible
	github.com/Azure/azure-sdk-for-go/sdk/azcore => github.com/Azure/azure-sdk-for-go/sdk/azcore v1.4.0
	github.com/Azure/azure-sdk-for-go/sdk/azidentity => github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.2.2
	github.com/Azure/azure-sdk-for-go/sdk/internal => github.com/Azure/azure-sdk-for-go/sdk/internal v1.3.0
	github.com/Azure/azure-service-bus-go => github.com/Azure/azure-service-bus-go v0.11.5
	github.com/Azure/azure-storage-blob-go => github.com/Azure/azure-storage-blob-go v0.15.0
	github.com/Azure/go-amqp => github.com/Azure/go-amqp v0.19.1
	github.com/Azure/go-ansiterm => github.com/Azure/go-ansiterm v0.0.0-20230124172434-306776ec8161
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v14.2.1-0.20230315184244-553a90ae65a6+incompatible
	github.com/Azure/go-autorest/autorest => github.com/Azure/go-autorest/autorest v0.11.29-0.20230315184244-553a90ae65a6
	github.com/Azure/go-autorest/autorest/adal => github.com/Azure/go-autorest/autorest/adal v0.9.23
	github.com/Azure/go-autorest/autorest/azure/auth => github.com/Azure/go-autorest/autorest/azure/auth v0.5.13-0.20230315184244-553a90ae65a6
	github.com/Azure/go-autorest/autorest/azure/cli => github.com/Azure/go-autorest/autorest/azure/cli v0.4.7-0.20230315184244-553a90ae65a6
	github.com/Azure/go-autorest/autorest/date => github.com/Azure/go-autorest/autorest/date v0.3.1-0.20230315184244-553a90ae65a6
	github.com/Azure/go-autorest/autorest/mocks => github.com/Azure/go-autorest/autorest/mocks v0.4.3-0.20230315184244-553a90ae65a6
	github.com/Azure/go-autorest/autorest/to => github.com/Azure/go-autorest/autorest/to v0.4.1-0.20230315184244-553a90ae65a6
	github.com/Azure/go-autorest/autorest/validation => github.com/Azure/go-autorest/autorest/validation v0.3.2-0.20230315184244-553a90ae65a6
	github.com/Azure/go-autorest/logger => github.com/Azure/go-autorest/logger v0.2.2-0.20230315184244-553a90ae65a6
	github.com/Azure/go-autorest/tracing => github.com/Azure/go-autorest/tracing v0.6.1-0.20230315184244-553a90ae65a6
	github.com/BurntSushi/toml => github.com/BurntSushi/toml v1.2.1
	github.com/DATA-DOG/go-sqlmock => github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/DataDog/datadog-go => github.com/DataDog/datadog-go v4.8.3+incompatible
	github.com/GoogleCloudPlatform/cloudsql-proxy => github.com/GoogleCloudPlatform/cloudsql-proxy v1.33.5
	github.com/Masterminds/semver/v3 => github.com/Masterminds/semver/v3 v3.2.0
	github.com/Microsoft/go-winio => github.com/Microsoft/go-winio v0.6.0
	github.com/Microsoft/hcsshim => github.com/Microsoft/hcsshim v0.9.8
	github.com/NYTimes/gziphandler => github.com/NYTimes/gziphandler v1.1.1
	github.com/ajstarks/deck => github.com/ajstarks/deck v0.0.0-20230403145838-746d569493ac
	github.com/ajstarks/deck/generate => github.com/ajstarks/deck/generate v0.0.0-20230403145838-746d569493ac
	github.com/ajstarks/svgo => github.com/ajstarks/svgo v0.0.0-20211024235047-1546f124cd8b
	github.com/antihax/optional => github.com/antihax/optional v1.0.0
	github.com/armon/circbuf => github.com/armon/circbuf v0.0.0-20190214190532-5111143e8da2
	github.com/armon/go-metrics => github.com/armon/go-metrics v0.4.1
	github.com/armon/go-radix => github.com/armon/go-radix v1.0.0
	github.com/armon/go-socks5 => github.com/armon/go-socks5 v0.0.0-20160902184237-e75332964ef5
	github.com/asaskevich/govalidator => github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2
	github.com/aws/aws-sdk-go => github.com/aws/aws-sdk-go v1.44.236
	github.com/aws/aws-sdk-go-v2 => github.com/aws/aws-sdk-go-v2 v1.17.7
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream => github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.4.10
	github.com/aws/aws-sdk-go-v2/config => github.com/aws/aws-sdk-go-v2/config v1.18.19
	github.com/aws/aws-sdk-go-v2/credentials => github.com/aws/aws-sdk-go-v2/credentials v1.13.18
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds => github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.13.1
	github.com/aws/aws-sdk-go-v2/feature/s3/manager => github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.11.60
	github.com/aws/aws-sdk-go-v2/internal/configsources => github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.31
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.25
	github.com/aws/aws-sdk-go-v2/internal/ini => github.com/aws/aws-sdk-go-v2/internal/ini v1.3.32
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.11
	github.com/aws/aws-sdk-go-v2/service/internal/checksum => github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.1.26
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.25
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared => github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.14.0
	github.com/aws/aws-sdk-go-v2/service/kms => github.com/aws/aws-sdk-go-v2/service/kms v1.20.8
	github.com/aws/aws-sdk-go-v2/service/s3 => github.com/aws/aws-sdk-go-v2/service/s3 v1.31.1
	github.com/aws/aws-sdk-go-v2/service/secretsmanager => github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.19.1
	github.com/aws/aws-sdk-go-v2/service/sns => github.com/aws/aws-sdk-go-v2/service/sns v1.20.6
	github.com/aws/aws-sdk-go-v2/service/sqs => github.com/aws/aws-sdk-go-v2/service/sqs v1.20.6
	github.com/aws/aws-sdk-go-v2/service/ssm => github.com/aws/aws-sdk-go-v2/service/ssm v1.36.0
	github.com/aws/aws-sdk-go-v2/service/sso => github.com/aws/aws-sdk-go-v2/service/sso v1.12.6
	github.com/aws/aws-sdk-go-v2/service/sts => github.com/aws/aws-sdk-go-v2/service/sts v1.18.7
	github.com/aws/smithy-go => github.com/aws/smithy-go v1.13.5
	github.com/benbjohnson/clock => github.com/benbjohnson/clock v1.3.0
	github.com/beorn7/perks => github.com/beorn7/perks v1.0.1
	github.com/bgentry/speakeasy => github.com/bgentry/speakeasy v0.1.0
	github.com/bitly/go-hostpool => github.com/bitly/go-hostpool v0.1.0
	github.com/blang/semver => github.com/blang/semver v3.5.1+incompatible
	github.com/blang/semver/v4 => github.com/blang/semver/v4 v4.0.0
	github.com/bmizerany/assert => github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869
	github.com/boombuler/barcode => github.com/boombuler/barcode v1.0.1
	github.com/buger/jsonparser => github.com/buger/jsonparser v1.1.1
	github.com/cenkalti/backoff/v4 => github.com/cenkalti/backoff/v4 v4.2.0
	github.com/census-instrumentation/opencensus-proto => github.com/census-instrumentation/opencensus-proto v0.4.1
	github.com/certifi/gocertifi => github.com/certifi/gocertifi v0.0.0-20210507211836-431795d63e8d
	github.com/cespare/xxhash/v2 => github.com/cespare/xxhash/v2 v2.2.0
	github.com/checkpoint-restore/go-criu/v5 => github.com/checkpoint-restore/go-criu/v5 v5.3.0
	github.com/chzyer/logex => github.com/chzyer/logex v1.2.1
	github.com/chzyer/readline => github.com/chzyer/readline v1.5.1
	github.com/chzyer/test => github.com/chzyer/test v1.0.0
	github.com/cncf/udpa/go => github.com/cncf/udpa/go v0.0.0-20220112060539-c52dc94e7fbe
	github.com/cncf/xds/go => github.com/cncf/xds/go v0.0.0-20230310173818-32f1caf87195
	github.com/cockroachdb/apd => github.com/cockroachdb/apd v1.1.0
	github.com/cockroachdb/datadriven => github.com/cockroachdb/datadriven v1.0.2
	github.com/containerd/aufs => github.com/containerd/aufs v1.0.0
	github.com/containerd/btrfs => github.com/containerd/btrfs v1.0.0
	github.com/containerd/cgroups => github.com/containerd/cgroups v1.1.0
	github.com/containerd/console => github.com/containerd/console v1.0.3
	github.com/containerd/containerd => github.com/containerd/containerd v1.7.0
	github.com/containerd/continuity => github.com/containerd/continuity v0.3.0
	github.com/containerd/fifo => github.com/containerd/fifo v1.1.0
	github.com/containerd/go-cni => github.com/containerd/go-cni v1.1.9
	github.com/containerd/go-runc => github.com/containerd/go-runc v1.0.0
	github.com/containerd/imgcrypt => github.com/containerd/imgcrypt v1.1.7
	github.com/containerd/nri => github.com/containerd/nri v0.3.0
	github.com/containerd/stargz-snapshotter/estargz => github.com/containerd/stargz-snapshotter/estargz v0.14.3
	github.com/containerd/ttrpc => github.com/containerd/ttrpc v1.2.1
	github.com/containerd/typeurl => github.com/containerd/typeurl v1.0.2
	github.com/containerd/zfs => github.com/containerd/zfs v1.0.0
	github.com/containernetworking/cni => github.com/containernetworking/cni v1.1.2
	github.com/containernetworking/plugins => github.com/containernetworking/plugins v1.2.0
	github.com/containers/ocicrypt => github.com/containers/ocicrypt v1.1.7
	github.com/coreos/go-iptables => github.com/coreos/go-iptables v0.6.0
	github.com/coreos/go-oidc => github.com/coreos/go-oidc v2.2.1+incompatible
	github.com/coreos/go-semver => github.com/coreos/go-semver v0.3.1
	github.com/coreos/go-systemd/v22 => github.com/coreos/go-systemd/v22 v22.5.0
	github.com/cpuguy83/go-md2man/v2 => github.com/cpuguy83/go-md2man/v2 v2.0.2
	github.com/creack/pty => github.com/creack/pty v1.1.18
	github.com/cyphar/filepath-securejoin => github.com/cyphar/filepath-securejoin v0.2.3
	github.com/d2g/dhcp4 => github.com/d2g/dhcp4 v0.0.0-20170904100407-a1d1b6c41b1c
	github.com/d2g/dhcp4client => github.com/d2g/dhcp4client v1.0.0
	github.com/d2g/dhcp4server => github.com/d2g/dhcp4server v0.0.0-20181031114812-7d4a0a7f59a5
	github.com/davecgh/go-spew => github.com/davecgh/go-spew v1.1.1
	github.com/denisenkom/go-mssqldb => github.com/denisenkom/go-mssqldb v0.12.3
	github.com/dennwc/varint => github.com/dennwc/varint v1.0.0
	github.com/devigned/tab => github.com/devigned/tab v0.1.1
	github.com/dgryski/go-rendezvous => github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f
	github.com/dnaeon/go-vcr => github.com/dnaeon/go-vcr v1.2.0
	github.com/docker/cli => github.com/docker/cli v23.0.3+incompatible
	github.com/docker/distribution => github.com/docker/distribution v2.8.1+incompatible
	github.com/docker/docker => github.com/docker/docker v23.0.3+incompatible
	github.com/docker/docker-credential-helpers => github.com/docker/docker-credential-helpers v0.7.0
	github.com/docker/go-connections => github.com/docker/go-connections v0.4.0
	github.com/docker/go-events => github.com/docker/go-events v0.0.0-20190806004212-e31b211e4f1c
	github.com/docker/go-metrics => github.com/docker/go-metrics v0.0.1
	github.com/docker/go-units => github.com/docker/go-units v0.5.0
	github.com/docopt/docopt-go => github.com/docopt/docopt-go v0.0.0-20180111231733-ee0de3bc6815
	github.com/dustin/go-humanize => github.com/dustin/go-humanize v1.0.1
	github.com/edsrzf/mmap-go => github.com/edsrzf/mmap-go v1.1.0
	github.com/elazarl/goproxy => github.com/elazarl/goproxy v0.0.0-20221015165544-a0805db90819
	github.com/emicklei/go-restful => github.com/emicklei/go-restful v2.16.0+incompatible
	github.com/emicklei/go-restful/v3 => github.com/emicklei/go-restful/v3 v3.10.2
	github.com/envoyproxy/go-control-plane => github.com/envoyproxy/go-control-plane v0.11.0
	github.com/envoyproxy/protoc-gen-validate => github.com/envoyproxy/protoc-gen-validate v0.10.1
	github.com/evanphx/json-patch => github.com/evanphx/json-patch v0.5.2
	github.com/fogleman/gg => github.com/fogleman/gg v1.3.0
	github.com/form3tech-oss/jwt-go => github.com/form3tech-oss/jwt-go v3.2.5+incompatible
	github.com/fortytw2/leaktest => github.com/fortytw2/leaktest v1.3.0
	github.com/frankban/quicktest => github.com/frankban/quicktest v1.14.4
	github.com/fsnotify/fsnotify => github.com/fsnotify/fsnotify v1.6.0
	github.com/gin-contrib/sse => github.com/gin-contrib/sse v0.1.0
	github.com/gin-gonic/gin => github.com/gin-gonic/gin v1.9.0
	github.com/go-errors/errors => github.com/go-errors/errors v1.4.2
	github.com/go-fonts/dejavu => github.com/go-fonts/dejavu v0.1.0
	github.com/go-fonts/latin-modern => github.com/go-fonts/latin-modern v0.3.0
	github.com/go-fonts/liberation => github.com/go-fonts/liberation v0.3.0
	github.com/go-fonts/stix => github.com/go-fonts/stix v0.1.0
	github.com/go-gl/gl => github.com/go-gl/gl v0.0.0-20211210172815-726fda9656d6
	github.com/go-gl/glfw/v3.3/glfw => github.com/go-gl/glfw/v3.3/glfw v0.0.0-20221017161538-93cebf72946b
	github.com/go-kit/kit => github.com/go-kit/kit v0.12.1-0.20230302021612-7f14cb4dc16c
	github.com/go-kit/log => github.com/go-kit/log v0.2.1
	github.com/go-latex/latex => github.com/go-latex/latex v0.0.0-20230307184459-12ec69307ad9
	github.com/go-logfmt/logfmt => github.com/go-logfmt/logfmt v0.6.0
	github.com/go-logr/logr => github.com/go-logr/logr v1.2.4
	github.com/go-logr/stdr => github.com/go-logr/stdr v1.2.2
	github.com/go-logr/zapr => github.com/go-logr/zapr v1.2.3
	github.com/go-openapi/jsonpointer => github.com/go-openapi/jsonpointer v0.19.6
	github.com/go-openapi/jsonreference => github.com/go-openapi/jsonreference v0.20.2
	github.com/go-openapi/loads => github.com/go-openapi/loads v0.21.2
	github.com/go-openapi/runtime => github.com/go-openapi/runtime v0.25.0
	github.com/go-openapi/spec => github.com/go-openapi/spec v0.20.8
	github.com/go-openapi/strfmt => github.com/go-openapi/strfmt v0.21.7
	github.com/go-openapi/swag => github.com/go-openapi/swag v0.22.3
	github.com/go-openapi/validate => github.com/go-openapi/validate v0.22.1
	github.com/go-pdf/fpdf => github.com/go-pdf/fpdf v1.4.3
	github.com/go-playground/assert/v2 => github.com/go-playground/assert/v2 v2.2.0
	github.com/go-playground/locales => github.com/go-playground/locales v0.14.1
	github.com/go-playground/universal-translator => github.com/go-playground/universal-translator v0.18.1
	github.com/go-playground/validator/v10 => github.com/go-playground/validator/v10 v10.12.0
	github.com/go-redis/redis/v8 => github.com/go-redis/redis/v8 v8.11.5
	github.com/go-resty/resty/v2 => github.com/go-resty/resty/v2 v2.7.0
	github.com/go-sql-driver/mysql => github.com/go-sql-driver/mysql v1.7.0
	github.com/go-stack/stack => github.com/go-stack/stack v1.8.1
	github.com/go-task/slim-sprig => github.com/go-task/slim-sprig v2.20.0+incompatible
	github.com/go-toolsmith/astcopy => github.com/go-toolsmith/astcopy v1.1.0
	github.com/go-toolsmith/astequal => github.com/go-toolsmith/astequal v1.1.0
	github.com/go-toolsmith/strparse => github.com/go-toolsmith/strparse v1.1.0
	github.com/go-zookeeper/zk => github.com/go-zookeeper/zk v1.0.3
	github.com/gobwas/httphead => github.com/gobwas/httphead v0.1.0
	github.com/gobwas/pool => github.com/gobwas/pool v0.2.1
	github.com/gobwas/ws => github.com/gobwas/ws v1.1.0
	github.com/goccy/go-json => github.com/goccy/go-json v0.10.2
	github.com/goccy/go-yaml => github.com/goccy/go-yaml v1.11.0
	github.com/gocql/gocql => github.com/gocql/gocql v1.3.2
	github.com/gocraft/dbr/v2 => github.com/gocraft/dbr/v2 v2.7.3
	github.com/godbus/dbus/v5 => github.com/godbus/dbus/v5 v5.1.0
	github.com/gofrs/uuid => github.com/gofrs/uuid v4.4.0+incompatible
	github.com/gogo/googleapis => github.com/gogo/googleapis v1.4.1
	github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.2
	github.com/golang-jwt/jwt/v4 => github.com/golang-jwt/jwt/v4 v4.5.0
	github.com/golang-sql/civil => github.com/golang-sql/civil v0.0.0-20220223132316-b832511892a9
	github.com/golang-sql/sqlexp => github.com/golang-sql/sqlexp v0.1.0
	github.com/golang/freetype => github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/golang/glog => github.com/golang/glog v1.1.1
	github.com/golang/groupcache => github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da
	github.com/golang/mock => github.com/golang/mock v1.6.0
	github.com/golang/protobuf => github.com/golang/protobuf v1.5.3
	github.com/golang/snappy => github.com/golang/snappy v0.0.4
	github.com/google/btree => github.com/google/btree v1.1.2
	github.com/google/gnostic => github.com/google/gnostic v0.6.9
	github.com/google/go-cmp => github.com/google/go-cmp v0.5.9
	github.com/google/go-containerregistry => github.com/google/go-containerregistry v0.14.0
	github.com/google/go-querystring => github.com/google/go-querystring v1.1.0
	github.com/google/go-replayers/grpcreplay => github.com/google/go-replayers/grpcreplay v1.1.0
	github.com/google/go-replayers/httpreplay => github.com/google/go-replayers/httpreplay v1.2.0
	github.com/google/gofuzz => github.com/google/gofuzz v1.2.0
	github.com/google/martian => github.com/google/martian v2.1.0+incompatible
	github.com/google/martian/v3 => github.com/google/martian/v3 v3.3.2
	github.com/google/pprof => github.com/google/pprof v0.0.0-20230323073829-e72429f035bd
	github.com/google/shlex => github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510
	github.com/google/subcommands => github.com/google/subcommands v1.2.0
	github.com/google/uuid => github.com/google/uuid v1.3.0
	github.com/google/wire => github.com/google/wire v0.5.0
	github.com/googleapis/gax-go/v2 => github.com/googleapis/gax-go/v2 v2.8.0
	github.com/gorilla/mux => github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket => github.com/gorilla/websocket v1.5.0
	github.com/grafana/regexp => github.com/grafana/regexp v0.0.0-20221122212121-6b5c0a4cb7fd
	github.com/gregjones/httpcache => github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79
	github.com/grpc-ecosystem/go-grpc-middleware => github.com/grpc-ecosystem/go-grpc-middleware v1.4.0
	github.com/grpc-ecosystem/go-grpc-prometheus => github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/grpc-ecosystem/grpc-gateway => github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/grpc-ecosystem/grpc-gateway/v2 => github.com/grpc-ecosystem/grpc-gateway/v2 v2.15.2
	github.com/hailocab/go-hostpool => github.com/kpango/go-hostpool v0.0.0-20210303030322-aab80263dcd0
	github.com/hanwen/go-fuse => github.com/hanwen/go-fuse v1.0.0
	github.com/hanwen/go-fuse/v2 => github.com/hanwen/go-fuse/v2 v2.2.0
	github.com/hashicorp/consul/api => github.com/hashicorp/consul/api v1.20.0
	github.com/hashicorp/consul/sdk => github.com/hashicorp/consul/sdk v0.13.1
	github.com/hashicorp/errwrap => github.com/hashicorp/errwrap v1.1.0
	github.com/hashicorp/go-cleanhttp => github.com/hashicorp/go-cleanhttp v0.5.2
	github.com/hashicorp/go-hclog => github.com/hashicorp/go-hclog v1.5.0
	github.com/hashicorp/go-immutable-radix => github.com/hashicorp/go-immutable-radix v1.3.1
	github.com/hashicorp/go-msgpack => github.com/hashicorp/go-msgpack v1.1.6
	github.com/hashicorp/go-multierror => github.com/hashicorp/go-multierror v1.1.1
	github.com/hashicorp/go-retryablehttp => github.com/hashicorp/go-retryablehttp v0.7.2
	github.com/hashicorp/go-rootcerts => github.com/hashicorp/go-rootcerts v1.0.2
	github.com/hashicorp/go-sockaddr => github.com/hashicorp/go-sockaddr v1.0.2
	github.com/hashicorp/go-syslog => github.com/hashicorp/go-syslog v1.0.0
	github.com/hashicorp/go-uuid => github.com/hashicorp/go-uuid v1.0.3
	github.com/hashicorp/go-version => github.com/hashicorp/go-version v1.6.0
	github.com/hashicorp/golang-lru => github.com/hashicorp/golang-lru v1.0.1
	github.com/hashicorp/logutils => github.com/hashicorp/logutils v1.0.0
	github.com/hashicorp/mdns => github.com/hashicorp/mdns v1.0.5
	github.com/hashicorp/memberlist => github.com/hashicorp/memberlist v0.5.0
	github.com/hashicorp/serf => github.com/hashicorp/serf v0.10.1
	github.com/hetznercloud/hcloud-go => github.com/hetznercloud/hcloud-go v1.41.0
	github.com/iancoleman/strcase => github.com/iancoleman/strcase v0.2.0
	github.com/ianlancetaylor/demangle => github.com/ianlancetaylor/demangle v0.0.0-20230322204757-857afb9054cd
	github.com/imdario/mergo => github.com/imdario/mergo v0.3.15
	github.com/inconshreveable/mousetrap => github.com/inconshreveable/mousetrap v1.1.0
	github.com/intel/goresctrl => github.com/intel/goresctrl v0.3.0
	github.com/jackc/chunkreader/v2 => github.com/jackc/chunkreader/v2 v2.0.1
	github.com/jackc/pgconn => github.com/jackc/pgconn v1.14.0
	github.com/jackc/pgio => github.com/jackc/pgio v1.0.0
	github.com/jackc/pgmock => github.com/jackc/pgmock v0.0.0-20210724152146-4ad1a8207f65
	github.com/jackc/pgpassfile => github.com/jackc/pgpassfile v1.0.0
	github.com/jackc/pgproto3/v2 => github.com/jackc/pgproto3/v2 v2.3.2
	github.com/jackc/pgservicefile => github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a
	github.com/jackc/pgtype => github.com/jackc/pgtype v1.14.0
	github.com/jackc/pgx/v4 => github.com/jackc/pgx/v4 v4.18.1
	github.com/jackc/puddle => github.com/jackc/puddle v1.3.0
	github.com/jessevdk/go-flags => github.com/jessevdk/go-flags v1.5.0
	github.com/jmespath/go-jmespath => github.com/jmespath/go-jmespath v0.4.0
	github.com/jmespath/go-jmespath/internal/testify => github.com/jmespath/go-jmespath/internal/testify v1.5.1
	github.com/jmoiron/sqlx => github.com/jmoiron/sqlx v1.3.5
	github.com/joho/godotenv => github.com/joho/godotenv v1.5.1
	github.com/josharian/intern => github.com/josharian/intern v1.0.0
	github.com/json-iterator/go => github.com/json-iterator/go v1.1.12
	github.com/jstemmer/go-junit-report => github.com/jstemmer/go-junit-report v1.0.0
	github.com/jtolds/gls => github.com/jtolds/gls v4.20.0+incompatible
	github.com/julienschmidt/httprouter => github.com/julienschmidt/httprouter v1.3.0
	github.com/kisielk/errcheck => github.com/kisielk/errcheck v1.6.3
	github.com/kisielk/gotool => github.com/kisielk/gotool v1.0.0
	github.com/klauspost/compress => github.com/klauspost/compress v1.16.4-0.20230403072145-9243a1faf01d
	github.com/klauspost/cpuid/v2 => github.com/klauspost/cpuid/v2 v2.2.4
	github.com/kolo/xmlrpc => github.com/kolo/xmlrpc v0.0.0-20220921171641-a4b6fa1dd06b
	github.com/kpango/fastime => github.com/kpango/fastime v1.1.9
	github.com/kpango/fuid => github.com/kpango/fuid v0.0.0-20221203053508-503b5ad89aa1
	github.com/kpango/gache => github.com/kpango/gache v1.2.8
	github.com/kpango/glg => github.com/kpango/glg v1.6.15
	github.com/kr/fs => github.com/kr/fs v0.1.0
	github.com/kr/pretty => github.com/kr/pretty v0.3.1
	github.com/kr/text => github.com/kr/text v0.2.0
	github.com/kylelemons/godebug => github.com/kylelemons/godebug v1.1.0
	github.com/leanovate/gopter => github.com/leanovate/gopter v0.2.9
	github.com/leodido/go-urn => github.com/leodido/go-urn v1.2.2
	github.com/lib/pq => github.com/lib/pq v1.10.7
	github.com/liggitt/tabwriter => github.com/liggitt/tabwriter v0.0.0-20181228230101-89fcab3d43de
	github.com/linode/linodego => github.com/linode/linodego v1.16.0
	github.com/linuxkit/virtsock => github.com/linuxkit/virtsock v0.0.0-20220523201153-1a23e78aa7a2
	github.com/lucasb-eyer/go-colorful => github.com/lucasb-eyer/go-colorful v1.2.0
	github.com/lyft/protoc-gen-star => github.com/lyft/protoc-gen-star v0.6.2
	github.com/mailru/easyjson => github.com/mailru/easyjson v0.7.7
	github.com/mattn/go-colorable => github.com/mattn/go-colorable v0.1.13
	github.com/mattn/go-isatty => github.com/mattn/go-isatty v0.0.18
	github.com/mattn/go-shellwords => github.com/mattn/go-shellwords v1.0.12
	github.com/mattn/go-sqlite3 => github.com/mattn/go-sqlite3 v1.14.16
	github.com/matttproud/golang_protobuf_extensions => github.com/matttproud/golang_protobuf_extensions v1.0.4
	github.com/mitchellh/colorstring => github.com/mitchellh/colorstring v0.0.0-20190213212951-d06e56a500db
	github.com/mitchellh/go-homedir => github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/go-wordwrap => github.com/mitchellh/go-wordwrap v1.0.1
	github.com/mitchellh/mapstructure => github.com/mitchellh/mapstructure v1.5.0
	github.com/moby/locker => github.com/moby/locker v1.0.1
	github.com/moby/spdystream => github.com/moby/spdystream v0.2.0
	github.com/moby/sys/mountinfo => github.com/moby/sys/mountinfo v0.6.2
	github.com/moby/sys/signal => github.com/moby/sys/signal v0.7.0
	github.com/moby/sys/symlink => github.com/moby/sys/symlink v0.2.0
	github.com/moby/term => github.com/moby/term v0.0.0-20221205130635-1aeaba878587
	github.com/modern-go/concurrent => github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd
	github.com/modern-go/reflect2 => github.com/modern-go/reflect2 v1.0.2
	github.com/modocache/gover => github.com/modocache/gover v0.0.0-20171022184752-b58185e213c5
	github.com/monochromegane/go-gitignore => github.com/monochromegane/go-gitignore v0.0.0-20200626010858-205db1a8cc00
	github.com/montanaflynn/stats => github.com/montanaflynn/stats v0.7.0
	github.com/morikuni/aec => github.com/morikuni/aec v1.0.0
	github.com/mrunalp/fileutils => github.com/mrunalp/fileutils v0.5.0
	github.com/munnerz/goautoneg => github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822
	github.com/niemeyer/pretty => github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e
	github.com/nxadm/tail => github.com/nxadm/tail v1.4.8
	github.com/oklog/run => github.com/oklog/run v1.1.0
	github.com/oklog/ulid => github.com/oklog/ulid v1.3.1
	github.com/onsi/ginkgo => github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega => github.com/onsi/gomega v1.27.6
	github.com/peterbourgon/diskv => github.com/peterbourgon/diskv v2.0.1+incompatible
	github.com/phpdave11/gofpdf => github.com/phpdave11/gofpdf v1.4.2
	github.com/phpdave11/gofpdi => github.com/phpdave11/gofpdi v1.0.13
	github.com/pierrec/cmdflag => github.com/pierrec/cmdflag v0.0.2
	github.com/pierrec/lz4/v3 => github.com/pierrec/lz4/v3 v3.3.5
	github.com/pkg/browser => github.com/pkg/browser v0.0.0-20210911075715-681adbf594b8
	github.com/pkg/errors => github.com/pkg/errors v0.9.1
	github.com/pkg/sftp => github.com/pkg/sftp v1.13.5
	github.com/pmezard/go-difflib => github.com/pmezard/go-difflib v1.0.0
	github.com/posener/complete => github.com/posener/complete v1.2.3
	github.com/pquerna/cachecontrol => github.com/pquerna/cachecontrol v0.1.0
	github.com/prashantv/gostub => github.com/prashantv/gostub v1.1.0
	github.com/prometheus/alertmanager => github.com/prometheus/alertmanager v0.25.0
	github.com/prometheus/client_golang => github.com/prometheus/client_golang v1.14.0
	github.com/prometheus/client_model => github.com/prometheus/client_model v0.3.0
	github.com/prometheus/common => github.com/prometheus/common v0.42.0
	github.com/prometheus/common/assets => github.com/prometheus/common/assets v0.2.0
	github.com/prometheus/common/sigv4 => github.com/prometheus/common/sigv4 v0.1.0
	github.com/prometheus/exporter-toolkit => github.com/prometheus/exporter-toolkit v0.9.1
	github.com/prometheus/procfs => github.com/prometheus/procfs v0.9.0
	github.com/prometheus/prometheus => github.com/prometheus/prometheus v1.99.0
	github.com/prometheus/prometheus/v2 => github.com/prometheus/prometheus/v2 v2.35.0-retract
	github.com/quasilyte/go-ruleguard => github.com/quasilyte/go-ruleguard v0.3.19
	github.com/quasilyte/go-ruleguard/dsl => github.com/quasilyte/go-ruleguard/dsl v0.3.22
	github.com/quasilyte/gogrep => github.com/quasilyte/gogrep v0.5.0
	github.com/quasilyte/stdinfo => github.com/quasilyte/stdinfo v0.0.0-20220114132959-f7386bf02567
	github.com/rogpeppe/fastuuid => github.com/rogpeppe/fastuuid v1.2.0
	github.com/rogpeppe/go-internal => github.com/rogpeppe/go-internal v1.10.0
	github.com/rs/cors => github.com/rs/cors v1.8.3
	github.com/rs/xid => github.com/rs/xid v1.4.0
	github.com/rs/zerolog => github.com/rs/zerolog v1.29.0
	github.com/russross/blackfriday/v2 => github.com/russross/blackfriday/v2 v2.1.0
	github.com/ruudk/golang-pdf417 => github.com/ruudk/golang-pdf417 v0.0.0-20201230142125-a7e3863a1245
	github.com/ryanuber/columnize => github.com/ryanuber/columnize v2.1.2+incompatible
	github.com/safchain/ethtool => github.com/safchain/ethtool v0.3.0
	github.com/satori/go.uuid => github.com/satori/go.uuid v1.2.0
	github.com/scaleway/scaleway-sdk-go => github.com/scaleway/scaleway-sdk-go v1.0.0-beta.15
	github.com/schollz/progressbar/v2 => github.com/schollz/progressbar/v2 v2.15.0
	github.com/scylladb/go-reflectx => github.com/scylladb/go-reflectx v1.0.1
	github.com/scylladb/gocqlx => github.com/scylladb/gocqlx v1.5.0
	github.com/sean-/seed => github.com/sean-/seed v0.0.0-20170313163322-e2103e2c3529
	github.com/seccomp/libseccomp-golang => github.com/seccomp/libseccomp-golang v0.10.0
	github.com/sergi/go-diff => github.com/sergi/go-diff v1.3.1
	github.com/shopspring/decimal => github.com/shopspring/decimal v1.3.1
	github.com/shurcooL/httpfs => github.com/shurcooL/httpfs v0.0.0-20190707220628-8d4bc4ba7749
	github.com/shurcooL/sanitized_anchor_name => github.com/shurcooL/sanitized_anchor_name v1.0.0
	github.com/shurcooL/vfsgen => github.com/shurcooL/vfsgen v0.0.0-20200824052919-0d455de96546
	github.com/sirupsen/logrus => github.com/sirupsen/logrus v1.9.0
	github.com/smartystreets/assertions => github.com/smartystreets/assertions v1.13.1
	github.com/smartystreets/goconvey => github.com/smartystreets/goconvey v1.7.2
	github.com/soheilhy/cmux => github.com/soheilhy/cmux v0.1.5
	github.com/spf13/afero => github.com/spf13/afero v1.9.5
	github.com/spf13/cast => github.com/spf13/cast v1.5.0
	github.com/spf13/cobra => github.com/spf13/cobra v1.7.0
	github.com/spf13/pflag => github.com/spf13/pflag v1.0.5
	github.com/stoewer/go-strcase => github.com/stoewer/go-strcase v1.3.0
	github.com/stretchr/objx => github.com/stretchr/objx v0.5.0
	github.com/stretchr/testify => github.com/stretchr/testify v1.8.2
	github.com/ugorji/go/codec => github.com/ugorji/go/codec v1.2.11
	github.com/vdaas/vald-client-go => github.com/vdaas/vald-client-go v1.7.4
	github.com/vishvananda/netlink => github.com/vishvananda/netlink v1.1.0
	github.com/vishvananda/netns => github.com/vishvananda/netns v0.0.4
	github.com/xdg-go/pbkdf2 => github.com/xdg-go/pbkdf2 v1.0.0
	github.com/xdg-go/scram => github.com/xdg-go/scram v1.1.2
	github.com/xdg-go/stringprep => github.com/xdg-go/stringprep v1.0.4
	github.com/xeipuuv/gojsonpointer => github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb
	github.com/xeipuuv/gojsonreference => github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415
	github.com/xeipuuv/gojsonschema => github.com/xeipuuv/gojsonschema v1.2.0
	github.com/xiang90/probing => github.com/xiang90/probing v0.0.0-20221125231312-a49e3df8f510
	github.com/xlab/treeprint => github.com/xlab/treeprint v1.2.0
	github.com/zeebo/assert => github.com/zeebo/assert v1.3.1
	github.com/zeebo/xxh3 => github.com/zeebo/xxh3 v1.0.2
	go.etcd.io/bbolt => go.etcd.io/bbolt v1.3.7
	go.etcd.io/etcd/api/v3 => go.etcd.io/etcd/api/v3 v3.5.7
	go.etcd.io/etcd/client/pkg/v3 => go.etcd.io/etcd/client/pkg/v3 v3.5.7
	go.etcd.io/etcd/client/v2 => go.etcd.io/etcd/client/v2 v2.305.7
	go.etcd.io/etcd/client/v3 => go.etcd.io/etcd/client/v3 v3.5.7
	go.etcd.io/etcd/pkg/v3 => go.etcd.io/etcd/pkg/v3 v3.5.7
	go.etcd.io/etcd/raft/v3 => go.etcd.io/etcd/raft/v3 v3.5.7
	go.etcd.io/etcd/server/v3 => go.etcd.io/etcd/server/v3 v3.5.7
	go.mongodb.org/mongo-driver => go.mongodb.org/mongo-driver v1.11.4
	go.mozilla.org/pkcs7 => go.mozilla.org/pkcs7 v0.0.0-20210826202110-33d05740a352
	go.opencensus.io => go.opencensus.io v0.24.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc => go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.37.0
	go.opentelemetry.io/otel => go.opentelemetry.io/otel v1.11.1
	go.opentelemetry.io/otel/exporters/otlp/internal/retry => go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.11.1
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric => go.opentelemetry.io/otel/exporters/otlp/otlpmetric v0.33.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc => go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v0.33.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace => go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.11.1
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc => go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.11.1
	go.opentelemetry.io/otel/metric => go.opentelemetry.io/otel/metric v0.33.0
	go.opentelemetry.io/otel/sdk => go.opentelemetry.io/otel/sdk v1.11.1
	go.opentelemetry.io/otel/sdk/metric => go.opentelemetry.io/otel/sdk/metric v0.33.0
	go.opentelemetry.io/otel/trace => go.opentelemetry.io/otel/trace v1.11.1
	go.opentelemetry.io/proto/otlp => go.opentelemetry.io/proto/otlp v0.19.0
	go.starlark.net => go.starlark.net v0.0.0-20230302034142-4b1e35fe2254
	go.uber.org/atomic => go.uber.org/atomic v1.10.0
	go.uber.org/automaxprocs => go.uber.org/automaxprocs v1.5.2
	go.uber.org/goleak => go.uber.org/goleak v1.2.1
	go.uber.org/multierr => go.uber.org/multierr v1.11.0
	go.uber.org/zap => go.uber.org/zap v1.24.0
	gocloud.dev => gocloud.dev v0.29.0
	golang.org/x/crypto => golang.org/x/crypto v0.7.0
	golang.org/x/exp => golang.org/x/exp v0.0.0-20230321023759-10a507213a29
	golang.org/x/exp/typeparams => golang.org/x/exp/typeparams v0.0.0-20230321023759-10a507213a29
	golang.org/x/image => golang.org/x/image v0.6.0
	golang.org/x/lint => golang.org/x/lint v0.0.0-20210508222113-6edffad5e616
	golang.org/x/mobile => golang.org/x/mobile v0.0.0-20230301163155-e0f57694e12c
	golang.org/x/mod => golang.org/x/mod v0.10.0
	golang.org/x/net => golang.org/x/net v0.8.0
	golang.org/x/oauth2 => golang.org/x/oauth2 v0.6.0
	golang.org/x/sync => golang.org/x/sync v0.1.0
	golang.org/x/sys => golang.org/x/sys v0.7.0
	golang.org/x/term => golang.org/x/term v0.7.0
	golang.org/x/text => golang.org/x/text v0.8.0
	golang.org/x/time => golang.org/x/time v0.3.0
	golang.org/x/tools => golang.org/x/tools v0.7.0
	golang.org/x/xerrors => golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2
	gomodules.xyz/jsonpatch/v2 => gomodules.xyz/jsonpatch/v2 v2.2.0
	gonum.org/v1/gonum => gonum.org/v1/gonum v0.12.0
	gonum.org/v1/hdf5 => gonum.org/v1/hdf5 v0.0.0-20210714002203-8c5d23bc6946
	gonum.org/v1/plot => gonum.org/v1/plot v0.12.0
	google.golang.org/api => google.golang.org/api v0.115.0
	google.golang.org/appengine => google.golang.org/appengine v1.6.7
	google.golang.org/genproto => google.golang.org/genproto v0.0.0-20230403163135-c38d8f061ccd
	google.golang.org/grpc => google.golang.org/grpc v1.54.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc => google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.3.0
	google.golang.org/protobuf => google.golang.org/protobuf v1.30.0
	gopkg.in/alecthomas/kingpin.v2 => gopkg.in/alecthomas/kingpin.v2 v2.3.2
	gopkg.in/check.v1 => gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c
	gopkg.in/gcfg.v1 => gopkg.in/gcfg.v1 v1.2.3
	gopkg.in/inconshreveable/log15.v2 => gopkg.in/inconshreveable/log15.v2 v2.16.0
	gopkg.in/inf.v0 => gopkg.in/inf.v0 v0.9.1
	gopkg.in/tomb.v1 => gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7
	gopkg.in/warnings.v0 => gopkg.in/warnings.v0 v0.1.2
	gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 => gopkg.in/yaml.v3 v3.0.1
	gotest.tools/v3 => gotest.tools/v3 v3.4.0
	honnef.co/go/tools => honnef.co/go/tools v0.4.3
	k8s.io/api => k8s.io/api v0.26.3
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.26.3
	k8s.io/apimachinery => k8s.io/apimachinery v0.26.3
	k8s.io/apiserver => k8s.io/apiserver v0.26.3
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.26.3
	k8s.io/client-go => k8s.io/client-go v0.26.3
	k8s.io/component-base => k8s.io/component-base v0.26.3
	k8s.io/cri-api => k8s.io/cri-api v0.26.3
	k8s.io/gengo => k8s.io/gengo v0.0.0-20230306165830-ab3349d207d4
	k8s.io/klog => k8s.io/klog v1.0.0
	k8s.io/klog/v2 => k8s.io/klog/v2 v2.90.1
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20230327201221-f5883ff37f0c
	k8s.io/kubernetes => k8s.io/kubernetes v0.26.3
	k8s.io/metrics => k8s.io/metrics v0.26.3
	nhooyr.io/websocket => nhooyr.io/websocket v1.8.7
	rsc.io/pdf => rsc.io/pdf v0.1.1
	sigs.k8s.io/apiserver-network-proxy/konnectivity-client => sigs.k8s.io/apiserver-network-proxy/konnectivity-client v0.1.2
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.14.6
	sigs.k8s.io/json => sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd
	sigs.k8s.io/kustomize => sigs.k8s.io/kustomize v2.0.3+incompatible
	sigs.k8s.io/structured-merge-diff/v4 => sigs.k8s.io/structured-merge-diff/v4 v4.2.3
	sigs.k8s.io/yaml => sigs.k8s.io/yaml v1.3.0
)

require (
	cloud.google.com/go/storage v1.29.0
	code.cloudfoundry.org/bytefmt v0.0.0-20190710193110-1eb035ffe2b6
	github.com/aws/aws-sdk-go v1.44.200
	github.com/envoyproxy/protoc-gen-validate v0.9.1
	github.com/fsnotify/fsnotify v1.6.0
	github.com/go-redis/redis/v8 v8.0.0-00010101000000-000000000000
	github.com/go-sql-driver/mysql v1.7.0
	github.com/goccy/go-json v0.10.2
	github.com/gocql/gocql v0.0.0-20200131111108-92af2e088537
	github.com/gocraft/dbr/v2 v2.0.0-00010101000000-000000000000
	github.com/google/go-cmp v0.5.9
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v0.0.0-00010101000000-000000000000
	github.com/hashicorp/go-version v0.0.0-00010101000000-000000000000
	github.com/klauspost/compress v1.15.9
	github.com/kpango/fastime v1.1.9
	github.com/kpango/fuid v0.0.0-00010101000000-000000000000
	github.com/kpango/gache v0.0.0-00010101000000-000000000000
	github.com/kpango/glg v1.6.14
	github.com/leanovate/gopter v0.0.0-00010101000000-000000000000
	github.com/lucasb-eyer/go-colorful v0.0.0-00010101000000-000000000000
	github.com/pierrec/lz4/v3 v3.0.0-00010101000000-000000000000
	github.com/quasilyte/go-ruleguard v0.0.0-00010101000000-000000000000
	github.com/quasilyte/go-ruleguard/dsl v0.3.22
	github.com/scylladb/gocqlx v0.0.0-00010101000000-000000000000
	github.com/vdaas/vald-client-go v0.0.0-00010101000000-000000000000
	github.com/zeebo/xxh3 v1.0.2
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.35.0
	go.opentelemetry.io/otel v1.11.2
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.11.1
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.10.0
	go.opentelemetry.io/otel/metric v0.34.0
	go.opentelemetry.io/otel/sdk v1.11.1
	go.opentelemetry.io/otel/sdk/metric v0.33.0
	go.opentelemetry.io/otel/trace v1.11.2
	go.uber.org/automaxprocs v0.0.0-00010101000000-000000000000
	go.uber.org/goleak v1.2.0
	go.uber.org/zap v1.24.0
	gocloud.dev v0.0.0-00010101000000-000000000000
	golang.org/x/exp v0.0.0-20220827204233-334a2380cb91
	golang.org/x/net v0.8.0
	golang.org/x/oauth2 v0.6.0
	golang.org/x/sync v0.1.0
	golang.org/x/sys v0.7.0
	golang.org/x/text v0.8.0
	golang.org/x/tools v0.7.0
	gonum.org/v1/hdf5 v0.0.0-00010101000000-000000000000
	gonum.org/v1/plot v0.10.1
	google.golang.org/genproto v0.0.0-20230331144136-dcfb400f0633
	google.golang.org/grpc v1.54.0
	google.golang.org/protobuf v1.30.0
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.26.3
	k8s.io/apimachinery v0.26.3
	k8s.io/cli-runtime v0.0.0-00010101000000-000000000000
	k8s.io/client-go v0.26.3
	k8s.io/metrics v0.0.0-00010101000000-000000000000
	sigs.k8s.io/controller-runtime v0.0.0-00010101000000-000000000000
)

require (
	cloud.google.com/go v0.110.0 // indirect
	cloud.google.com/go/compute v1.19.0 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	cloud.google.com/go/iam v0.13.0 // indirect
	git.sr.ht/~sbinet/gg v0.3.1 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/ajstarks/svgo v0.0.0-20211024235047-1546f124cd8b // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cenkalti/backoff/v4 v4.1.3 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/emicklei/go-restful/v3 v3.9.0 // indirect
	github.com/evanphx/json-patch v4.12.0+incompatible // indirect
	github.com/evanphx/json-patch/v5 v5.6.0 // indirect
	github.com/go-errors/errors v1.0.1 // indirect
	github.com/go-fonts/liberation v0.3.0 // indirect
	github.com/go-latex/latex v0.0.0-20210823091927-c0d11ff05a81 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-openapi/jsonpointer v0.19.6 // indirect
	github.com/go-openapi/jsonreference v0.20.1 // indirect
	github.com/go-openapi/swag v0.22.3 // indirect
	github.com/go-pdf/fpdf v0.7.0 // indirect
	github.com/go-toolsmith/astcopy v1.0.2 // indirect
	github.com/go-toolsmith/astequal v1.1.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/btree v1.0.1 // indirect
	github.com/google/gnostic v0.5.7-v3refs // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/google/wire v0.5.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.2.3 // indirect
	github.com/googleapis/gax-go/v2 v2.8.0 // indirect
	github.com/gregjones/httpcache v0.0.0-20180305231024-9cad4c3443a7 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.11.3 // indirect
	github.com/hailocab/go-hostpool v0.0.0-20160125115350-e80d13ce29ed // indirect
	github.com/huandu/xstrings v1.4.0 // indirect
	github.com/imdario/mergo v0.3.6 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.0.9 // indirect
	github.com/liggitt/tabwriter v0.0.0-20181228230101-89fcab3d43de // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/moby/spdystream v0.2.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/monochromegane/go-gitignore v0.0.0-20200626010858-205db1a8cc00 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_golang v1.14.0 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/quasilyte/gogrep v0.5.0 // indirect
	github.com/quasilyte/stdinfo v0.0.0-20220114132959-f7386bf02567 // indirect
	github.com/scylladb/go-reflectx v1.0.1 // indirect
	github.com/spf13/cobra v1.6.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/xlab/treeprint v1.1.0 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.11.1 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric v0.33.0 // indirect
	go.opentelemetry.io/proto/otlp v0.19.0 // indirect
	go.starlark.net v0.0.0-20200306205701-8dd3e2ee1dd5 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	golang.org/x/exp/typeparams v0.0.0-20230203172020-98cc5a0785f9 // indirect
	golang.org/x/image v0.6.0 // indirect
	golang.org/x/mod v0.9.0 // indirect
	golang.org/x/term v0.6.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	gomodules.xyz/jsonpatch/v2 v2.2.0 // indirect
	google.golang.org/api v0.114.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/apiextensions-apiserver v0.26.1 // indirect
	k8s.io/component-base v0.26.3 // indirect
	k8s.io/klog/v2 v2.80.1 // indirect
	k8s.io/kube-openapi v0.0.0-20221012153701-172d655c2280 // indirect
	k8s.io/utils v0.0.0-20221128185143-99ec85e7a448 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/kustomize/api v0.12.1 // indirect
	sigs.k8s.io/kustomize/kyaml v0.13.9 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)
