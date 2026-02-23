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
	"fmt"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// log is for logging in this package.
var agentworkloadlog = logf.Log.WithName("agentworkload-resource")

// Default implements DefaultingWebhook so a webhook will be registered for the type
func (r *AgentWorkload) Default() {
	agentworkloadlog.Info("default", "name", r.Name)

	// Set defaults if not provided
	if r.Spec.AutoApproveThreshold == nil {
		r.Spec.AutoApproveThreshold = stringPtr("0.95")
	}

	if r.Spec.OPAPolicy == nil {
		r.Spec.OPAPolicy = stringPtr("strict")
	}

	// Initialize status if empty
	if r.Status.Phase == "" {
		r.Status.Phase = "Pending"
	}
	if r.Status.ProposedActions == nil {
		r.Status.ProposedActions = []Action{}
	}
	if r.Status.ExecutedActions == nil {
		r.Status.ExecutedActions = []Action{}
	}
}

// ValidateCreate validates the resource on creation
func (r *AgentWorkload) ValidateCreate() error {
	agentworkloadlog.Info("validate create", "name", r.Name)
	return r.validate()
}

// ValidateUpdate validates the resource on update
func (r *AgentWorkload) ValidateUpdate(old runtime.Object) error {
	agentworkloadlog.Info("validate update", "name", r.Name)
	return r.validate()
}

// ValidateDelete validates the resource on deletion
func (r *AgentWorkload) ValidateDelete() error {
	agentworkloadlog.Info("validate delete", "name", r.Name)
	// Allow deletion
	return nil
}

// validate performs all validation checks on AgentWorkload
func (r *AgentWorkload) validate() error {
	var allErrs []string

	// 1. Validate workloadType
	validWorkloadTypes := []string{"generic", "ceph", "minio", "postgres", "aws", "kubernetes"}
	if !isStringInSlice(r.Spec.WorkloadType, validWorkloadTypes) {
		allErrs = append(allErrs, fmt.Sprintf("workloadType must be one of %v, got %q", validWorkloadTypes, r.Spec.WorkloadType))
	}

	// 2. Validate mcpServerEndpoint
	if err := validateMCPEndpoint(r.Spec.MCPServerEndpoint); err != nil {
		allErrs = append(allErrs, err.Error())
	}

	// 3. Validate objective
	if len(r.Spec.Objective) == 0 {
		allErrs = append(allErrs, "objective must not be empty")
	}
	if len(r.Spec.Objective) > 1000 {
		allErrs = append(allErrs, fmt.Sprintf("objective must be ≤ 1000 characters, got %d", len(r.Spec.Objective)))
	}

	// 4. Validate agents list
	if len(r.Spec.Agents) == 0 {
		allErrs = append(allErrs, "agents list must not be empty")
	}

	// 5. Validate autoApproveThreshold
	if r.Spec.AutoApproveThreshold != nil {
		if err := validateThreshold(*r.Spec.AutoApproveThreshold); err != nil {
			allErrs = append(allErrs, err.Error())
		}
	}

	// 6. Validate opaPolicy
	if r.Spec.OPAPolicy != nil {
		validPolicies := []string{"strict", "permissive"}
		if !isStringInSlice(*r.Spec.OPAPolicy, validPolicies) {
			allErrs = append(allErrs, fmt.Sprintf("opaPolicy must be one of %v, got %q", validPolicies, *r.Spec.OPAPolicy))
		}
	}

	// Combine errors
	if len(allErrs) > 0 {
		errMsg := strings.Join(allErrs, "; ")
		return apierrors.NewInvalid(
			r.GroupVersionKind().GroupKind(),
			r.Name,
			field.ErrorList{
				field.InternalError(field.NewPath("spec"), fmt.Errorf("%s", errMsg)),
			},
		)
	}

	return nil
}

// validateMCPEndpoint validates the MCP server endpoint
func validateMCPEndpoint(endpoint string) error {
	if endpoint == "" {
		return fmt.Errorf("mcpServerEndpoint must not be empty")
	}

	// Check URL format
	parsedURL, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("mcpServerEndpoint is not a valid URL: %v", err)
	}

	// Check scheme is http or https
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("mcpServerEndpoint scheme must be http or https, got %q", parsedURL.Scheme)
	}

	// Check host is not empty
	if parsedURL.Host == "" {
		return fmt.Errorf("mcpServerEndpoint host is empty")
	}

	// Try to reach the endpoint (with short timeout)
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Head(endpoint + "/tools")
	if err != nil {
		// Warning only - don't fail validation if endpoint unreachable (might be created later)
		// This is expected in dev environments where MCP server starts after operator
		return nil // Allow creation; endpoint will be checked at runtime
	}
	defer resp.Body.Close()

	return nil
}

// validateThreshold validates the autoApproveThreshold value
func validateThreshold(threshold string) error {
	if threshold == "" {
		return fmt.Errorf("autoApproveThreshold cannot be empty")
	}

	// Parse as float
	val, err := strconv.ParseFloat(threshold, 64)
	if err != nil {
		return fmt.Errorf("autoApproveThreshold must be a valid number, got %q", threshold)
	}

	// Check range 0.0-1.0
	if val < 0.0 || val > 1.0 {
		return fmt.Errorf("autoApproveThreshold must be between 0.0 and 1.0, got %v", val)
	}

	// Validate format (optional: check decimal places)
	if !isValidThresholdFormat(threshold) {
		return fmt.Errorf("autoApproveThreshold format invalid, expected format like '0.95'")
	}

	return nil
}

// isValidThresholdFormat checks if threshold has valid format (max 2 decimal places)
func isValidThresholdFormat(s string) bool {
	pattern := `^[01](\.\d{1,2})?$`
	matched, _ := regexp.MatchString(pattern, s)
	return matched
}

// isStringInSlice checks if a string is in a slice
func isStringInSlice(s string, slice []string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

// stringPtr returns a pointer to a string
func stringPtr(s string) *string {
	return &s
}

// IsEndpointReachable attempts to verify the MCP endpoint is reachable
func IsEndpointReachable(endpoint string) bool {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Head(endpoint + "/tools")
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode >= 200 && resp.StatusCode < 400
}

// ResolveEndpointIP attempts to resolve DNS and return IP
func ResolveEndpointIP(endpoint string) (string, error) {
	parsedURL, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}

	// Extract hostname (without port)
	host := parsedURL.Hostname()
	if host == "" {
		return "", fmt.Errorf("no hostname in endpoint")
	}

	// Resolve DNS
	ips, err := net.LookupIP(host)
	if err != nil {
		return "", err
	}

	if len(ips) == 0 {
		return "", fmt.Errorf("no IPs found for hostname %s", host)
	}

	return ips[0].String(), nil
}
