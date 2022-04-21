// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

// Code generated by ack-generate. DO NOT EDIT.

package app

import (
	"context"
	"errors"
	"reflect"
	"strings"

	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackcondition "github.com/aws-controllers-k8s/runtime/pkg/condition"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	"github.com/aws/aws-sdk-go/aws"
	svcsdk "github.com/aws/aws-sdk-go/service/sagemaker"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	svcapitypes "github.com/aws-controllers-k8s/sagemaker-controller/apis/v1alpha1"
)

// Hack to avoid import errors during build...
var (
	_ = &metav1.Time{}
	_ = strings.ToLower("")
	_ = &aws.JSONValue{}
	_ = &svcsdk.SageMaker{}
	_ = &svcapitypes.App{}
	_ = ackv1alpha1.AWSAccountID("")
	_ = &ackerr.NotFound
	_ = &ackcondition.NotManagedMessage
	_ = &reflect.Value{}
)

// sdkFind returns SDK-specific information about a supplied resource
func (rm *resourceManager) sdkFind(
	ctx context.Context,
	r *resource,
) (latest *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkFind")
	defer exit(err)
	// If any required fields in the input shape are missing, AWS resource is
	// not created yet. Return NotFound here to indicate to callers that the
	// resource isn't yet created.
	if rm.requiredFieldsMissingFromReadOneInput(r) {
		return nil, ackerr.NotFound
	}

	input, err := rm.newDescribeRequestPayload(r)
	if err != nil {
		return nil, err
	}

	var resp *svcsdk.DescribeAppOutput
	resp, err = rm.sdkapi.DescribeAppWithContext(ctx, input)
	rm.metrics.RecordAPICall("READ_ONE", "DescribeApp", err)
	if err != nil {
		if awsErr, ok := ackerr.AWSError(err); ok && awsErr.Code() == "ResourceNotFound" {
			return nil, ackerr.NotFound
		}
		return nil, err
	}

	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := r.ko.DeepCopy()

	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if resp.AppArn != nil {
		arn := ackv1alpha1.AWSResourceName(*resp.AppArn)
		ko.Status.ACKResourceMetadata.ARN = &arn
	}
	if resp.AppName != nil {
		ko.Spec.AppName = resp.AppName
	} else {
		ko.Spec.AppName = nil
	}
	if resp.AppType != nil {
		ko.Spec.AppType = resp.AppType
	} else {
		ko.Spec.AppType = nil
	}
	if resp.DomainId != nil {
		ko.Spec.DomainID = resp.DomainId
	} else {
		ko.Spec.DomainID = nil
	}
	if resp.ResourceSpec != nil {
		f8 := &svcapitypes.ResourceSpec{}
		if resp.ResourceSpec.InstanceType != nil {
			f8.InstanceType = resp.ResourceSpec.InstanceType
		}
		if resp.ResourceSpec.LifecycleConfigArn != nil {
			f8.LifecycleConfigARN = resp.ResourceSpec.LifecycleConfigArn
		}
		if resp.ResourceSpec.SageMakerImageArn != nil {
			f8.SageMakerImageARN = resp.ResourceSpec.SageMakerImageArn
		}
		if resp.ResourceSpec.SageMakerImageVersionArn != nil {
			f8.SageMakerImageVersionARN = resp.ResourceSpec.SageMakerImageVersionArn
		}
		ko.Spec.ResourceSpec = f8
	} else {
		ko.Spec.ResourceSpec = nil
	}
	if resp.Status != nil {
		ko.Status.Status = resp.Status
	} else {
		ko.Status.Status = nil
	}
	if resp.UserProfileName != nil {
		ko.Spec.UserProfileName = resp.UserProfileName
	} else {
		ko.Spec.UserProfileName = nil
	}

	rm.setStatusDefaults(ko)
	rm.customDescribeAppSetOutput(ko)
	return &resource{ko}, nil
}

// requiredFieldsMissingFromReadOneInput returns true if there are any fields
// for the ReadOne Input shape that are required but not present in the
// resource's Spec or Status
func (rm *resourceManager) requiredFieldsMissingFromReadOneInput(
	r *resource,
) bool {
	return r.ko.Spec.DomainID == nil || r.ko.Spec.UserProfileName == nil || r.ko.Spec.AppType == nil || r.ko.Spec.AppName == nil

}

// newDescribeRequestPayload returns SDK-specific struct for the HTTP request
// payload of the Describe API call for the resource
func (rm *resourceManager) newDescribeRequestPayload(
	r *resource,
) (*svcsdk.DescribeAppInput, error) {
	res := &svcsdk.DescribeAppInput{}

	if r.ko.Spec.AppName != nil {
		res.SetAppName(*r.ko.Spec.AppName)
	}
	if r.ko.Spec.AppType != nil {
		res.SetAppType(*r.ko.Spec.AppType)
	}
	if r.ko.Spec.DomainID != nil {
		res.SetDomainId(*r.ko.Spec.DomainID)
	}
	if r.ko.Spec.UserProfileName != nil {
		res.SetUserProfileName(*r.ko.Spec.UserProfileName)
	}

	return res, nil
}

