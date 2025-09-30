module github.com/vdaas/vald

go 1.24.5

tool (
	github.com/bufbuild/buf/cmd/buf
	github.com/cockroachdb/crlfmt
	github.com/cweill/gotests/gotests
	github.com/derailed/k9s
	github.com/fatih/gomodifytags
	github.com/go-delve/delve/cmd/dlv
	github.com/google/yamlfmt/cmd/yamlfmt
	github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt
	github.com/haya14busa/goplay/cmd/goplay
	github.com/josharian/impl
	github.com/mfridman/tparse
	github.com/momotaro98/strictgoimports/cmd/strictgoimports
	github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc
	github.com/segmentio/golines
	github.com/stern/stern
	golang.org/x/tools/cmd/goimports
	honnef.co/go/tools/cmd/staticcheck
	mvdan.cc/gofumpt
)

replace (
	cloud.google.com/go => cloud.google.com/go v0.123.0
	cloud.google.com/go/bigquery => cloud.google.com/go/bigquery v1.70.0
	cloud.google.com/go/compute => cloud.google.com/go/compute v1.47.0
	cloud.google.com/go/datastore => cloud.google.com/go/datastore v1.20.0
	cloud.google.com/go/firestore => cloud.google.com/go/firestore v1.18.0
	cloud.google.com/go/iam => cloud.google.com/go/iam v1.5.2
	cloud.google.com/go/kms => cloud.google.com/go/kms v1.23.0
	cloud.google.com/go/monitoring => cloud.google.com/go/monitoring v1.24.2
	cloud.google.com/go/pubsub => cloud.google.com/go/pubsub v1.50.1
	cloud.google.com/go/secretmanager => cloud.google.com/go/secretmanager v1.15.0
	cloud.google.com/go/storage => cloud.google.com/go/storage v1.57.0
	cloud.google.com/go/trace => cloud.google.com/go/trace v1.11.6
	code.cloudfoundry.org/bytefmt => code.cloudfoundry.org/bytefmt v0.53.0
	contrib.go.opencensus.io/exporter/aws => contrib.go.opencensus.io/exporter/aws v0.0.0-20230502192102-15967c811cec
	contrib.go.opencensus.io/exporter/prometheus => contrib.go.opencensus.io/exporter/prometheus v0.4.2
	contrib.go.opencensus.io/integrations/ocsql => contrib.go.opencensus.io/integrations/ocsql v0.1.7
	git.sr.ht/~sbinet/gg => git.sr.ht/~sbinet/gg v0.7.0
	github.com/Azure/azure-amqp-common-go/v3 => github.com/Azure/azure-amqp-common-go/v3 v3.2.3
	github.com/Azure/azure-sdk-for-go => github.com/Azure/azure-sdk-for-go v68.0.0+incompatible
	github.com/Azure/azure-sdk-for-go/sdk/azcore => github.com/Azure/azure-sdk-for-go/sdk/azcore v1.19.1
	github.com/Azure/azure-sdk-for-go/sdk/azidentity => github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.12.0
	github.com/Azure/azure-sdk-for-go/sdk/internal => github.com/Azure/azure-sdk-for-go/sdk/internal v1.11.2
	github.com/Azure/go-amqp => github.com/Azure/go-amqp v1.5.0
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v14.2.1-0.20250128162915-33e12ab7683c+incompatible
	github.com/Azure/go-autorest/autorest => github.com/Azure/go-autorest/autorest v0.11.30
	github.com/Azure/go-autorest/autorest/adal => github.com/Azure/go-autorest/autorest/adal v0.9.24
	github.com/Azure/go-autorest/autorest/date => github.com/Azure/go-autorest/autorest/date v0.3.1
	github.com/Azure/go-autorest/autorest/mocks => github.com/Azure/go-autorest/autorest/mocks v0.4.3
	github.com/Azure/go-autorest/autorest/to => github.com/Azure/go-autorest/autorest/to v0.4.1
	github.com/Azure/go-autorest/logger => github.com/Azure/go-autorest/logger v0.2.2
	github.com/Azure/go-autorest/tracing => github.com/Azure/go-autorest/tracing v0.6.1
	github.com/BurntSushi/toml => github.com/BurntSushi/toml v1.5.0
	github.com/DATA-DOG/go-sqlmock => github.com/DATA-DOG/go-sqlmock v1.5.2
	github.com/GoogleCloudPlatform/cloudsql-proxy => github.com/GoogleCloudPlatform/cloudsql-proxy v1.37.9
	github.com/Masterminds/semver/v3 => github.com/Masterminds/semver/v3 v3.4.0
	github.com/ajstarks/deck => github.com/ajstarks/deck v0.0.0-20250603153621-d300efa64c01
	github.com/ajstarks/deck/generate => github.com/ajstarks/deck/generate v0.0.0-20250603153621-d300efa64c01
	github.com/ajstarks/svgo => github.com/ajstarks/svgo v0.0.0-20211024235047-1546f124cd8b
	github.com/akrylysov/pogreb => github.com/akrylysov/pogreb v0.10.2
	github.com/antihax/optional => github.com/antihax/optional v1.0.0
	github.com/armon/go-socks5 => github.com/armon/go-socks5 v0.0.0-20160902184237-e75332964ef5
	github.com/aws/aws-sdk-go => github.com/aws/aws-sdk-go v1.55.8
	github.com/aws/aws-sdk-go-v2 => github.com/aws/aws-sdk-go-v2 v1.39.2
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream => github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.1
	github.com/aws/aws-sdk-go-v2/config => github.com/aws/aws-sdk-go-v2/config v1.31.12
	github.com/aws/aws-sdk-go-v2/credentials => github.com/aws/aws-sdk-go-v2/credentials v1.18.16
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds => github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.9
	github.com/aws/aws-sdk-go-v2/feature/s3/manager => github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.19.10
	github.com/aws/aws-sdk-go-v2/internal/configsources => github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.9
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 => github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.9
	github.com/aws/aws-sdk-go-v2/internal/ini => github.com/aws/aws-sdk-go-v2/internal/ini v1.8.3
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding => github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.1
	github.com/aws/aws-sdk-go-v2/service/internal/checksum => github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.8.9
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url => github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.9
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared => github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.19.9
	github.com/aws/aws-sdk-go-v2/service/kms => github.com/aws/aws-sdk-go-v2/service/kms v1.45.6
	github.com/aws/aws-sdk-go-v2/service/s3 => github.com/aws/aws-sdk-go-v2/service/s3 v1.88.3
	github.com/aws/aws-sdk-go-v2/service/secretsmanager => github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.39.6
	github.com/aws/aws-sdk-go-v2/service/sns => github.com/aws/aws-sdk-go-v2/service/sns v1.38.5
	github.com/aws/aws-sdk-go-v2/service/sqs => github.com/aws/aws-sdk-go-v2/service/sqs v1.42.8
	github.com/aws/aws-sdk-go-v2/service/ssm => github.com/aws/aws-sdk-go-v2/service/ssm v1.65.1
	github.com/aws/aws-sdk-go-v2/service/sso => github.com/aws/aws-sdk-go-v2/service/sso v1.29.6
	github.com/aws/aws-sdk-go-v2/service/sts => github.com/aws/aws-sdk-go-v2/service/sts v1.38.6
	github.com/aws/smithy-go => github.com/aws/smithy-go v1.23.0
	github.com/benbjohnson/clock => github.com/benbjohnson/clock v1.3.5
	github.com/beorn7/perks => github.com/beorn7/perks v1.0.1
	github.com/bmizerany/assert => github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869
	github.com/boombuler/barcode => github.com/boombuler/barcode v1.1.0
	github.com/buger/jsonparser => github.com/buger/jsonparser v1.1.1
	github.com/cenkalti/backoff/v4 => github.com/cenkalti/backoff/v4 v4.3.0
	github.com/census-instrumentation/opencensus-proto => github.com/census-instrumentation/opencensus-proto v0.4.1
	github.com/cespare/xxhash/v2 => github.com/cespare/xxhash/v2 v2.3.0
	github.com/chzyer/logex => github.com/chzyer/logex v1.2.1
	github.com/chzyer/readline => github.com/chzyer/readline v1.5.1
	github.com/chzyer/test => github.com/chzyer/test v1.0.0
	github.com/cncf/udpa/go => github.com/cncf/udpa/go v0.0.0-20220112060539-c52dc94e7fbe
	github.com/cncf/xds/go => github.com/cncf/xds/go v0.0.0-20250501225837-2ac532fd4443
	github.com/cockroachdb/apd => github.com/cockroachdb/apd v1.1.0
	github.com/coreos/go-systemd/v22 => github.com/coreos/go-systemd/v22 v22.6.0
	github.com/cpuguy83/go-md2man/v2 => github.com/cpuguy83/go-md2man/v2 v2.0.7
	github.com/creack/pty => github.com/creack/pty v1.1.24
	github.com/davecgh/go-spew => github.com/davecgh/go-spew v1.1.1
	github.com/denisenkom/go-mssqldb => github.com/denisenkom/go-mssqldb v0.12.3
	github.com/devigned/tab => github.com/devigned/tab v0.1.1
	github.com/dgryski/go-rendezvous => github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f
	github.com/dnaeon/go-vcr => github.com/dnaeon/go-vcr v1.2.0
	github.com/docopt/docopt-go => github.com/docopt/docopt-go v0.0.0-20180111231733-ee0de3bc6815
	github.com/dustin/go-humanize => github.com/dustin/go-humanize v1.0.1
	github.com/emicklei/go-restful/v3 => github.com/emicklei/go-restful/v3 v3.13.0
	github.com/envoyproxy/go-control-plane => github.com/envoyproxy/go-control-plane v0.13.4
	github.com/envoyproxy/protoc-gen-validate => github.com/envoyproxy/protoc-gen-validate v1.2.1
	github.com/evanphx/json-patch => github.com/evanphx/json-patch v0.5.2
	github.com/fogleman/gg => github.com/fogleman/gg v1.3.0
	github.com/fortytw2/leaktest => github.com/fortytw2/leaktest v1.3.0
	github.com/frankban/quicktest => github.com/frankban/quicktest v1.14.6
	github.com/fsnotify/fsnotify => github.com/fsnotify/fsnotify v1.9.0
	github.com/gin-contrib/sse => github.com/gin-contrib/sse v1.1.0
	github.com/gin-gonic/gin => github.com/gin-gonic/gin v1.11.0
	github.com/go-errors/errors => github.com/go-errors/errors v1.5.1
	github.com/go-fonts/dejavu => github.com/go-fonts/dejavu v0.3.4
	github.com/go-fonts/latin-modern => github.com/go-fonts/latin-modern v0.3.3
	github.com/go-fonts/liberation => github.com/go-fonts/liberation v0.3.3
	github.com/go-fonts/stix => github.com/go-fonts/stix v0.2.2
	github.com/go-gl/gl => github.com/go-gl/gl v0.0.0-20231021071112-07e5d0ea2e71
	github.com/go-gl/glfw/v3.3/glfw => github.com/go-gl/glfw/v3.3/glfw v0.0.0-20250301202403-da16c1255728
	github.com/go-kit/log => github.com/go-kit/log v0.2.1
	github.com/go-latex/latex => github.com/go-latex/latex v0.0.0-20250304174226-2790903426af
	github.com/go-logfmt/logfmt => github.com/go-logfmt/logfmt v0.6.0
	github.com/go-logr/logr => github.com/go-logr/logr v1.4.3
	github.com/go-logr/stdr => github.com/go-logr/stdr v1.2.2
	github.com/go-logr/zapr => github.com/go-logr/zapr v1.3.0
	github.com/go-openapi/jsonpointer => github.com/go-openapi/jsonpointer v0.22.1
	github.com/go-openapi/jsonreference => github.com/go-openapi/jsonreference v0.21.2
	github.com/go-openapi/swag => github.com/go-openapi/swag v0.25.1
	github.com/go-pdf/fpdf => github.com/go-pdf/fpdf v1.4.3
	github.com/go-playground/assert/v2 => github.com/go-playground/assert/v2 v2.2.0
	github.com/go-playground/locales => github.com/go-playground/locales v0.14.1
	github.com/go-playground/universal-translator => github.com/go-playground/universal-translator v0.18.1
	github.com/go-playground/validator/v10 => github.com/go-playground/validator/v10 v10.27.0
	github.com/go-redis/redis/v8 => github.com/go-redis/redis/v8 v8.11.5
	github.com/go-sql-driver/mysql => github.com/go-sql-driver/mysql v1.9.3
	github.com/go-task/slim-sprig => github.com/go-task/slim-sprig v2.20.0+incompatible
	github.com/go-toolsmith/astcopy => github.com/go-toolsmith/astcopy v1.1.0
	github.com/go-toolsmith/astequal => github.com/go-toolsmith/astequal v1.2.0
	github.com/go-toolsmith/strparse => github.com/go-toolsmith/strparse v1.1.0
	github.com/gobwas/httphead => github.com/gobwas/httphead v0.1.0
	github.com/gobwas/pool => github.com/gobwas/pool v0.2.1
	github.com/gobwas/ws => github.com/gobwas/ws v1.4.0
	github.com/goccy/go-json => github.com/goccy/go-json v0.10.5
	github.com/gocql/gocql => github.com/gocql/gocql v1.7.0
	github.com/gocraft/dbr/v2 => github.com/gocraft/dbr/v2 v2.7.7
	github.com/godbus/dbus/v5 => github.com/godbus/dbus/v5 v5.1.0
	github.com/gofrs/uuid => github.com/gofrs/uuid v4.4.0+incompatible
	github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.2
	github.com/golang-jwt/jwt/v4 => github.com/golang-jwt/jwt/v4 v4.5.2
	github.com/golang-sql/civil => github.com/golang-sql/civil v0.0.0-20220223132316-b832511892a9
	github.com/golang-sql/sqlexp => github.com/golang-sql/sqlexp v0.1.0
	github.com/golang/freetype => github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/golang/glog => github.com/golang/glog v1.2.5
	github.com/golang/groupcache => github.com/golang/groupcache v0.0.0-20241129210726-2c02b8208cf8
	github.com/golang/mock => github.com/golang/mock v1.6.0
	github.com/golang/protobuf => github.com/golang/protobuf v1.5.4
	github.com/golang/snappy => github.com/golang/snappy v1.0.0
	github.com/google/btree => github.com/google/btree v1.1.3
	github.com/google/gnostic => github.com/google/gnostic v0.7.1
	github.com/google/go-cmp => github.com/google/go-cmp v0.7.0
	github.com/google/go-replayers/grpcreplay => github.com/google/go-replayers/grpcreplay v1.3.0
	github.com/google/go-replayers/httpreplay => github.com/google/go-replayers/httpreplay v1.2.0
	github.com/google/gofuzz => github.com/google/gofuzz v1.2.0
	github.com/google/martian => github.com/google/martian v2.1.0+incompatible
	github.com/google/martian/v3 => github.com/google/martian/v3 v3.3.3
	github.com/google/pprof => github.com/google/pprof v0.0.0-20250923004556-9e5a51aed1e8
	github.com/google/shlex => github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510
	github.com/google/subcommands => github.com/google/subcommands v1.2.0
	github.com/google/uuid => github.com/google/uuid v1.6.0
	github.com/google/wire => github.com/google/wire v0.7.0
	github.com/googleapis/gax-go/v2 => github.com/googleapis/gax-go/v2 v2.15.0
	github.com/gorilla/mux => github.com/gorilla/mux v1.8.1
	github.com/gorilla/websocket => github.com/gorilla/websocket v1.5.3
	github.com/grafana/grafana-foundation-sdk/go => github.com/grafana/grafana-foundation-sdk/go v0.0.0-20250916165541-37dba040e63e
	github.com/grafana/pyroscope-go/godeltaprof => github.com/grafana/pyroscope-go/godeltaprof v0.1.9
	github.com/gregjones/httpcache => github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79
	github.com/grpc-ecosystem/grpc-gateway/v2 => github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.2
	github.com/hailocab/go-hostpool => github.com/kpango/go-hostpool v0.0.0-20210303030322-aab80263dcd0
	github.com/hanwen/go-fuse/v2 => github.com/hanwen/go-fuse/v2 v2.8.0
	github.com/hashicorp/go-uuid => github.com/hashicorp/go-uuid v1.0.3
	github.com/hashicorp/go-version => github.com/hashicorp/go-version v1.7.0
	github.com/iancoleman/strcase => github.com/iancoleman/strcase v0.3.0
	github.com/ianlancetaylor/demangle => github.com/ianlancetaylor/demangle v0.0.0-20250628045327-2d64ad6b7ec5
	github.com/inconshreveable/mousetrap => github.com/inconshreveable/mousetrap v1.1.0
	github.com/jackc/chunkreader/v2 => github.com/jackc/chunkreader/v2 v2.0.1
	github.com/jackc/pgconn => github.com/jackc/pgconn v1.14.3
	github.com/jackc/pgio => github.com/jackc/pgio v1.0.0
	github.com/jackc/pgmock => github.com/jackc/pgmock v0.0.0-20210724152146-4ad1a8207f65
	github.com/jackc/pgpassfile => github.com/jackc/pgpassfile v1.0.0
	github.com/jackc/pgproto3/v2 => github.com/jackc/pgproto3/v2 v2.3.3
	github.com/jackc/pgservicefile => github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761
	github.com/jackc/pgtype => github.com/jackc/pgtype v1.14.4
	github.com/jackc/pgx/v4 => github.com/jackc/pgx/v4 v4.18.3
	github.com/jackc/puddle => github.com/jackc/puddle v1.3.0
	github.com/jessevdk/go-flags => github.com/jessevdk/go-flags v1.6.1
	github.com/jmespath/go-jmespath => github.com/jmespath/go-jmespath v0.4.0
	github.com/jmespath/go-jmespath/internal/testify => github.com/jmespath/go-jmespath/internal/testify v1.5.1
	github.com/jmoiron/sqlx => github.com/jmoiron/sqlx v1.4.0
	github.com/joho/godotenv => github.com/joho/godotenv v1.5.1
	github.com/josharian/intern => github.com/josharian/intern v1.0.0
	github.com/json-iterator/go => github.com/json-iterator/go v1.1.12
	github.com/jstemmer/go-junit-report => github.com/jstemmer/go-junit-report v1.0.0
	github.com/kisielk/errcheck => github.com/kisielk/errcheck v1.9.0
	github.com/kisielk/gotool => github.com/kisielk/gotool v1.0.0
	github.com/klauspost/compress => github.com/klauspost/compress v1.18.1-0.20250921124417-3c0d30844ced
	github.com/klauspost/cpuid/v2 => github.com/klauspost/cpuid/v2 v2.3.0
	github.com/kpango/fastime => github.com/kpango/fastime v1.1.10
	github.com/kpango/fuid => github.com/kpango/fuid v0.0.0-20221203053508-503b5ad89aa1
	github.com/kpango/gache/v2 => github.com/kpango/gache/v2 v2.1.1
	github.com/kpango/glg => github.com/kpango/glg v1.6.15
	github.com/kr/fs => github.com/kr/fs v0.1.0
	github.com/kr/pretty => github.com/kr/pretty v0.3.1
	github.com/kr/text => github.com/kr/text v0.2.0
	github.com/kubernetes-csi/external-snapshotter/client/v6 => github.com/kubernetes-csi/external-snapshotter/client/v6 v6.3.0
	github.com/kylelemons/godebug => github.com/kylelemons/godebug v1.1.0
	github.com/leanovate/gopter => github.com/leanovate/gopter v0.2.11
	github.com/leodido/go-urn => github.com/leodido/go-urn v1.4.0
	github.com/lib/pq => github.com/lib/pq v1.10.9
	github.com/liggitt/tabwriter => github.com/liggitt/tabwriter v0.0.0-20181228230101-89fcab3d43de
	github.com/lucasb-eyer/go-colorful => github.com/lucasb-eyer/go-colorful v1.3.0
	github.com/mailru/easyjson => github.com/mailru/easyjson v0.9.1
	github.com/mattn/go-colorable => github.com/mattn/go-colorable v0.1.14
	github.com/mattn/go-isatty => github.com/mattn/go-isatty v0.0.20
	github.com/mattn/go-sqlite3 => github.com/mattn/go-sqlite3 v1.14.32
	github.com/matttproud/golang_protobuf_extensions => github.com/matttproud/golang_protobuf_extensions v1.0.4
	github.com/mitchellh/colorstring => github.com/mitchellh/colorstring v0.0.0-20190213212951-d06e56a500db
	github.com/moby/spdystream => github.com/moby/spdystream v0.5.0
	github.com/moby/sys/mountinfo => github.com/moby/sys/mountinfo v0.7.2
	github.com/modern-go/concurrent => github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd
	github.com/modern-go/reflect2 => github.com/modern-go/reflect2 v1.0.2
	github.com/modocache/gover => github.com/modocache/gover v0.0.0-20171022184752-b58185e213c5
	github.com/monochromegane/go-gitignore => github.com/monochromegane/go-gitignore v0.0.0-20200626010858-205db1a8cc00
	github.com/montanaflynn/stats => github.com/montanaflynn/stats v0.7.1
	github.com/munnerz/goautoneg => github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822
	github.com/niemeyer/pretty => github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e
	github.com/nxadm/tail => github.com/nxadm/tail v1.4.11
	github.com/onsi/ginkgo => github.com/onsi/ginkgo v1.16.5
	github.com/onsi/ginkgo/v2 => github.com/onsi/ginkgo/v2 v2.25.3
	github.com/onsi/gomega => github.com/onsi/gomega v1.38.2
	github.com/peterbourgon/diskv => github.com/peterbourgon/diskv v2.0.1+incompatible
	github.com/phpdave11/gofpdf => github.com/phpdave11/gofpdf v1.4.3
	github.com/phpdave11/gofpdi => github.com/phpdave11/gofpdi v1.0.15
	github.com/pierrec/cmdflag => github.com/pierrec/cmdflag v0.0.2
	github.com/pierrec/lz4/v3 => github.com/pierrec/lz4/v3 v3.3.5
	github.com/pkg/browser => github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c
	github.com/pkg/errors => github.com/pkg/errors v0.9.1
	github.com/pkg/sftp => github.com/pkg/sftp v1.13.9
	github.com/pmezard/go-difflib => github.com/pmezard/go-difflib v1.0.0
	github.com/prashantv/gostub => github.com/prashantv/gostub v1.1.0
	github.com/prometheus/client_golang => github.com/prometheus/client_golang v1.23.2
	github.com/prometheus/client_model => github.com/prometheus/client_model v0.6.2
	github.com/prometheus/procfs => github.com/prometheus/procfs v0.17.0
	github.com/prometheus/prometheus => github.com/prometheus/prometheus v1.99.0
	github.com/quasilyte/go-ruleguard => github.com/quasilyte/go-ruleguard v0.4.5
	github.com/quasilyte/go-ruleguard/dsl => github.com/quasilyte/go-ruleguard/dsl v0.3.23
	github.com/quasilyte/gogrep => github.com/quasilyte/gogrep v0.5.0
	github.com/quasilyte/stdinfo => github.com/quasilyte/stdinfo v0.0.0-20220114132959-f7386bf02567
	github.com/rogpeppe/fastuuid => github.com/rogpeppe/fastuuid v1.2.0
	github.com/rogpeppe/go-internal => github.com/rogpeppe/go-internal v1.14.1
	github.com/rs/xid => github.com/rs/xid v1.6.0
	github.com/rs/zerolog => github.com/rs/zerolog v1.34.0
	github.com/russross/blackfriday/v2 => github.com/russross/blackfriday/v2 v2.1.0
	github.com/ruudk/golang-pdf417 => github.com/ruudk/golang-pdf417 v0.0.0-20201230142125-a7e3863a1245
	github.com/schollz/progressbar/v2 => github.com/schollz/progressbar/v2 v2.15.0
	github.com/scylladb/go-reflectx => github.com/scylladb/go-reflectx v1.0.1
	github.com/scylladb/gocqlx => github.com/scylladb/gocqlx v1.5.0
	github.com/sergi/go-diff => github.com/sergi/go-diff v1.4.0
	github.com/shopspring/decimal => github.com/shopspring/decimal v1.4.0
	github.com/shurcooL/httpfs => github.com/shurcooL/httpfs v0.0.0-20230704072500-f1e31cf0ba5c
	github.com/shurcooL/vfsgen => github.com/shurcooL/vfsgen v0.0.0-20230704071429-0000e147ea92
	github.com/sirupsen/logrus => github.com/sirupsen/logrus v1.9.3
	github.com/spf13/afero => github.com/spf13/afero v1.15.0
	github.com/spf13/cobra => github.com/spf13/cobra v1.10.1
	github.com/spf13/pflag => github.com/spf13/pflag v1.0.10
	github.com/stern/stern => github.com/stern/stern v1.32.0
	github.com/stoewer/go-strcase => github.com/stoewer/go-strcase v1.3.1
	github.com/stretchr/objx => github.com/stretchr/objx v0.5.2
	github.com/stretchr/testify => github.com/stretchr/testify v1.11.1
	github.com/ugorji/go/codec => github.com/ugorji/go/codec v1.3.0
	github.com/xeipuuv/gojsonpointer => github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb
	github.com/xeipuuv/gojsonreference => github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415
	github.com/xeipuuv/gojsonschema => github.com/xeipuuv/gojsonschema v1.2.0
	github.com/xlab/treeprint => github.com/xlab/treeprint v1.2.0
	github.com/zeebo/assert => github.com/zeebo/assert v1.3.1
	github.com/zeebo/xxh3 => github.com/zeebo/xxh3 v1.0.2
	go.etcd.io/bbolt => go.etcd.io/bbolt v1.4.3
	go.opencensus.io => go.opencensus.io v0.24.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc => go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.63.0
	go.opentelemetry.io/otel => go.opentelemetry.io/otel v1.38.0
	go.opentelemetry.io/otel/exporters/otlp/internal/retry => go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.17.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric => go.opentelemetry.io/otel/exporters/otlp/otlpmetric v0.43.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc => go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.38.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace => go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.38.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc => go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.38.0
	go.opentelemetry.io/otel/metric => go.opentelemetry.io/otel/metric v1.38.0
	go.opentelemetry.io/otel/sdk => go.opentelemetry.io/otel/sdk v1.38.0
	go.opentelemetry.io/otel/sdk/metric => go.opentelemetry.io/otel/sdk/metric v1.38.0
	go.opentelemetry.io/otel/trace => go.opentelemetry.io/otel/trace v1.38.0
	go.opentelemetry.io/proto/otlp => go.opentelemetry.io/proto/otlp v1.8.0
	go.starlark.net => go.starlark.net v0.0.0-20250906160240-bf296ed553ea
	go.uber.org/atomic => go.uber.org/atomic v1.11.0
	go.uber.org/automaxprocs => go.uber.org/automaxprocs v1.6.0
	go.uber.org/goleak => go.uber.org/goleak v1.3.0
	go.uber.org/multierr => go.uber.org/multierr v1.11.0
	go.uber.org/zap => go.uber.org/zap v1.27.0
	gocloud.dev => gocloud.dev v0.43.0
	golang.org/x/crypto => golang.org/x/crypto v0.42.0
	golang.org/x/exp => golang.org/x/exp v0.0.0-20250911091902-df9299821621
	golang.org/x/exp/typeparams => golang.org/x/exp/typeparams v0.0.0-20250911091902-df9299821621
	golang.org/x/image => golang.org/x/image v0.31.0
	golang.org/x/lint => golang.org/x/lint v0.0.0-20241112194109-818c5a804067
	golang.org/x/mobile => golang.org/x/mobile v0.0.0-20250911085028-6912353760cf
	golang.org/x/mod => golang.org/x/mod v0.28.0
	golang.org/x/net => golang.org/x/net v0.44.0
	golang.org/x/oauth2 => golang.org/x/oauth2 v0.31.0
	golang.org/x/sync => golang.org/x/sync v0.17.0
	golang.org/x/sys => golang.org/x/sys v0.36.0
	golang.org/x/term => golang.org/x/term v0.35.0
	golang.org/x/text => golang.org/x/text v0.29.0
	golang.org/x/time => golang.org/x/time v0.13.0
	golang.org/x/tools => golang.org/x/tools v0.37.0
	golang.org/x/xerrors => golang.org/x/xerrors v0.0.0-20240903120638-7835f813f4da
	gomodules.xyz/jsonpatch/v2 => gomodules.xyz/jsonpatch/v2 v2.5.0
	gonum.org/v1/gonum => gonum.org/v1/gonum v0.16.0
	gonum.org/v1/hdf5 => gonum.org/v1/hdf5 v0.0.0-20210714002203-8c5d23bc6946
	gonum.org/v1/plot => gonum.org/v1/plot v0.16.0
	google.golang.org/api => google.golang.org/api v0.250.0
	google.golang.org/appengine => google.golang.org/appengine v1.6.8
	google.golang.org/genproto => google.golang.org/genproto v0.0.0-20250929231259-57b25ae835d4
	google.golang.org/genproto/googleapis/api => google.golang.org/genproto/googleapis/api v0.0.0-20250929231259-57b25ae835d4
	google.golang.org/genproto/googleapis/rpc => google.golang.org/genproto/googleapis/rpc v0.0.0-20250929231259-57b25ae835d4
	google.golang.org/grpc => google.golang.org/grpc v1.75.1
	google.golang.org/grpc/cmd/protoc-gen-go-grpc => google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.5.1
	google.golang.org/protobuf => google.golang.org/protobuf v1.36.9
	gopkg.in/check.v1 => gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c
	gopkg.in/inconshreveable/log15.v2 => gopkg.in/inconshreveable/log15.v2 v2.16.0
	gopkg.in/inf.v0 => gopkg.in/inf.v0 v0.9.1
	gopkg.in/tomb.v1 => gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7
	gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 => gopkg.in/yaml.v3 v3.0.1
	honnef.co/go/tools => honnef.co/go/tools v0.6.1
	nhooyr.io/websocket => nhooyr.io/websocket v1.8.17
	rsc.io/pdf => rsc.io/pdf v0.1.1
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.21.0
	sigs.k8s.io/json => sigs.k8s.io/json v0.0.0-20250730193827-2d320260d730
	sigs.k8s.io/kustomize => sigs.k8s.io/kustomize v2.0.3+incompatible
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.9-20250912141014-52f32327d4b0.1
	cloud.google.com/go/storage v1.56.0
	code.cloudfoundry.org/bytefmt v0.0.0-20190710193110-1eb035ffe2b6
	github.com/akrylysov/pogreb v0.0.0-00010101000000-000000000000
	github.com/aws/aws-sdk-go v1.55.7
	github.com/felixge/fgprof v0.9.5
	github.com/fsnotify/fsnotify v1.9.0
	github.com/go-redis/redis/v8 v8.0.0-00010101000000-000000000000
	github.com/go-sql-driver/mysql v1.9.3
	github.com/goccy/go-json v0.10.2
	github.com/gocql/gocql v0.0.0-20200131111108-92af2e088537
	github.com/gocraft/dbr/v2 v2.0.0-00010101000000-000000000000
	github.com/google/go-cmp v0.7.0
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1
	github.com/grafana/grafana-foundation-sdk/go v0.0.0-00010101000000-000000000000
	github.com/grafana/promql-builder/go v0.0.0-20250916111012-8fa9625b89a3
	github.com/grafana/pyroscope-go/godeltaprof v0.0.0-00010101000000-000000000000
	github.com/hashicorp/go-version v1.7.0
	github.com/klauspost/compress v1.18.0
	github.com/kpango/fastime v1.1.9
	github.com/kpango/gache/v2 v2.0.0-00010101000000-000000000000
	github.com/kpango/glg v1.6.15
	github.com/kubernetes-csi/external-snapshotter/client/v6 v6.0.0-00010101000000-000000000000
	github.com/leanovate/gopter v0.0.0-00010101000000-000000000000
	github.com/lucasb-eyer/go-colorful v1.2.0
	github.com/pierrec/lz4/v3 v3.0.0-00010101000000-000000000000
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10
	github.com/quasilyte/go-ruleguard v0.0.0-00010101000000-000000000000
	github.com/quasilyte/go-ruleguard/dsl v0.3.22
	github.com/quic-go/quic-go v0.54.0
	github.com/scylladb/gocqlx v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.11.1
	github.com/unum-cloud/usearch/golang v0.0.0-20250904130807-fd6279af6bc2
	github.com/zeebo/xxh3 v1.0.2
	go.etcd.io/bbolt v1.4.2
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.62.0
	go.opentelemetry.io/otel v1.38.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.37.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.38.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.37.0
	go.opentelemetry.io/otel/metric v1.38.0
	go.opentelemetry.io/otel/sdk v1.38.0
	go.opentelemetry.io/otel/sdk/metric v1.38.0
	go.opentelemetry.io/otel/trace v1.38.0
	go.uber.org/automaxprocs v1.6.0
	go.uber.org/goleak v1.3.0
	go.uber.org/ratelimit v0.3.1
	go.uber.org/zap v1.27.0
	gocloud.dev v0.0.0-00010101000000-000000000000
	golang.org/x/net v0.44.0
	golang.org/x/oauth2 v0.31.0
	golang.org/x/sync v0.17.0
	golang.org/x/sys v0.36.0
	golang.org/x/text v0.29.0
	golang.org/x/time v0.13.0
	golang.org/x/tools v0.37.0
	gonum.org/v1/hdf5 v0.0.0-00010101000000-000000000000
	gonum.org/v1/plot v0.15.2
	google.golang.org/genproto/googleapis/api v0.0.0-20250922171735-9219d122eba9
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250922171735-9219d122eba9
	google.golang.org/grpc v1.75.1
	google.golang.org/protobuf v1.36.9
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.1
	k8s.io/api v0.34.1
	k8s.io/apimachinery v0.34.1
	k8s.io/cli-runtime v0.34.1
	k8s.io/client-go v0.34.1
	k8s.io/metrics v0.34.1
	k8s.io/utils v0.0.0-20250604170112-4c0f3b243397
	sigs.k8s.io/controller-runtime v0.0.0-00010101000000-000000000000
	sigs.k8s.io/yaml v1.6.0
)

