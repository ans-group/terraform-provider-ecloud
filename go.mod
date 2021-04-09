module github.com/ukfast/terraform-provider-ecloud

go 1.13

require (
	github.com/hashicorp/go-hclog v0.12.1 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/terraform v0.14.3
	github.com/stretchr/testify v1.6.1
	github.com/ukfast/sdk-go v1.3.48
)

// replace github.com/ukfast/sdk-go => ../sdk-go
