/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package iamauth

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	iamauthv1 "sigs.k8s.io/aws-iam-authenticator/pkg/mapper/crd/apis/iamauthenticator/v1alpha1"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type crdBackend struct {
	client crclient.Client
}

func (b *crdBackend) MapRole(mapping RoleMapping) error {
	ctx := context.TODO()

	if err := mapping.Validate(); err != nil {
		return err
	}

	mappingList := iamauthv1.IAMIdentityMappingList{}
	err := b.client.List(ctx, &mappingList)
	if err != nil {
		return fmt.Errorf("getting list of mappings: %w", err)
	}

	for _, existingMapping := range mappingList.Items {
		existing := existingMapping
		if roleMappingMatchesIAMMap(mapping, &existing) {
			// We already have a mapping so do nothing
			return nil
		}
	}

	iamMapping := &iamauthv1.IAMIdentityMapping{
		ObjectMeta: v1.ObjectMeta{
			Namespace:    metav1.NamespaceSystem,
			GenerateName: "capa-iamauth-",
		},
		Spec: iamauthv1.IAMIdentityMappingSpec{
			ARN:      mapping.RoleARN,
			Username: mapping.UserName,
			Groups:   mapping.Groups,
		},
	}

	return b.client.Create(ctx, iamMapping)
}

func (b *crdBackend) MapUser(mapping UserMapping) error {
	ctx := context.TODO()

	if err := mapping.Validate(); err != nil {
		return err
	}

	mappingList := iamauthv1.IAMIdentityMappingList{}
	err := b.client.List(ctx, &mappingList)
	if err != nil {
		return fmt.Errorf("getting list of mappings: %w", err)
	}

	for _, existingMapping := range mappingList.Items {
		existing := existingMapping
		if userMappingMatchesIAMMap(mapping, &existing) {
			// We already have a mapping so do nothing
			return nil
		}
	}

	iamMapping := &iamauthv1.IAMIdentityMapping{
		ObjectMeta: v1.ObjectMeta{
			Namespace:    metav1.NamespaceSystem,
			GenerateName: "capa-iamauth-",
		},
		Spec: iamauthv1.IAMIdentityMappingSpec{
			ARN:      mapping.UserARN,
			Username: mapping.UserName,
			Groups:   mapping.Groups,
		},
	}

	return b.client.Create(ctx, iamMapping)
}

func roleMappingMatchesIAMMap(mapping RoleMapping, iamMapping *iamauthv1.IAMIdentityMapping) bool {
	if mapping.RoleARN != iamMapping.Spec.ARN {
		return false
	}

	if mapping.UserName != iamMapping.Spec.Username {
		return false
	}

	if len(mapping.Groups) != len(iamMapping.Spec.Groups) {
		return false
	}

	for _, mappingGroup := range mapping.Groups {
		found := false
		for _, iamGroup := range iamMapping.Spec.Groups {
			if iamGroup == mappingGroup {
				found = true
			}
		}
		if !found {
			return false
		}
	}

	return true
}

func userMappingMatchesIAMMap(mapping UserMapping, iamMapping *iamauthv1.IAMIdentityMapping) bool {
	if mapping.UserARN != iamMapping.Spec.ARN {
		return false
	}

	if mapping.UserName != iamMapping.Spec.Username {
		return false
	}

	if len(mapping.Groups) != len(iamMapping.Spec.Groups) {
		return false
	}

	for _, mappingGroup := range mapping.Groups {
		found := false
		for _, iamGroup := range iamMapping.Spec.Groups {
			if iamGroup == mappingGroup {
				found = true
			}
		}
		if !found {
			return false
		}
	}

	return true
}
