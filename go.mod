module github.com/giantswarm/app/v7

go 1.21

toolchain go1.22.4

require (
	github.com/giantswarm/apiextensions-application v0.6.2
	github.com/giantswarm/k8smetadata v0.25.0
	github.com/giantswarm/microerror v0.4.1
	github.com/giantswarm/micrologger v1.1.1
	github.com/giantswarm/to v0.4.0
	github.com/google/go-cmp v0.6.0
	github.com/google/go-github/v62 v62.0.0
	github.com/imdario/mergo v0.3.16
	golang.org/x/oauth2 v0.21.0
	k8s.io/api v0.20.15
	k8s.io/apiextensions-apiserver v0.20.15
	k8s.io/apimachinery v0.20.15
	k8s.io/client-go v0.20.15
	sigs.k8s.io/controller-runtime v0.6.5
	sigs.k8s.io/yaml v1.4.0
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/evanphx/json-patch v4.11.0+incompatible // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/googleapis/gnostic v0.5.5 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/crypto v0.24.0 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/term v0.21.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/klog/v2 v2.120.1 // indirect
	k8s.io/kube-openapi v0.0.0-20211110013926-83f114cd0513 // indirect
	k8s.io/utils v0.0.0-20210802155522-efc7438f0176 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.1.2 // indirect
)

replace (
	github.com/dgrijalva/jwt-go => github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/docker/docker => github.com/moby/moby v26.1.4+incompatible
	github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.2
	github.com/golang/protobuf => github.com/golang/protobuf v1.5.4
	github.com/gorilla/websocket => github.com/gorilla/websocket v1.5.1
	github.com/miekg/dns => github.com/miekg/dns v1.1.59
	github.com/prometheus/client_golang => github.com/prometheus/client_golang v1.19.1
	github.com/spf13/viper => github.com/spf13/viper v1.19.0
	go.mongodb.org/mongo-driver => go.mongodb.org/mongo-driver v1.15.0
	golang.org/x/net => golang.org/x/net v0.26.0
	google.golang.org/protobuf v1.32.0 => google.golang.org/protobuf v1.33.0
	// Use fork of CAPI with Kubernetes 1.18 support.
	sigs.k8s.io/cluster-api => github.com/giantswarm/cluster-api v1.5.0
)
