package test

import (
	"testing"

	"github.com/cloudposse/test-helpers/pkg/atmos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGuardDutyDelegatedAdministrator tests the GuardDuty detector creation
// in the delegated administrator account (Step 1)
func TestGuardDutyDelegatedAdministrator(t *testing.T) {
	t.Parallel()

	fixture := atmos.Fixture{
		TestFolder:          "fixtures",
		StackName:           "default-test",
		ComponentName:       "guardduty/delegated-administrator",
		TerraformDir:        "../src",
		RandomSeed:          "guardduty-delegated-admin",
		DestroyOnCompletion: true,
	}

	defer fixture.TearDown(t)
	fixture.SetUp(t)

	fixture.Plan(t)
	fixture.Apply(t)

	// Validate outputs
	outputs := fixture.Outputs(t)

	// Verify detector was created
	detectorArn := outputs["guardduty_detector_arn"]
	assert.NotNil(t, detectorArn, "GuardDuty detector ARN should not be nil")
	assert.NotEmpty(t, detectorArn, "GuardDuty detector ARN should not be empty")

	detectorId := outputs["guardduty_detector_id"]
	assert.NotNil(t, detectorId, "GuardDuty detector ID should not be nil")
	assert.NotEmpty(t, detectorId, "GuardDuty detector ID should not be empty")
}

// TestGuardDutyWithSNSNotifications tests GuardDuty with SNS notifications enabled
func TestGuardDutyWithSNSNotifications(t *testing.T) {
	t.Parallel()

	fixture := atmos.Fixture{
		TestFolder:          "fixtures",
		StackName:           "default-test",
		ComponentName:       "guardduty/with-sns",
		TerraformDir:        "../src",
		RandomSeed:          "guardduty-with-sns",
		DestroyOnCompletion: true,
	}

	defer fixture.TearDown(t)
	fixture.SetUp(t)

	fixture.Plan(t)
	fixture.Apply(t)

	// Validate outputs
	outputs := fixture.Outputs(t)

	// Verify detector was created
	detectorArn := outputs["guardduty_detector_arn"]
	assert.NotNil(t, detectorArn, "GuardDuty detector ARN should not be nil")

	// Verify SNS topic was created
	snsTopicArn := outputs["root_sns_topic_arn"]
	assert.NotNil(t, snsTopicArn, "SNS topic ARN should not be nil")
	assert.NotEmpty(t, snsTopicArn, "SNS topic ARN should not be empty")

	snsTopicName := outputs["sns_topic_name"]
	assert.NotNil(t, snsTopicName, "SNS topic name should not be nil")

	// Verify SQS queue was created
	sqsQueueArn := outputs["root_sqs_queue_arn"]
	assert.NotNil(t, sqsQueueArn, "SQS queue ARN should not be nil")
	assert.NotEmpty(t, sqsQueueArn, "SQS queue ARN should not be empty")

	// Verify KMS key was created
	kmsKeyArn := outputs["root_kms_key_arn"]
	assert.NotNil(t, kmsKeyArn, "KMS key ARN should not be nil")
	assert.NotEmpty(t, kmsKeyArn, "KMS key ARN should not be empty")

	// Verify CloudWatch event rule was created
	cloudwatchRuleArn := outputs["cloudwatch_event_rule_arn"]
	assert.NotNil(t, cloudwatchRuleArn, "CloudWatch event rule ARN should not be nil")
	assert.NotEmpty(t, cloudwatchRuleArn, "CloudWatch event rule ARN should not be empty")
}

// TestGuardDutyWithProtectionFeatures tests GuardDuty with various protection features enabled
func TestGuardDutyWithProtectionFeatures(t *testing.T) {
	t.Parallel()

	fixture := atmos.Fixture{
		TestFolder:          "fixtures",
		StackName:           "default-test",
		ComponentName:       "guardduty/with-features",
		TerraformDir:        "../src",
		RandomSeed:          "guardduty-with-features",
		DestroyOnCompletion: true,
	}

	defer fixture.TearDown(t)
	fixture.SetUp(t)

	fixture.Plan(t)
	fixture.Apply(t)

	// Validate outputs
	outputs := fixture.Outputs(t)

	// Verify detector was created
	detectorArn := outputs["guardduty_detector_arn"]
	assert.NotNil(t, detectorArn, "GuardDuty detector ARN should not be nil")
	assert.NotEmpty(t, detectorArn, "GuardDuty detector ARN should not be empty")

	detectorId := outputs["guardduty_detector_id"]
	assert.NotNil(t, detectorId, "GuardDuty detector ID should not be nil")
	assert.NotEmpty(t, detectorId, "GuardDuty detector ID should not be empty")
}

// TestGuardDutyDisabled tests that the component creates nothing when disabled
func TestGuardDutyDisabled(t *testing.T) {
	t.Parallel()

	fixture := atmos.Fixture{
		TestFolder:          "fixtures",
		StackName:           "default-test",
		ComponentName:       "guardduty/disabled",
		TerraformDir:        "../src",
		RandomSeed:          "guardduty-disabled",
		DestroyOnCompletion: true,
	}

	defer fixture.TearDown(t)
	fixture.SetUp(t)

	fixture.Plan(t)
	fixture.Apply(t)

	// Validate outputs - all should be null when disabled
	outputs := fixture.Outputs(t)

	detectorArn := outputs["guardduty_detector_arn"]
	assert.Nil(t, detectorArn, "GuardDuty detector ARN should be nil when disabled")

	snsTopicArn := outputs["root_sns_topic_arn"]
	assert.Nil(t, snsTopicArn, "SNS topic ARN should be nil when disabled")
}

// TestGuardDutyValidationRuntimeMonitoringConflict tests that enabling both
// runtime_monitoring_enabled and eks_runtime_monitoring_enabled fails validation
func TestGuardDutyValidationRuntimeMonitoringConflict(t *testing.T) {
	t.Parallel()

	fixture := atmos.Fixture{
		TestFolder:          "fixtures",
		StackName:           "default-test",
		ComponentName:       "guardduty/validation-conflict",
		TerraformDir:        "../src",
		RandomSeed:          "guardduty-validation-conflict",
		DestroyOnCompletion: true,
	}

	defer fixture.TearDown(t)
	fixture.SetUp(t)

	// This should fail during plan due to the precondition
	_, err := fixture.PlanE(t)
	require.Error(t, err, "Plan should fail when both runtime monitoring options are enabled")
}
