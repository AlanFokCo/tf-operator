// Copyright 2018 The Kubeflow Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tensorflow

import (
	"fmt"
	common "github.com/kubeflow/tf-operator/pkg/apis/common/v1"

	tfv1 "github.com/kubeflow/tf-operator/pkg/apis/tensorflow/v1"
)

var (
	errPortNotFound = fmt.Errorf("failed to found the port")
)

// GetPortFromTFJob gets the port of tensorflow container.
func GetPortFromTFJob(tfJob *tfv1.TFJob, rtype tfv1.TFReplicaType) (int32, error) {
	containers := tfJob.Spec.TFReplicaSpecs[rtype].Template.Spec.Containers
	for _, container := range containers {
		if container.Name == tfv1.DefaultContainerName {
			ports := container.Ports
			for _, port := range ports {
				if port.Name == tfv1.DefaultPortName {
					return port.ContainerPort, nil
				}
			}
		}
	}
	return -1, errPortNotFound
}

// ContainChieforMasterSpec returns true if the tfjob contains chief or master spec.
func ContainChieforMasterSpec(tfJob *tfv1.TFJob) bool {
	if _, ok := tfJob.Spec.TFReplicaSpecs[tfv1.TFReplicaTypeChief]; ok {
		return true
	} else if _, ok := tfJob.Spec.TFReplicaSpecs[tfv1.TFReplicaTypeMaster]; ok {
		return true
	}
	return false
}

func CheckTFJobIsNotPending(tfJob *tfv1.TFJob) bool {
	checkResult := false

	if tfJob.Status.Conditions == nil {
		return checkResult
	}
	tfJobConditions := tfJob.Status.Conditions
	for _, condition := range tfJobConditions {
		if condition.Type == common.JobRunning {
			checkResult = true
		}
	}

	return checkResult
}
