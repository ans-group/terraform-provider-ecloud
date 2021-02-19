# terraform-provider-ecloud

## Getting Started

This provider is available via the [Terraform Registry](https://registry.terraform.io/providers/ukfast/ecloud/latest) with Terraform v0.13+

> :warning: We strongly recommend pinning the provider version to a target major version, as to ensure future breaking changes do not affect workflows and automated CI pipelines

```
terraform {
  required_providers {
    ecloud = {
      source  = "ukfast/ecloud"
      version = "~> 2.0"
    }
  }
}

provider "ecloud" {
  api_key = "abc"
}
```

## Getting Started (manual)

To get started, the `terraform-provider-ecloud` binary (`.exe` extension if Windows) should be downloaded from [Releases](https://github.com/ukfast/terraform-provider-ecloud/releases) and placed in the plugins directory (see [here](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins) for more information). For this example, we'll place it in `~/.terraform.d/plugins/`

We'll then need to initialise terraform with our provider:

```console
terraform init
```

Finally, we can invoke `terraform apply` to apply our terraform configuration:

```console
terraform apply
```

## Documentation

Documentation is located within this repository at `/docs`, and is published in the [Terraform Registry](https://registry.terraform.io/providers/ukfast/ecloud/latest/docs)

## Development

### Testing

Acceptance tests can be executed using `make` as below:

```
export UKF_TEST_VPC_REGION_ID=<region_id>
make testacc TEST=VPC_basic
```

Note `UKF_TEST_VPC_REGION_ID` environment is required for executing tests


### Releasing 

`goreleaser` is used to release the provider on Github. First, obtain your GPG fingerprint:

```
gpg -k
```

Cache GPG passphrase:

```
gpg --armor --detach-sign .
```

Finally tag and invoke `goreleaser`:

```
git tag v2.0.0
git push --tags
export GITHUB_TOKEN=<token>
export GPG_FINGERPRINT=<fingerprint>
goreleaser --rm-dist
```
