# aws-go-console
A helper package provides short-lived (scoped based) token/url for AWS console. It is based on the [documentation](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_providers_enable-console-custom-url.html) provided by AWS.

Package `console` have following use cases:
- Give federated access to a user via AWS Management Console without an IAM User.
- Allow users who sign in to your organization's network securely access the AWS Management Console.
