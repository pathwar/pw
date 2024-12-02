module pathwar.land/pathwar/v2

require (
	cloud.google.com/go v0.37.4 // indirect
	github.com/Bearer/bearer-go v1.2.1
	github.com/Microsoft/go-winio v0.4.14 // indirect
	github.com/adrg/xdg v0.3.3
	github.com/alessio/shellescape v1.2.2
	github.com/brianvoe/gofakeit v3.18.0+incompatible
	github.com/bwmarrin/snowflake v0.3.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v1.13.1
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-units v0.4.0 // indirect
	github.com/dustin/go-humanize v1.0.0
	github.com/fatih/color v1.9.0 // indirect
	github.com/getsentry/sentry-go v0.6.1
	github.com/githubnemo/CompileDaemon v1.2.1
	github.com/go-chi/chi v4.1.1+incompatible
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gobuffalo/packr/v2 v2.8.0
	github.com/gogo/gateway v1.1.0
	github.com/gogo/protobuf v1.3.2
	github.com/golang/protobuf v1.5.0
	github.com/google/go-cmp v0.5.5 // indirect
	github.com/google/go-querystring v1.0.0
	github.com/gosimple/slug v1.9.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0
	github.com/grpc-ecosystem/grpc-gateway v1.14.5
	github.com/jinzhu/gorm v1.9.12
	github.com/karrick/godirwalk v1.15.6 // indirect
	github.com/martinlindhe/base36 v1.0.0
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/moby/moby v20.10.20+incompatible
	github.com/oklog/run v1.1.1-0.20200508094559-c7096881717e
	github.com/olekukonko/tablewriter v0.0.4
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opentracing/opentracing-go v1.1.0
	github.com/openzipkin-contrib/zipkin-go-opentracing v0.4.5
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/peterbourgon/ff/v3 v3.0.0
	github.com/pkg/errors v0.9.1
	github.com/rogpeppe/go-internal v1.6.0 // indirect
	github.com/rs/cors v1.7.0
	github.com/sirupsen/logrus v1.6.0 // indirect
	github.com/soheilhy/cmux v0.1.4
	github.com/stretchr/testify v1.7.0
	github.com/tailscale/depaware v0.0.0-20201214215404-77d1e9757027
	github.com/treastech/logger v0.0.0-20180705232552-e381e9ecf2e3
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0
	go.uber.org/zap v1.21.0
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	golang.org/x/lint v0.0.0-20201208152925-83fdc39ff7b5 // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
	google.golang.org/appengine v1.6.6 // indirect
	google.golang.org/genproto v0.0.0-20190927181202-20e1ac93f88c
	google.golang.org/grpc v1.28.0-pre
	gopkg.in/gormigrate.v1 v1.6.0
	gopkg.in/yaml.v2 v2.3.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	moul.io/banner v1.0.1
	moul.io/godev v1.6.0
	moul.io/motd v1.0.0
	moul.io/roundtripper v1.0.0
	moul.io/srand v1.4.0
	moul.io/u v1.23.0
	moul.io/zapconfig v1.4.0
	moul.io/zapgorm v1.0.0
)

replace (
	//github.com/Bearer/bearer-go => ../github.com/Bearer/bearer-go
	github.com/golang/lint => golang.org/x/lint v0.0.0-20190409202823-959b441ac422
	github.com/golang/protobuf => github.com/golang/protobuf v1.3.2
	gopkg.in/jcmturner/rpc.v1 => gopkg.in/jcmturner/rpc.v1 v1.1.0
	//moul.io/godev => ../moul.io/godev
	sourcegraph.com/sourcegraph/go-diff => github.com/sourcegraph/go-diff v0.5.1
)

go 1.13
