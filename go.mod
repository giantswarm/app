module github.com/giantswarm/app/v6

go 1.19

require (
	github.com/giantswarm/apiextensions-application v0.6.0
	github.com/giantswarm/k8smetadata v0.13.0
	github.com/giantswarm/microerror v0.4.0
	github.com/giantswarm/micrologger v0.6.0
	github.com/giantswarm/to v0.4.0
	github.com/google/go-cmp v0.5.9
	github.com/google/go-github/v39 v39.2.0
	github.com/imdario/mergo v0.3.13
	golang.org/x/oauth2 v0.3.0
	k8s.io/api v0.20.15
	k8s.io/apiextensions-apiserver v0.20.15
	k8s.io/apimachinery v0.20.15
	k8s.io/client-go v0.20.15
	sigs.k8s.io/controller-runtime v0.6.5
	sigs.k8s.io/yaml v1.3.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/evanphx/json-patch v4.11.0+incompatible // indirect
	github.com/go-kit/log v0.2.0 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/go-logr/logr v0.4.0 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/googleapis/gnostic v0.5.5 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/crypto v0.0.0-20220525230936-793ad666bf5e // indirect
	golang.org/x/net v0.3.0 // indirect
	golang.org/x/sys v0.3.0 // indirect
	golang.org/x/term v0.3.0 // indirect
	golang.org/x/text v0.5.0 // indirect
	golang.org/x/time v0.0.0-20220609170525-579cf78fd858 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/klog/v2 v2.9.0 // indirect
	k8s.io/kube-openapi v0.0.0-20211110013926-83f114cd0513 // indirect
	k8s.io/utils v0.0.0-20210802155522-efc7438f0176 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.1.2 // indirect
)

replace (
	github.com/Microsoft/hcsshim v0.8.7 => github.com/Microsoft/hcsshim v0.8.10
	github.com/bketelsen/crypt => github.com/bketelsen/crypt v0.0.5
	github.com/dgrijalva/jwt-go => github.com/dgrijalva/jwt-go/v4 v4.0.0-preview1
	github.com/docker/docker => github.com/moby/moby v20.10.22+incompatible
	github.com/gogo/protobuf v1.3.1 => github.com/gogo/protobuf v1.3.2
	github.com/gorilla/websocket => github.com/gorilla/websocket v1.5.0
	github.com/miekg/dns v1.0.14 => github.com/miekg/dns v1.1.50
	github.com/miekg/dns v1.1.41 => github.com/miekg/dns v1.1.50
	github.com/opencontainers/runc v0.1.1 => github.com/opencontainers/runc v1.0.0-rc7
	github.com/prometheus/client_golang v1.11.0 => github.com/prometheus/client_golang v1.12.2
	github.com/prometheus/client_golang v1.7.1 => github.com/prometheus/client_golang v1.12.2
	github.com/spf13/viper => github.com/spf13/viper v1.14.0
	go.mongodb.org/mongo-driver v1.1.2 => go.mongodb.org/mongo-driver v1.9.1
	golang.org/x/net v0.3.0 => golang.org/x/net v0.4.0
	// Use fork of CAPI with Kubernetes 1.18 support.
	sigs.k8s.io/cluster-api => github.com/giantswarm/cluster-api v0.3.10-gs
)
