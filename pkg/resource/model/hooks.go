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

package model

import (
	"context"

	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	svcapitypes "github.com/aws-controllers-k8s/sagemaker-controller/apis/v1alpha1"
	svcsdk "github.com/aws/aws-sdk-go/service/sagemaker"
)

// func (rm *resourceManager) customUpdateModel(
// 	ctx context.Context,
// 	desired *resource,
// 	latest *resource,
// 	delta *ackcompare.Delta,
// ) (updated *resource, err error) {
// 	rlog := ackrtlog.FromContext(ctx)
// 	exit := rlog.Trace("rm.customUpdateModel")
// 	defer exit(err)

// 	// Default `updated` to `desired` because it is likely
// 	// EC2 `modify` APIs do NOT return output, only errors.
// 	// If the `modify` calls (i.e. `sync`) do NOT return
// 	// an error, then the update was successful and desired.Spec
// 	// (now updated.Spec) reflects the latest resource state.
// 	updated = rm.concreteResource(desired.DeepCopy())

// 	if delta.DifferentAt("Spec.Tags") {
// 		if err := rm.syncTags(ctx, desired, latest); err != nil {
// 			return nil, err
// 		}
// 	}

// 	return updated, nil
// }

// syncTags used to keep tags in sync by calling Create and Delete API's
func (rm *resourceManager) syncTags(
	ctx context.Context,
	desired *resource,
	latest *resource,
) (resp *svcsdk.DeleteTagsOutput, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.syncTags")
	defer func(err error) {
		exit(err)
	}(err)

	resourceId := (*string)(latest.ko.Status.ACKResourceMetadata.ARN)

	toDelete := computeTagsDelta(
		desired.ko.Spec.Tags, latest.ko.Spec.Tags,
	)

	if len(toDelete) > 0 {
		rlog.Debug("removing tags from model resource", "tags", toDelete)

		// := reflect.ValueOf(toDelete).MapKeys()
		keys := make([]*string, len(toDelete))
		for i, raw_key := range toDelete {
			keys[i] = raw_key.Key
		}
		resp, err = rm.sdkapi.DeleteTagsWithContext(
			ctx,
			&svcsdk.DeleteTagsInput{
				ResourceArn: resourceId,
				TagKeys:     keys,
			},
		)
		rm.metrics.RecordAPICall("UPDATE", "DeleteTags", err)
		if err != nil {
			return nil, err
		}

		return resp, nil

	}

	// if len(toAdd) > 0 {
	// 	rlog.Debug("adding tags to model resource", "tags", toAdd)
	// 	_, err = rm.sdkapi.AddTagsWithContext(
	// 		ctx,
	// 		&svcsdk.AddTagsInput{
	// 			ResourceArn: resourceId,
	// 			Tags:        rm.sdkTags(toAdd),
	// 		},
	// 	)
	// 	rm.metrics.RecordAPICall("UPDATE", "CreateTags", err)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil, nil
}

// sdkTags converts *svcapitypes.Tag array to a *svcsdk.Tag array
func (rm *resourceManager) sdkTags(
	tags []*svcapitypes.Tag,
) (sdktags []*svcsdk.Tag) {

	for _, i := range tags {
		sdktag := rm.newTag(*i)
		sdktags = append(sdktags, sdktag)
	}

	return sdktags
}

func (rm *resourceManager) newTag(
	c svcapitypes.Tag,
) *svcsdk.Tag {
	res := &svcsdk.Tag{}
	if c.Key != nil {
		res.SetKey(*c.Key)
	}
	if c.Value != nil {
		res.SetValue(*c.Value)
	}

	return res
}

// computeTagsDelta returns tags to be added and removed from the resource
func computeTagsDelta(
	desired []*svcapitypes.Tag,
	latest []*svcapitypes.Tag,
) (toDelete []*svcapitypes.Tag) {

	desiredTags := map[string]string{}
	for _, tag := range desired {
		desiredTags[*tag.Key] = *tag.Value
	}

	latestTags := map[string]string{}
	for _, tag := range latest {
		latestTags[*tag.Key] = *tag.Value
	}

	// for _, tag := range desired {
	// 	val, ok := latestTags[*tag.Key]
	// 	if !ok || val != *tag.Value {
	// 		toAdd = append(toAdd, tag)
	// 	}
	// }

	for _, tag := range latest {
		_, ok := desiredTags[*tag.Key]
		if !ok {
			toDelete = append(toDelete, tag)
		}
	}

	return toDelete

}
