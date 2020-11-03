module github.com/giantswarm/app/v3

go 1.15

require (
	github.com/giantswarm/apiextensions/v3 v3.4.2-0.20201103105314-dd1460e43e4f
	github.com/giantswarm/microerror v0.2.1
	github.com/google/go-cmp v0.5.2 // indirect
	github.com/stretchr/testify v1.6.1 // indirect
	k8s.io/apimachinery v0.18.9
	sigs.k8s.io/yaml v1.2.0
)

replace sigs.k8s.io/cluster-api => github.com/giantswarm/cluster-api v0.3.10-gs
