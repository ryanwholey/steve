module github.com/TakeScoop/steve

go 1.15

require (
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
	helm.sh/helm/v3 v3.5.3
	k8s.io/client-go v0.20.2
)

replace (
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d
	github.com/docker/docker => github.com/moby/moby v0.7.3-0.20190826074503-38ab9da00309
)
