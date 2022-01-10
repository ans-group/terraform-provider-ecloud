module github.com/ukfast/terraform-provider-ecloud

go 1.13

require (
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.7.1
	github.com/stretchr/testify v1.7.0
	github.com/ukfast/sdk-go v1.4.29
)

// replace github.com/ukfast/sdk-go => ../sdk-go
