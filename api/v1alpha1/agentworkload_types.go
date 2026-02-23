/*
Copyright 2026.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AgentWorkloadSpec defines the desired state of AgentWorkload
type AgentWorkloadSpec struct {
	// workloadType defines the infrastructure type (generic, ceph, minio, postgres, etc.)
	// +kubebuilder:validation:Enum=generic;ceph;minio;postgres;aws;kubernetes
	// +required
	WorkloadType string `json:"workloadType"`

	// mcpServerEndpoint is the HTTP endpoint of the MCP server (e.g. "http://mcp-server:8000")
	// +kubebuilder:validation:Pattern=`^https?://[a-zA-Z0-9.-]+(:[0-9]+)?$`
	// +required
	MCPServerEndpoint string `json:"mcpServerEndpoint"`

	// objective is the high-level goal for the agent (e.g. "optimize cluster performance")
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1000
	// +required
	Objective string `json:"objective"`

	// agents is a list of agent names to run for this workload
	// +optional
	Agents []string `json:"agents,omitempty"`

	// autoApproveThreshold is the confidence level (0-1) above which actions are auto-approved (as string, e.g. "0.95")
	// +kubebuilder:validation:Pattern=`^0(\.[0-9]{1,2})?$|^1(\.0{1,2})?$`
	// +kubebuilder:default="0.95"
	// +optional
	AutoApproveThreshold *string `json:"autoApproveThreshold,omitempty"`

	// opaPolicy controls the safety policy for action execution
	// +kubebuilder:validation:Enum=strict;permissive
	// +kubebuilder:default=strict
	// +optional
	OPAPolicy *string `json:"opaPolicy,omitempty"`
}

// Action represents a proposed or executed action by an agent
type Action struct {
	// name is a unique identifier for this action
	Name string `json:"name"`

	// description explains what the action does
	Description string `json:"description"`

	// confidence is the agent's confidence in this action (0-1, as string e.g. "0.95")
	// +kubebuilder:validation:Pattern=`^0(\.[0-9]{1,2})?$|^1(\.0{1,2})?$`
	Confidence string `json:"confidence"`

	// timestamp is when the action was proposed
	// +optional
	Timestamp *metav1.Time `json:"timestamp,omitempty"`

	// approved indicates if this action was approved
	Approved *bool `json:"approved,omitempty"`
}

// AgentWorkloadStatus defines the observed state of AgentWorkload.
type AgentWorkloadStatus struct {
	// phase is the current lifecycle phase: Pending, Running, Completed, Failed
	// +optional
	Phase string `json:"phase,omitempty"`

	// readyAgents is the number of agents ready to execute
	// +optional
	ReadyAgents int32 `json:"readyAgents,omitempty"`

	// lastReconcileTime is the last time the resource was reconciled
	// +optional
	LastReconcileTime *metav1.Time `json:"lastReconcileTime,omitempty"`

	// proposedActions is a list of actions proposed by agents
	// +optional
	ProposedActions []Action `json:"proposedActions,omitempty"`

	// executedActions is a list of actions that were approved and executed
	// +optional
	ExecutedActions []Action `json:"executedActions,omitempty"`

	// conditions represent the current state of the AgentWorkload resource.
	// Each condition has a unique type and reflects the status of a specific aspect of the resource.
	//
	// Standard condition types include:
	// - "Available": the resource is fully functional
	// - "Progressing": the resource is being created or updated
	// - "Degraded": the resource failed to reach or maintain its desired state
	//
	// The status of each condition is one of True, False, or Unknown.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// AgentWorkload is the Schema for the agentworkloads API
type AgentWorkload struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of AgentWorkload
	// +required
	Spec AgentWorkloadSpec `json:"spec"`

	// status defines the observed state of AgentWorkload
	// +optional
	Status AgentWorkloadStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// AgentWorkloadList contains a list of AgentWorkload
type AgentWorkloadList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []AgentWorkload `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AgentWorkload{}, &AgentWorkloadList{})
}
