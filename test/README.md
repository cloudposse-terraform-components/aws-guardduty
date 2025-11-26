# Tests

This directory contains integration tests for the `aws-guardduty` component
using [Terratest](https://terratest.gruntwork.io/) and [Atmos](https://atmos.tools/).

## Prerequisites

- Go 1.25+
- Terraform 1.3+
- [Atmos](https://atmos.tools/install) CLI
- AWS credentials configured with appropriate permissions
- `TEST_ACCOUNT_ID` environment variable set to your AWS account ID

## Test Structure

```
test/
├── component_test.go          # Go test file with all test cases
├── go.mod                     # Go module definition
├── go.sum                     # Go module checksums (generated)
├── vendor.yaml                # Atmos vendor configuration
├── fixtures/
│   ├── atmos.yaml            # Atmos CLI configuration
│   ├── components/terraform/  # Vendored components (generated)
│   └── stacks/
│       ├── catalog/
│       │   ├── account-map.yaml
│       │   └── usecase/
│       │       ├── delegated-administrator.yaml
│       │       ├── with-sns.yaml
│       │       ├── with-features.yaml
│       │       └── disabled.yaml
│       └── orgs/default/test/
│           ├── _defaults.yaml
│           └── test.yaml
```

## Test Cases

| Test Name                             | Description                                                        |
|---------------------------------------|--------------------------------------------------------------------|
| `TestGuardDutyDelegatedAdministrator` | Tests basic GuardDuty detector creation in delegated admin account |
| `TestGuardDutyWithSNSNotifications`   | Tests GuardDuty with SNS topic, SQS queue, and KMS key creation    |
| `TestGuardDutyWithProtectionFeatures` | Tests GuardDuty with various protection features enabled           |
| `TestGuardDutyDisabled`               | Tests that nothing is created when the component is disabled       |

## Running Tests

### Setup

1. Set the required environment variable:

```bash
export TEST_ACCOUNT_ID="123456789012"  # Your AWS account ID
```

2. Install Go dependencies:

```bash
cd test
go mod download
```

3. Vendor required components:

```bash
cd fixtures
atmos vendor pull
cd ..
```

### Run All Tests

```bash
go test -v -timeout 60m ./...
```

### Run Specific Test

```bash
go test -v -timeout 30m -run TestGuardDutyDelegatedAdministrator ./...
```

> **Note:** Tests run sequentially (not in parallel) because GuardDuty organization-level resources are singletons.
> The `aws_guardduty_organization_admin_account` resource can only exist once per organization, so parallel
> test execution would cause race conditions.

## Test Configuration

Each test case uses a different component configuration defined in `fixtures/stacks/catalog/usecase/`:

- **delegated-administrator.yaml**: Basic GuardDuty detector with S3 protection only
- **with-sns.yaml**: GuardDuty with SNS notifications and CloudWatch events
- **with-features.yaml**: GuardDuty with all protection features enabled (S3, EKS, Lambda, Runtime Monitoring)
- **disabled.yaml**: Component disabled, should create no resources

## Cleanup

Tests are configured with `DestroyOnCompletion: true`, so resources are automatically cleaned up after each test. If a
test fails mid-execution, you may need to manually clean up resources:

```bash
cd fixtures
atmos terraform destroy guardduty/<component-name> -s default-test
```

## Troubleshooting

### Test Timeout

If tests timeout, increase the timeout value:

```bash
go test -v -timeout 120m ./...
```

### AWS Permissions

Ensure your AWS credentials have permissions for:

- GuardDuty (create/delete detectors, features)
- SNS (create/delete topics, subscriptions)
- SQS (create/delete queues, policies)
- KMS (create/delete keys, policies)
- CloudWatch Events (create/delete rules, targets)
- IAM (for service-linked roles)

### State Conflicts

If you encounter state conflicts, remove the local state files:

```bash
rm -rf fixtures/state/
```