require (
	buf.build/gen/go/bufbuild/bufplugin/protocolbuffers/go v1.36.9-20250718181942-e35f9b667443.1 // indirect
	buf.build/gen/go/bufbuild/registry/connectrpc/go v1.18.1-20250903170917-c4be0f57e197.1 // indirect
	buf.build/gen/go/bufbuild/registry/protocolbuffers/go v1.36.9-20250903170917-c4be0f57e197.1 // indirect
	buf.build/gen/go/pluginrpc/pluginrpc/protocolbuffers/go v1.36.8-20241007202033-cf42259fcbfc.1 // indirect
	buf.build/go/app v0.1.0 // indirect
	buf.build/go/bufplugin v0.9.0 // indirect
	buf.build/go/interrupt v1.1.0 // indirect
	buf.build/go/protovalidate v1.0.0 // indirect
	buf.build/go/protoyaml v0.6.0 // indirect
	buf.build/go/spdx v0.2.0 // indirect
	buf.build/go/standard v0.1.0 // indirect
	cel.dev/expr v0.24.0 // indirect
	cloud.google.com/go v0.121.6 // indirect
	cloud.google.com/go/auth v0.16.5 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.8 // indirect
	cloud.google.com/go/compute/metadata v0.8.4 // indirect
	cloud.google.com/go/iam v1.5.2 // indirect
	cloud.google.com/go/monitoring v1.24.2 // indirect
	codeberg.org/go-fonts/liberation v0.5.0 // indirect
	codeberg.org/go-latex/latex v0.1.0 // indirect
	codeberg.org/go-pdf/fpdf v0.11.1 // indirect
	connectrpc.com/connect v1.18.1 // indirect
	connectrpc.com/otelconnect v0.8.0 // indirect
	dario.cat/mergo v1.0.1 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	git.sr.ht/~sbinet/gg v0.6.0 // indirect
	github.com/AdaLogics/go-fuzz-headers v0.0.0-20230811130428-ced1acdcaa24 // indirect
	github.com/AdamKorcz/go-118-fuzz-build v0.0.0-20230306123547-8075edf89bb0 // indirect
	github.com/AlecAivazis/survey/v2 v2.3.7 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20250102033503-faa5f7b0171c // indirect
	github.com/BurntSushi/toml v1.5.0 // indirect
	github.com/CycloneDX/cyclonedx-go v0.9.2 // indirect
	github.com/DataDog/zstd v1.5.5 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp v1.29.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric v0.53.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/internal/resourcemapping v0.53.0 // indirect
	github.com/IGLOU-EU/go-wildcard v1.0.3 // indirect
	github.com/MakeNowJust/heredoc v1.0.0 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver v1.4.2 // indirect
	github.com/Masterminds/semver/v3 v3.4.0 // indirect
	github.com/Masterminds/sprig v2.15.0+incompatible // indirect
	github.com/Masterminds/sprig/v3 v3.3.0 // indirect
	github.com/Masterminds/squirrel v1.5.4 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/Microsoft/hcsshim v0.11.7 // indirect
	github.com/OneOfOne/xxhash v1.2.8 // indirect
	github.com/ProtonMail/go-crypto v1.3.0 // indirect
	github.com/STARRY-S/zip v0.2.1 // indirect
	github.com/acobaugh/osrelease v0.1.0 // indirect
	github.com/adrg/xdg v0.5.3 // indirect
	github.com/agext/levenshtein v1.2.1 // indirect
	github.com/ajstarks/svgo v0.0.0-20211024235047-1546f124cd8b // indirect
	github.com/alecthomas/kingpin/v2 v2.4.0 // indirect
	github.com/alecthomas/units v0.0.0-20240927000941-0f3dac36c52b // indirect
	github.com/anchore/archiver/v3 v3.5.3-0.20241210171143-5b1d8d1c7c51 // indirect
	github.com/anchore/clio v0.0.0-20250408180537-ec8fa27f0d9f // indirect
	github.com/anchore/fangs v0.0.0-20250402135612-96e29e45f3fe // indirect
	github.com/anchore/go-collections v0.0.0-20240216171411-9321230ce537 // indirect
	github.com/anchore/go-homedir v0.0.0-20250319154043-c29668562e4d // indirect
	github.com/anchore/go-logger v0.0.0-20250318195838-07ae343dd722 // indirect
	github.com/anchore/go-lzo v0.1.0 // indirect
	github.com/anchore/go-macholibre v0.0.0-20220308212642-53e6d0aaf6fb // indirect
	github.com/anchore/go-rpmdb v0.0.0-20250516171929-f77691e1faec // indirect
	github.com/anchore/go-struct-converter v0.0.0-20221118182256-c68fdcfa2092 // indirect
	github.com/anchore/go-sync v0.0.0-20250326131806-4eda43a485b6 // indirect
	github.com/anchore/go-version v1.2.2-0.20210903204242-51efa5b487c4 // indirect
	github.com/anchore/grype v0.96.0 // indirect
	github.com/anchore/packageurl-go v0.1.1-0.20250220190351-d62adb6e1115 // indirect
	github.com/anchore/stereoscope v0.1.10 // indirect
	github.com/anchore/syft v1.33.0 // indirect
	github.com/andybalholm/brotli v1.1.2-0.20250424173009-453214e765f3 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.1 // indirect
	github.com/aokoli/goutils v1.0.1 // indirect
	github.com/apparentlymart/go-textseg/v15 v15.0.0 // indirect
	github.com/aquasecurity/go-pep440-version v0.0.1 // indirect
	github.com/aquasecurity/go-version v0.0.1 // indirect
	github.com/araddon/dateparse v0.0.0-20210429162001-6b43995a97de // indirect
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2 // indirect
	github.com/atotto/clipboard v0.1.4 // indirect
	github.com/aws/aws-sdk-go-v2 v1.39.2 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.1 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.29.17 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.18.16 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.9 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.9 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.9 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.8.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.19.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/s3 v1.84.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.29.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.35.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.38.6 // indirect
	github.com/aws/smithy-go v1.23.0 // indirect
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/becheran/wildmatch-go v1.0.0 // indirect
	github.com/benbjohnson/clock v1.3.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bgentry/go-netrc v0.0.0-20140422174119-9fd32a8b3d3d // indirect
	github.com/bitnami/go-version v0.0.0-20250131085805-b1f57a8634ef // indirect
	github.com/blakesmith/ar v0.0.0-20190502131153-809d4375e1fb // indirect
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/bmatcuk/doublestar/v2 v2.0.4 // indirect
	github.com/bmatcuk/doublestar/v4 v4.9.1 // indirect
	github.com/bodgit/plumbing v1.3.0 // indirect
	github.com/bodgit/sevenzip v1.6.0 // indirect
	github.com/bodgit/windows v1.0.1 // indirect
	github.com/bufbuild/buf v1.57.2 // indirect
	github.com/bufbuild/protocompile v0.14.1 // indirect
	github.com/bufbuild/protoplugin v0.0.0-20250218205857-750e09ce93e1 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cenkalti/backoff/v5 v5.0.3 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/chai2010/gettext-go v1.0.2 // indirect
	github.com/charmbracelet/colorprofile v0.3.1 // indirect
	github.com/charmbracelet/lipgloss v1.1.0 // indirect
	github.com/charmbracelet/x/ansi v0.10.1 // indirect
	github.com/charmbracelet/x/cellbuf v0.0.13 // indirect
	github.com/charmbracelet/x/term v0.2.1 // indirect
	github.com/cilium/ebpf v0.11.0 // indirect
	github.com/cloudflare/circl v1.6.1 // indirect
	github.com/cncf/xds/go v0.0.0-20250501225837-2ac532fd4443 // indirect
	github.com/cockroachdb/crlfmt v0.3.0 // indirect
	github.com/cockroachdb/gostdlib v1.19.0 // indirect
	github.com/containerd/cgroups v1.1.0 // indirect
	github.com/containerd/containerd v1.7.28 // indirect
	github.com/containerd/containerd/api v1.8.0 // indirect
	github.com/containerd/continuity v0.4.4 // indirect
	github.com/containerd/errdefs v1.0.0 // indirect
	github.com/containerd/errdefs/pkg v0.3.0 // indirect
	github.com/containerd/fifo v1.1.0 // indirect
	github.com/containerd/log v0.1.0 // indirect
	github.com/containerd/platforms v0.2.1 // indirect
	github.com/containerd/stargz-snapshotter/estargz v0.17.0 // indirect
	github.com/containerd/ttrpc v1.2.7 // indirect
	github.com/containerd/typeurl/v2 v2.2.2 // indirect
	github.com/cosiner/argv v0.1.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.7 // indirect
	github.com/cweill/gotests v1.6.0 // indirect
	github.com/cyphar/filepath-securejoin v0.4.1 // indirect
	github.com/dave/dst v0.27.3 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/deitch/magic v0.0.0-20230404182410-1ff89d7342da // indirect
	github.com/derailed/k9s v0.50.13 // indirect
	github.com/derailed/tcell/v2 v2.3.1-rc.4 // indirect
	github.com/derailed/tview v0.8.5 // indirect
	github.com/derekparker/trie v0.0.0-20230829180723-39f4de51ef7d // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/diskfs/go-diskfs v1.7.0 // indirect
	github.com/distribution/reference v0.6.0 // indirect
	github.com/docker/cli v28.4.0+incompatible // indirect
	github.com/docker/distribution v2.8.3+incompatible // indirect
	github.com/docker/docker v28.4.0+incompatible // indirect
	github.com/docker/docker-credential-helpers v0.9.3 // indirect
	github.com/docker/go-connections v0.6.0 // indirect
	github.com/docker/go-events v0.0.0-20190806004212-e31b211e4f1c // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/dsnet/compress v0.0.2-0.20230904184137-39efe44ab707 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/elliotchance/phpserialize v1.4.0 // indirect
	github.com/emicklei/go-restful/v3 v3.12.2 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/envoyproxy/go-control-plane/envoy v1.32.4 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.2.1 // indirect
	github.com/evanphx/json-patch v5.9.11+incompatible // indirect
	github.com/evanphx/json-patch/v5 v5.9.11 // indirect
	github.com/exponent-io/jsonpath v0.0.0-20210407135951-1de76d718b3f // indirect
	github.com/facebookincubator/nvdtools v0.1.5 // indirect
	github.com/fatih/camelcase v1.0.0 // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/fatih/gomodifytags v1.17.0 // indirect
	github.com/fatih/structtag v1.2.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/fvbommel/sortorder v1.1.0 // indirect
	github.com/fxamacker/cbor/v2 v2.9.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.9 // indirect
	github.com/gdamore/encoding v1.0.1 // indirect
	github.com/github/go-spdx/v2 v2.3.3 // indirect
	github.com/glebarez/go-sqlite v1.22.0 // indirect
	github.com/glebarez/sqlite v1.11.0 // indirect
	github.com/go-chi/chi/v5 v5.2.3 // indirect
	github.com/go-delve/delve v1.25.2 // indirect
	github.com/go-delve/liner v1.2.3-0.20231231155935-4726ab1d7f62 // indirect
	github.com/go-errors/errors v1.5.1 // indirect
	github.com/go-git/gcfg v1.5.1-0.20230307220236-3a3c6141e376 // indirect
	github.com/go-git/go-billy/v5 v5.6.2 // indirect
	github.com/go-git/go-git/v5 v5.16.2 // indirect
	github.com/go-gorp/gorp/v3 v3.1.0 // indirect
	github.com/go-jose/go-jose/v4 v4.1.1 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-openapi/jsonpointer v0.22.1 // indirect
	github.com/go-openapi/jsonreference v0.21.0 // indirect
	github.com/go-openapi/swag v0.23.0 // indirect
	github.com/go-openapi/swag/cmdutils v0.25.1 // indirect
	github.com/go-openapi/swag/conv v0.25.1 // indirect
	github.com/go-openapi/swag/fileutils v0.25.1 // indirect
	github.com/go-openapi/swag/jsonname v0.25.1 // indirect
	github.com/go-openapi/swag/jsonutils v0.25.1 // indirect
	github.com/go-openapi/swag/loading v0.25.1 // indirect
	github.com/go-openapi/swag/mangling v0.25.1 // indirect
	github.com/go-openapi/swag/netutils v0.25.1 // indirect
	github.com/go-openapi/swag/stringutils v0.25.1 // indirect
	github.com/go-openapi/swag/typeutils v0.25.1 // indirect
	github.com/go-openapi/swag/yamlutils v0.25.1 // indirect
	github.com/go-restruct/restruct v1.2.0-alpha // indirect
	github.com/go-toolsmith/astcopy v1.0.2 // indirect
	github.com/go-toolsmith/astequal v1.1.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/goccy/go-yaml v1.18.0 // indirect
	github.com/gofrs/flock v0.12.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/gohugoio/hashstructure v0.5.0 // indirect
	github.com/golang-sql/sqlexp v0.0.0-00010101000000-000000000000 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golang/groupcache v0.0.0-20241129210726-2c02b8208cf8 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/btree v1.1.3 // indirect
	github.com/google/cel-go v0.26.1 // indirect
	github.com/google/gnostic-models v0.7.0 // indirect
	github.com/google/go-containerregistry v0.20.6 // indirect
	github.com/google/go-dap v0.12.0 // indirect
	github.com/google/licensecheck v0.3.1 // indirect
	github.com/google/pprof v0.0.0-20250923004556-9e5a51aed1e8 // indirect
	github.com/google/s2a-go v0.1.9 // indirect
	github.com/google/wire v0.6.0 // indirect
	github.com/google/yamlfmt v0.17.2 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.6 // indirect
	github.com/googleapis/gax-go/v2 v2.15.0 // indirect
	github.com/gookit/color v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.4-0.20250319132907-e064f32e3674 // indirect
	github.com/gosuri/uitable v0.0.4 // indirect
	github.com/gotesttools/gotestfmt/v2 v2.5.0 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.2 // indirect
	github.com/hailocab/go-hostpool v0.0.0-20160125115350-e80d13ce29ed // indirect
	github.com/hako/durafmt v0.0.0-20210608085754-5c1018a4e16b // indirect
	github.com/hashicorp/aws-sdk-go-base/v2 v2.0.0-beta.65 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-getter v1.8.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-safetemp v1.0.0 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/hashicorp/hcl/v2 v2.24.0 // indirect
	github.com/haya14busa/goplay v1.0.0 // indirect
	github.com/huandu/xstrings v1.5.0 // indirect
	github.com/iancoleman/strcase v0.3.0 // indirect
	github.com/imdario/mergo v0.3.13 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/itchyny/gojq v0.12.17 // indirect
	github.com/itchyny/timefmt-go v0.1.6 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/jdx/go-netrc v1.0.0 // indirect
	github.com/jinzhu/copier v0.4.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/jmoiron/sqlx v1.4.0 // indirect
	github.com/josharian/impl v1.4.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kastenhq/goversion v0.0.0-20230811215019-93b2f8823953 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/kevinburke/ssh_config v1.2.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.5 // indirect
	github.com/klauspost/pgzip v1.2.6 // indirect
	github.com/knqyf263/go-apk-version v0.0.0-20200609155635-041fdbb8563f // indirect
	github.com/knqyf263/go-deb-version v0.0.0-20190517075300-09fca494f03d // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/liggitt/tabwriter v0.0.0-20181228230101-89fcab3d43de // indirect
	github.com/lmittmann/tint v1.0.7 // indirect
	github.com/masahiro331/go-mvn-version v0.0.0-20210429150710-d3157d602a08 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/mfridman/tparse v0.18.0 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/mholt/archives v0.1.3 // indirect
	github.com/mikelolasagasti/xz v1.0.1 // indirect
	github.com/minio/minlz v1.0.0 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/moby/docker-image-spec v1.3.1 // indirect
	github.com/moby/locker v1.0.1 // indirect
	github.com/moby/spdystream v0.5.0 // indirect
	github.com/moby/sys/mountinfo v0.7.2 // indirect
	github.com/moby/sys/sequential v0.6.0 // indirect
	github.com/moby/sys/signal v0.7.0 // indirect
	github.com/moby/sys/user v0.3.0 // indirect
	github.com/moby/sys/userns v0.1.0 // indirect
	github.com/moby/term v0.5.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.3-0.20250322232337-35a7c28c31ee // indirect
	github.com/momotaro98/strictgoimports v1.2.2 // indirect
	github.com/monochromegane/go-gitignore v0.0.0-20200626010858-205db1a8cc00 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/muesli/termenv v0.16.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/mwitkow/go-proto-validators v0.0.0-20180403085117-0950a7990007 // indirect
	github.com/mxk/go-flowrate v0.0.0-20140419014527-cca7078d478f // indirect
	github.com/ncruces/go-strftime v0.1.9 // indirect
	github.com/nix-community/go-nix v0.0.0-20250101154619-4bdde671e0a1 // indirect
	github.com/nwaples/rardecode v1.1.3 // indirect
	github.com/nwaples/rardecode/v2 v2.1.0 // indirect
	github.com/olekukonko/errors v1.1.0 // indirect
	github.com/olekukonko/ll v0.0.9 // indirect
	github.com/olekukonko/tablewriter v1.1.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.1 // indirect
	github.com/opencontainers/runtime-spec v1.2.0 // indirect
	github.com/opencontainers/selinux v1.11.1 // indirect
	github.com/openvex/go-vex v0.2.5 // indirect
	github.com/owenrumney/go-sarif v1.1.2-0.20231003122901-1000f5e05554 // indirect
	github.com/package-url/packageurl-go v0.1.1 // indirect
	github.com/pandatix/go-cvss v0.6.2 // indirect
	github.com/pborman/indent v1.2.1 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/pjbgf/sha1cd v0.3.2 // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pkg/profile v1.7.0 // indirect
	github.com/pkg/xattr v0.4.9 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/prometheus/client_golang v1.22.0 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.66.1 // indirect
	github.com/prometheus/procfs v0.16.1 // indirect
	github.com/pseudomuto/protoc-gen-doc v1.5.1 // indirect
	github.com/pseudomuto/protokit v0.2.0 // indirect
	github.com/quasilyte/gogrep v0.5.0 // indirect
	github.com/quasilyte/stdinfo v0.0.0-20220114132959-f7386bf02567 // indirect
	github.com/quic-go/qpack v0.5.1 // indirect
	github.com/rakyll/hey v0.1.4 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/rs/cors v1.11.1 // indirect
	github.com/rubenv/sql-migrate v1.8.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/rust-secure-code/go-rustaudit v0.0.0-20250226111315-e20ec32e963c // indirect
	github.com/sabhiram/go-gitignore v0.0.0-20210923224102-525f6e181f06 // indirect
	github.com/sagikazarmark/locafero v0.9.0 // indirect
	github.com/sahilm/fuzzy v0.1.1 // indirect
	github.com/saintfish/chardet v0.0.0-20230101081208-5e3ef4b5456d // indirect
	github.com/santhosh-tekuri/jsonschema/v6 v6.0.2 // indirect
	github.com/sassoftware/go-rpmutils v0.4.0 // indirect
	github.com/scylladb/go-reflectx v1.0.1 // indirect
	github.com/scylladb/go-set v1.0.3-0.20200225121959-cc7b2070d91e // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	github.com/segmentio/encoding v0.5.3 // indirect
	github.com/segmentio/golines v0.13.0 // indirect
	github.com/sergi/go-diff v1.4.0 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/sirupsen/logrus v1.9.4-0.20230606125235-dd1b4c2e81af // indirect
	github.com/skeema/knownhosts v1.3.1 // indirect
	github.com/skratchdot/open-golang v0.0.0-20200116055534-eef842397966 // indirect
	github.com/sorairolake/lzip-go v0.3.5 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spdx/gordf v0.0.0-20201111095634-7098f93598fb // indirect
	github.com/spdx/tools-golang v0.5.5 // indirect
	github.com/spf13/afero v1.15.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/spf13/cobra v1.10.1 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/spf13/viper v1.20.1 // indirect
	github.com/spiffe/go-spiffe/v2 v2.5.0 // indirect
	github.com/stern/stern v0.0.0-00010101000000-000000000000 // indirect
	github.com/stoewer/go-strcase v1.3.1 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/sylabs/sif/v2 v2.22.0 // indirect
	github.com/sylabs/squashfs v1.0.6 // indirect
	github.com/tetratelabs/wazero v1.9.0 // indirect
	github.com/therootcompany/xz v1.0.1 // indirect
	github.com/ulikunitz/xz v0.5.15 // indirect
	github.com/vbatts/go-mtree v0.6.0 // indirect
	github.com/vbatts/tar-split v0.12.1 // indirect
	github.com/vifraa/gopom v1.0.0 // indirect
	github.com/wagoodman/go-partybus v0.0.0-20230516145632-8ccac152c651 // indirect
	github.com/wagoodman/go-presenter v0.0.0-20211015174752-f9c01afc824b // indirect
	github.com/wagoodman/go-progress v0.0.0-20230925121702-07e42b3cdba0 // indirect
	github.com/x-cray/logrus-prefixed-formatter v0.5.2 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	github.com/xanzy/ssh-agent v0.3.3 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	github.com/xhit/go-str2duration/v2 v2.1.0 // indirect
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8 // indirect
	github.com/xlab/treeprint v1.2.0 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	github.com/zclconf/go-cty v1.16.3 // indirect
	github.com/zeebo/errs v1.4.0 // indirect
	go.lsp.dev/jsonrpc2 v0.10.0 // indirect
	go.lsp.dev/pkg v0.0.0-20210717090340-384b27a52fb2 // indirect
	go.lsp.dev/protocol v0.12.0 // indirect
	go.lsp.dev/uri v0.3.0 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/contrib/detectors/gcp v1.37.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.63.0 // indirect
	go.opentelemetry.io/proto/otlp v1.8.0 // indirect
	go.starlark.net v0.0.0-20231101134539-556fd59b42f6 // indirect
	go.uber.org/mock v0.6.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.yaml.in/yaml/v2 v2.4.2 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	go4.org v0.0.0-20230225012048-214862532bf5 // indirect
	golang.org/x/arch v0.11.0 // indirect
	golang.org/x/crypto v0.42.0 // indirect
	golang.org/x/exp v0.0.0-20250819193227-8b4c13bb791b // indirect
	golang.org/x/exp/typeparams v0.0.0-20240213143201-ec583247a57a // indirect
	golang.org/x/image v0.31.0 // indirect
	golang.org/x/mod v0.28.0 // indirect
	golang.org/x/telemetry v0.0.0-20250908211612-aef8a434d053 // indirect
	golang.org/x/term v0.35.0 // indirect
	golang.org/x/xerrors v0.0.0-20240903120638-7835f813f4da // indirect
	gomodules.xyz/jsonpatch/v2 v2.4.0 // indirect
	google.golang.org/api v0.247.0 // indirect
	google.golang.org/genproto v0.0.0-20250715232539-7130f93afb79 // indirect
	gopkg.in/evanphx/json-patch.v4 v4.12.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gorm.io/gorm v1.30.0 // indirect
	helm.sh/helm/v3 v3.19.0 // indirect
	honnef.co/go/tools v0.1.3 // indirect
	k8s.io/apiextensions-apiserver v0.34.1 // indirect
	k8s.io/apiserver v0.34.1 // indirect
	k8s.io/component-base v0.34.1 // indirect
	k8s.io/component-helpers v0.34.1 // indirect
	k8s.io/klog/v2 v2.130.1 // indirect
	k8s.io/kube-openapi v0.0.0-20250710124328-f3f2b991d03b // indirect
	k8s.io/kubectl v0.34.1 // indirect
	modernc.org/libc v1.66.3 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.11.0 // indirect
	modernc.org/sqlite v1.39.0 // indirect
	mvdan.cc/gofumpt v0.9.1 // indirect
	oras.land/oras-go/v2 v2.6.0 // indirect
	pluginrpc.com/pluginrpc v0.5.0 // indirect
	sigs.k8s.io/json v0.0.0-20241014173422-cfa47c3a1cc8 // indirect
	sigs.k8s.io/kustomize/api v0.20.1 // indirect
	sigs.k8s.io/kustomize/kyaml v0.20.1 // indirect
	sigs.k8s.io/randfill v1.0.0 // indirect
	sigs.k8s.io/structured-merge-diff/v6 v6.3.0 // indirect
)
