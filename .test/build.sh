go build -o terraform-provider-ecloud
export OS_ARCH="$(go env GOHOSTOS)_$(go env GOHOSTARCH)"
mkdir -p ~/.terraform.d/plugins/ukfast.io/test/ecloud/0.1/$OS_ARCH
mv terraform-provider-ecloud ~/.terraform.d/plugins/ukfast.io/test/ecloud/0.1/$OS_ARCH