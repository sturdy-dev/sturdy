module getsturdy.com/api

go 1.17

// Fork of go-diff
// Patched to introduce Windows compatible APIs
// Patched to add API to always quote filenames
replace github.com/sourcegraph/go-diff => github.com/zegl/go-diff v0.6.2-0.20211129194800-6d29830744dc

// replace github.com/sourcegraph/go-diff => /Users/gustav/src/go-diff

// Custom fork of graphql-transport-ws
// https://github.com/sturdy-dev/graphql-transport-ws/commits/sturdy
replace github.com/graph-gophers/graphql-transport-ws => github.com/sturdy-dev/graphql-transport-ws v0.0.0-20211122094650-15c742155db6

// Custom fork of go-flags to avoid tags name conflict
replace github.com/jessevdk/go-flags => github.com/sturdy-dev/go-flags v1.5.1-0.20220203104421-967e8bff1baf

require (
	github.com/aws/aws-sdk-go v1.38.47
	github.com/bmatcuk/doublestar/v4 v4.0.2
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869
	github.com/bradleyfalzon/ghinstallation v1.1.1
	github.com/buildkite/go-buildkite/v3 v3.0.0
	github.com/cenkalti/backoff/v4 v4.1.2
	github.com/disintegration/imaging v1.6.2
	github.com/fatih/color v1.13.0
	github.com/getsentry/raven-go v0.2.0
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-contrib/gzip v0.0.5
	github.com/gin-contrib/zap v0.0.2
	github.com/gin-gonic/gin v1.7.7
	github.com/go-playground/validator/v10 v10.10.0
	github.com/gofrs/flock v0.8.1
	github.com/golang-migrate/migrate/v4 v4.15.1
	github.com/golang/mock v1.6.0
	github.com/google/go-github/v39 v39.2.0
	github.com/google/uuid v1.3.0
	github.com/gosimple/slug v1.9.0
	github.com/graph-gophers/dataloader/v6 v6.0.0
	github.com/graph-gophers/graphql-go v1.3.0
	github.com/graph-gophers/graphql-transport-ws v0.0.1
	github.com/hashicorp/golang-lru v0.5.4
	github.com/jessevdk/go-flags v1.5.0
	github.com/jmoiron/sqlx v1.3.4
	github.com/jxskiss/base62 v0.0.0-20191017122030-4f11678b909b
	github.com/lib/pq v1.10.4
	github.com/libgit2/git2go/v33 v33.0.0
	github.com/mergestat/timediff v0.0.2
	github.com/microcosm-cc/bluemonday v1.0.16
	github.com/posthog/posthog-go v0.0.0-20211028072449-93c17c49e2b0
	github.com/prometheus/client_golang v1.11.0
	github.com/psanford/memfs v0.0.0-20210214183328-a001468d78ef
	github.com/sourcegraph/go-diff v0.6.2-0.20210526090523-35b24a7eb480
	github.com/stretchr/testify v1.7.0
	github.com/tailscale/hujson v0.0.0-20210818175511-7360507a6e88
	github.com/tidwall/match v1.0.3
	github.com/yuin/goldmark v1.4.4
	go.uber.org/dig v1.13.0
	go.uber.org/multierr v1.7.0
	go.uber.org/zap v1.19.1
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
	golang.org/x/oauth2 v0.0.0-20210628180205-a41e5a781914
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211
	gopkg.in/square/go-jose.v2 v2.6.0
)

require (
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/benbjohnson/clock v1.3.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cenkalti/backoff v2.0.0+incompatible // indirect
	github.com/certifi/gocertifi v0.0.0-20210507211836-431795d63e8d // indirect
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-github/v29 v29.0.2 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/gorilla/css v1.0.0 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mattn/go-colorable v0.1.9 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.26.0 // indirect
	github.com/prometheus/procfs v0.6.0 // indirect
	github.com/rainycape/unidecode v0.0.0-20150907023854-cb7f23ec59be // indirect
	github.com/rogpeppe/go-internal v1.8.0 // indirect
	github.com/ugorji/go/codec v1.1.7 // indirect
	github.com/xtgo/uuid v0.0.0-20140804021211-a0b114877d4c // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/goleak v1.1.11 // indirect
	golang.org/x/image v0.0.0-20210216034530-4410531fe030 // indirect
	golang.org/x/mod v0.4.2 // indirect
	golang.org/x/net v0.0.0-20211013171255-e13a2654a71e // indirect
	golang.org/x/sys v0.0.0-20211025201205-69cdffdb9359 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/tools v0.1.5 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