// sdkCreate creates the supplied resource in the backend AWS service API and
// returns a copy of the resource with resource fields (in both Spec and
// Status) filled in with values from the CREATE API operation's Output shape.
func (rm *resourceManager) sdkCreate(
	ctx context.Context,
	desired *resource,
) (created *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkCreate")
	defer exit(err)
	input, err := rm.newCreateRequestPayload(ctx, desired)
	if err != nil {
		return nil, err
	}

	var resp *svcsdk.CreateAppOutput
	_ = resp
	resp, err = rm.sdkapi.CreateAppWithContext(ctx, input)
	rm.metrics.RecordAPICall("CREATE", "CreateApp", err)
	if err != nil {
		return nil, err
	}
	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := desired.ko.DeepCopy()

	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if resp.AppArn != nil {
		arn := ackv1alpha1.AWSResourceName(*resp.AppArn)
		ko.Status.ACKResourceMetadata.ARN = &arn
	}

	rm.setStatusDefaults(ko)
	return &resource{ko}, nil
}

// newCreateRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Create API call for the resource
func (rm *resourceManager) newCreateRequestPayload(
	ctx context.Context,
	r *resource,
) (*svcsdk.CreateAppInput, error) {
	res := &svcsdk.CreateAppInput{}

	if r.ko.Spec.AppName != nil {
		res.SetAppName(*r.ko.Spec.AppName)
	}
	if r.ko.Spec.AppType != nil {
		res.SetAppType(*r.ko.Spec.AppType)
	}
	if r.ko.Spec.DomainID != nil {
		res.SetDomainId(*r.ko.Spec.DomainID)
	}
	if r.ko.Spec.ResourceSpec != nil {
		f3 := &svcsdk.ResourceSpec{}
		if r.ko.Spec.ResourceSpec.InstanceType != nil {
			f3.SetInstanceType(*r.ko.Spec.ResourceSpec.InstanceType)
		}
		if r.ko.Spec.ResourceSpec.LifecycleConfigARN != nil {
			f3.SetLifecycleConfigArn(*r.ko.Spec.ResourceSpec.LifecycleConfigARN)
		}
		if r.ko.Spec.ResourceSpec.SageMakerImageARN != nil {
			f3.SetSageMakerImageArn(*r.ko.Spec.ResourceSpec.SageMakerImageARN)
		}
		if r.ko.Spec.ResourceSpec.SageMakerImageVersionARN != nil {
			f3.SetSageMakerImageVersionArn(*r.ko.Spec.ResourceSpec.SageMakerImageVersionARN)
		}
		res.SetResourceSpec(f3)
	}
	if r.ko.Spec.Tags != nil {
		f4 := []*svcsdk.Tag{}
		for _, f4iter := range r.ko.Spec.Tags {
			f4elem := &svcsdk.Tag{}
			if f4iter.Key != nil {
				f4elem.SetKey(*f4iter.Key)
			}
			if f4iter.Value != nil {
				f4elem.SetValue(*f4iter.Value)
			}
			f4 = append(f4, f4elem)
		}
		res.SetTags(f4)
	}
	if r.ko.Spec.UserProfileName != nil {
		res.SetUserProfileName(*r.ko.Spec.UserProfileName)
	}

	return res, nil
}

// sdkUpdate patches the supplied resource in the backend AWS service API and
// returns a new resource with updated fields.
func (rm *resourceManager) sdkUpdate(
	ctx context.Context,
	desired *resource,
	latest *resource,
	delta *ackcompare.Delta,
) (*resource, error) {
	// TODO(jaypipes): Figure this out...
	return nil, ackerr.NotImplemented
}

// sdkDelete deletes the supplied resource in the backend AWS service API
func (rm *resourceManager) sdkDelete(
	ctx context.Context,
	r *resource,
) (latest *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkDelete")
	defer exit(err)
	latestStatus := r.ko.Status.Status
	if latestStatus != nil && *latestStatus == svcsdk.AppStatusDeleted {
		return nil, nil
	}

	if err = rm.requeueUntilCanModify(ctx, r); err != nil {
		return r, err
	}

	input, err := rm.newDeleteRequestPayload(r)
	if err != nil {
		return nil, err
	}
	var resp *svcsdk.DeleteAppOutput
	_ = resp
	resp, err = rm.sdkapi.DeleteAppWithContext(ctx, input)
	rm.metrics.RecordAPICall("DELETE", "DeleteApp", err)

	if err == nil {
		if observed, err := rm.sdkFind(ctx, r); err != ackerr.NotFound {
			if err != nil {
				return nil, err
			}
			r.SetStatus(observed)
			return r, requeueWaitWhileDeleting
		}
	}

	return nil, err
}

// newDeleteRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Delete API call for the resource
func (rm *resourceManager) newDeleteRequestPayload(
	r *resource,
) (*svcsdk.DeleteAppInput, error) {
	res := &svcsdk.DeleteAppInput{}

	if r.ko.Spec.AppName != nil {
		res.SetAppName(*r.ko.Spec.AppName)
	}
	if r.ko.Spec.AppType != nil {
		res.SetAppType(*r.ko.Spec.AppType)
	}
	if r.ko.Spec.DomainID != nil {
		res.SetDomainId(*r.ko.Spec.DomainID)
	}
	if r.ko.Spec.UserProfileName != nil {
		res.SetUserProfileName(*r.ko.Spec.UserProfileName)
	}

	return res, nil
}

// setStatusDefaults sets default properties into supplied custom resource
func (rm *resourceManager) setStatusDefaults(
	ko *svcapitypes.App,
) {
	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if ko.Status.ACKResourceMetadata.Region == nil {
		ko.Status.ACKResourceMetadata.Region = &rm.awsRegion
	}
	if ko.Status.ACKResourceMetadata.OwnerAccountID == nil {
		ko.Status.ACKResourceMetadata.OwnerAccountID = &rm.awsAccountID
	}
	if ko.Status.Conditions == nil {
		ko.Status.Conditions = []*ackv1alpha1.Condition{}
	}
}

// updateConditions returns updated resource, true; if conditions were updated
// else it returns nil, false
func (rm *resourceManager) updateConditions(
	r *resource,
	onSuccess bool,
	err error,
) (*resource, bool) {
	ko := r.ko.DeepCopy()
	rm.setStatusDefaults(ko)

	// Terminal condition
	var terminalCondition *ackv1alpha1.Condition = nil
	var recoverableCondition *ackv1alpha1.Condition = nil
	var syncCondition *ackv1alpha1.Condition = nil
	for _, condition := range ko.Status.Conditions {
		if condition.Type == ackv1alpha1.ConditionTypeTerminal {
			terminalCondition = condition
		}
		if condition.Type == ackv1alpha1.ConditionTypeRecoverable {
			recoverableCondition = condition
		}
		if condition.Type == ackv1alpha1.ConditionTypeResourceSynced {
			syncCondition = condition
		}
	}
	var termError *ackerr.TerminalError
	if rm.terminalAWSError(err) || err == ackerr.SecretTypeNotSupported || err == ackerr.SecretNotFound || errors.As(err, &termError) {
		if terminalCondition == nil {
			terminalCondition = &ackv1alpha1.Condition{
				Type: ackv1alpha1.ConditionTypeTerminal,
			}
			ko.Status.Conditions = append(ko.Status.Conditions, terminalCondition)
		}
		var errorMessage = ""
		if err == ackerr.SecretTypeNotSupported || err == ackerr.SecretNotFound || errors.As(err, &termError) {
			errorMessage = err.Error()
		} else {
			awsErr, _ := ackerr.AWSError(err)
			errorMessage = awsErr.Error()
		}
		terminalCondition.Status = corev1.ConditionTrue
		terminalCondition.Message = &errorMessage
	} else {
		// Clear the terminal condition if no longer present
		if terminalCondition != nil {
			terminalCondition.Status = corev1.ConditionFalse
			terminalCondition.Message = nil
		}
		// Handling Recoverable Conditions
		if err != nil {
			if recoverableCondition == nil {
				// Add a new Condition containing a non-terminal error
				recoverableCondition = &ackv1alpha1.Condition{
					Type: ackv1alpha1.ConditionTypeRecoverable,
				}
				ko.Status.Conditions = append(ko.Status.Conditions, recoverableCondition)
			}
			recoverableCondition.Status = corev1.ConditionTrue
			awsErr, _ := ackerr.AWSError(err)
			errorMessage := err.Error()
			if awsErr != nil {
				errorMessage = awsErr.Error()
			}
			recoverableCondition.Message = &errorMessage
		} else if recoverableCondition != nil {
			recoverableCondition.Status = corev1.ConditionFalse
			recoverableCondition.Message = nil
		}
	}
	// Required to avoid the "declared but not used" error in the default case
	_ = syncCondition
	if terminalCondition != nil || recoverableCondition != nil || syncCondition != nil {
		return &resource{ko}, true // updated
	}
	return nil, false // not updated
}

// terminalAWSError returns awserr, true; if the supplied error is an aws Error type
// and if the exception indicates that it is a Terminal exception
// 'Terminal' exception are specified in generator configuration
func (rm *resourceManager) terminalAWSError(err error) bool {
	if err == nil {
		return false
	}
	awsErr, ok := ackerr.AWSError(err)
	if !ok {
		return false
	}
	switch awsErr.Code() {
	case "ResourceNotFound",
		"InvalidParameterCombination",
		"InvalidParameterValue",
		"MissingParameter":
		return true
	default:
		return false
	}
}
