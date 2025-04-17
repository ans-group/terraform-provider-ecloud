FROM golang:1.24 AS builder
COPY . /build
WORKDIR /build
RUN go mod download
RUN CGO_ENABLED=0 go build

FROM hashicorp/terraform:latest
COPY --from=builder /build/terraform-provider-ecloud /root/.terraform.d/plugins/linux_amd64/terraform-provider-ecloud
