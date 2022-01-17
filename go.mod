module github.com/red-hat-storage/hpessa-exporter

go 1.16

require (
	github.com/go-logr/logr v0.4.0
	github.com/prometheus/client_golang v1.11.0
	github.com/spf13/cobra v1.2.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/sys v0.0.0-20211124211545-fe61309f8881
	k8s.io/api v0.22.4
	k8s.io/apimachinery v0.22.4
	k8s.io/client-go v0.22.4
	sigs.k8s.io/controller-runtime v0.10.3
)
