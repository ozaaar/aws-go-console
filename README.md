[![Go Report Card](https://goreportcard.com/badge/github.com/ozaaar/aws-go-console)](https://goreportcard.com/report/github.com/ozaaar/aws-go-console)
[![Go Reference](https://pkg.go.dev/badge/github.com/ozaaar/aws-go-console.svg)](https://pkg.go.dev/github.com/ozaaar/aws-go-console)

# aws-go-console
A helper package provides short-lived (scoped based) token/url for AWS console. It is based on the [documentation](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_providers_enable-console-custom-url.html) provided by AWS.

Package `console` have following use cases:
- Give federated access to a user via AWS Management Console without an IAM User.
- Allow users who sign in to your organization's network securely access the AWS Management Console.
