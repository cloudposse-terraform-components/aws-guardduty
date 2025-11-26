package test

import (
	"strings"
	"testing"

	"github.com/cloudposse/test-helpers/pkg/atmos"
	helper "github.com/cloudposse/test-helpers/pkg/atmos/component-helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ComponentSuite defines the test suite for the GuardDuty component
type ComponentSuite struct {
	helper.TestSuite
}

const (
	defaultStack  = "default-test"
	defaultRegion = "us-east-2"
)

// TestRunGuardDutySuite is the entry point for running all GuardDuty component tests
func TestRunGuardDutySuite(t *testing.T) {
	suite := new(ComponentSuite)
	helper.Run(t, suite)
}

// TestGuardDutyDelegatedAdministrator tests basic GuardDuty detector creation
// in the delegated administrator account (Step 1)
func (s *ComponentSuite) TestGuardDutyDelegatedAdministrator() {
	const component = "guardduty/delegated-administrator"

	defer s.DestroyAtmosComponent(s.T(), component, defaultStack, nil)
	options, _ := s.DeployAtmosComponent(s.T(), component, defaultStack, nil)

	// Verify detector was created
	detectorArn := atmos.Output(s.T(), options, "guardduty_detector_arn")
	assert.NotEmpty(s.T(), detectorArn, "GuardDuty detector ARN should not be empty")
	require.True(s.T(), strings.HasPrefix(detectorArn, "arn:aws:guardduty:"), "Detector ARN should have correct prefix")

	detectorId := atmos.Output(s.T(), options, "guardduty_detector_id")
	assert.NotEmpty(s.T(), detectorId, "GuardDuty detector ID should not be empty")

	// Verify delegated administrator account ID output
	delegatedAdminAccountId := atmos.Output(s.T(), options, "delegated_administrator_account_id")
	assert.NotEmpty(s.T(), delegatedAdminAccountId, "Delegated administrator account ID should not be empty")

	// Verify SNS resources are not created (create_sns_topic is false)
	snsTopicArn := atmos.Output(s.T(), options, "root_sns_topic_arn")
	assert.Empty(s.T(), snsTopicArn, "SNS topic ARN should be empty when create_sns_topic is false")
}

// TestGuardDutyWithSNSNotifications tests GuardDuty with SNS notifications enabled
func (s *ComponentSuite) TestGuardDutyWithSNSNotifications() {
	const component = "guardduty/with-sns"

	defer s.DestroyAtmosComponent(s.T(), component, defaultStack, nil)
	options, _ := s.DeployAtmosComponent(s.T(), component, defaultStack, nil)

	// Verify detector was created
	detectorArn := atmos.Output(s.T(), options, "guardduty_detector_arn")
	assert.NotEmpty(s.T(), detectorArn, "GuardDuty detector ARN should not be empty")

	// Verify SNS topic was created
	snsTopicArn := atmos.Output(s.T(), options, "root_sns_topic_arn")
	assert.NotEmpty(s.T(), snsTopicArn, "SNS topic ARN should not be empty")
	require.True(s.T(), strings.HasPrefix(snsTopicArn, "arn:aws:sns:"), "SNS topic ARN should have correct prefix")

	snsTopicName := atmos.Output(s.T(), options, "sns_topic_name")
	assert.NotEmpty(s.T(), snsTopicName, "SNS topic name should not be empty")

	// Verify SQS queue was created
	sqsQueueArn := atmos.Output(s.T(), options, "root_sqs_queue_arn")
	assert.NotEmpty(s.T(), sqsQueueArn, "SQS queue ARN should not be empty")
	require.True(s.T(), strings.HasPrefix(sqsQueueArn, "arn:aws:sqs:"), "SQS queue ARN should have correct prefix")

	sqsQueueName := atmos.Output(s.T(), options, "root_sqs_queue_name")
	assert.NotEmpty(s.T(), sqsQueueName, "SQS queue name should not be empty")

	// Verify KMS key was created
	kmsKeyArn := atmos.Output(s.T(), options, "root_kms_key_arn")
	assert.NotEmpty(s.T(), kmsKeyArn, "KMS key ARN should not be empty")
	require.True(s.T(), strings.HasPrefix(kmsKeyArn, "arn:aws:kms:"), "KMS key ARN should have correct prefix")

	kmsKeyId := atmos.Output(s.T(), options, "root_kms_key_id")
	assert.NotEmpty(s.T(), kmsKeyId, "KMS key ID should not be empty")

	// Verify CloudWatch event rule was created
	cloudwatchRuleArn := atmos.Output(s.T(), options, "cloudwatch_event_rule_arn")
	assert.NotEmpty(s.T(), cloudwatchRuleArn, "CloudWatch event rule ARN should not be empty")
	require.True(s.T(), strings.HasPrefix(cloudwatchRuleArn, "arn:aws:events:"), "CloudWatch event rule ARN should have correct prefix")
}

// TestGuardDutyWithProtectionFeatures tests GuardDuty with various protection features enabled
func (s *ComponentSuite) TestGuardDutyWithProtectionFeatures() {
	const component = "guardduty/with-features"

	defer s.DestroyAtmosComponent(s.T(), component, defaultStack, nil)
	options, _ := s.DeployAtmosComponent(s.T(), component, defaultStack, nil)

	// Verify detector was created
	detectorArn := atmos.Output(s.T(), options, "guardduty_detector_arn")
	assert.NotEmpty(s.T(), detectorArn, "GuardDuty detector ARN should not be empty")
	require.True(s.T(), strings.HasPrefix(detectorArn, "arn:aws:guardduty:"), "Detector ARN should have correct prefix")

	detectorId := atmos.Output(s.T(), options, "guardduty_detector_id")
	assert.NotEmpty(s.T(), detectorId, "GuardDuty detector ID should not be empty")
}

// TestGuardDutyDisabled tests that the component creates nothing when disabled
func (s *ComponentSuite) TestGuardDutyDisabled() {
	const component = "guardduty/disabled"

	defer s.DestroyAtmosComponent(s.T(), component, defaultStack, nil)
	options, _ := s.DeployAtmosComponent(s.T(), component, defaultStack, nil)

	// Validate outputs - all should be empty when disabled
	detectorArn := atmos.Output(s.T(), options, "guardduty_detector_arn")
	assert.Empty(s.T(), detectorArn, "GuardDuty detector ARN should be empty when disabled")

	snsTopicArn := atmos.Output(s.T(), options, "root_sns_topic_arn")
	assert.Empty(s.T(), snsTopicArn, "SNS topic ARN should be empty when disabled")

	delegatedAdminAccountId := atmos.Output(s.T(), options, "delegated_administrator_account_id")
	assert.Empty(s.T(), delegatedAdminAccountId, "Delegated administrator account ID should be empty when disabled")
}
