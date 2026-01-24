// Package gtm provides helper functions for building GTM API parameters.
package gtm

import (
	"fmt"
	"regexp"
	"strings"
)

// Parameter type constants
const (
	ParameterTypeTemplate         = "TEMPLATE"
	ParameterTypeBoolean          = "BOOLEAN"
	ParameterTypeInteger          = "INTEGER"
	ParameterTypeList             = "LIST"
	ParameterTypeMap              = "MAP"
	ParameterTypeTagReference     = "TAG_REFERENCE"
	ParameterTypeTriggerReference = "TRIGGER_REFERENCE"
)

// BuildTemplateParameter creates a template parameter.
// Template parameters can contain variable references like {{Variable Name}}.
func BuildTemplateParameter(key, value string) Parameter {
	return Parameter{
		Type:  ParameterTypeTemplate,
		Key:   key,
		Value: value,
	}
}

// BuildBooleanParameter creates a boolean parameter.
func BuildBooleanParameter(key string, value bool) Parameter {
	boolStr := "false"
	if value {
		boolStr = "true"
	}
	return Parameter{
		Type:  ParameterTypeBoolean,
		Key:   key,
		Value: boolStr,
	}
}

// BuildIntegerParameter creates an integer parameter.
func BuildIntegerParameter(key string, value int) Parameter {
	return Parameter{
		Type:  ParameterTypeInteger,
		Key:   key,
		Value: fmt.Sprintf("%d", value),
	}
}

// BuildListParameter creates a list parameter.
func BuildListParameter(key string, items []Parameter) Parameter {
	return Parameter{
		Type: ParameterTypeList,
		Key:  key,
		List: items,
	}
}

// BuildMapParameter creates a map parameter.
func BuildMapParameter(keyValuePairs []Parameter) Parameter {
	return Parameter{
		Type: ParameterTypeMap,
		Map:  keyValuePairs,
	}
}

// BuildTagReferenceParameter creates a tag reference parameter.
func BuildTagReferenceParameter(key, tagName string) Parameter {
	return Parameter{
		Type:  ParameterTypeTagReference,
		Key:   key,
		Value: tagName,
	}
}

// BuildTriggerReferenceParameter creates a trigger reference parameter.
func BuildTriggerReferenceParameter(triggerID string) Parameter {
	return Parameter{
		Type:  ParameterTypeTriggerReference,
		Value: triggerID,
	}
}

// BuildEventParameter creates a GA4 event parameter map.
// This is the nested map structure required for GA4 event parameters.
func BuildEventParameter(name, value string) Parameter {
	return BuildMapParameter([]Parameter{
		BuildTemplateParameter("name", name),
		BuildTemplateParameter("value", value),
	})
}

// BuildEventParametersList creates a list of GA4 event parameters.
// Each input map should have "name" and "value" keys.
func BuildEventParametersList(params map[string]string) []Parameter {
	result := make([]Parameter, 0, len(params))
	for name, value := range params {
		result = append(result, BuildEventParameter(name, value))
	}
	return result
}

// BuildScrollPercentageList creates the scroll percentage list parameter for scroll depth triggers.
func BuildScrollPercentageList(percentages []int) Parameter {
	items := make([]Parameter, 0, len(percentages))
	for _, pct := range percentages {
		items = append(items, Parameter{
			Type:  ParameterTypeTemplate,
			Value: fmt.Sprintf("%d", pct),
		})
	}
	return Parameter{
		Type: ParameterTypeList,
		Key:  "verticalScrollPercentageList",
		List: items,
	}
}

// BuildCustomEventFilter creates a custom event filter for customEvent triggers.
// The matchType parameter can be "EQUALS", "CONTAINS", "MATCHES_REGEX", etc.
func BuildCustomEventFilter(eventName string, matchType string) []Condition {
	if matchType == "" {
		matchType = "EQUALS"
	}
	return []Condition{
		{
			Type: matchType,
			Parameter: []Parameter{
				BuildTemplateParameter("arg0", "{{_event}}"),
				BuildTemplateParameter("arg1", eventName),
			},
		},
	}
}

// BuildURLFilter creates a URL-based filter for triggers.
func BuildURLFilter(variable, matchType, pattern string) []Condition {
	return []Condition{
		{
			Type: matchType,
			Parameter: []Parameter{
				BuildTemplateParameter("arg0", variable),
				BuildTemplateParameter("arg1", pattern),
			},
		},
	}
}

// BuildPageViewFilter creates a filter for the Page URL variable.
func BuildPageViewFilter(matchType, pattern string) []Condition {
	return BuildURLFilter("{{Page URL}}", matchType, pattern)
}

// BuildClickFilter creates a filter for click triggers.
func BuildClickFilter(matchType, pattern string, clickProperty string) []Condition {
	if clickProperty == "" {
		clickProperty = "{{Click URL}}"
	}
	return BuildURLFilter(clickProperty, matchType, pattern)
}

// Path building functions
// Note: BuildTagPath, BuildTriggerPath, BuildVariablePath are defined in mutations.go

// BuildVersionPath constructs a version path from IDs.
func BuildVersionPath(accountID, containerID, versionID string) string {
	return fmt.Sprintf("accounts/%s/containers/%s/versions/%s",
		accountID, containerID, versionID)
}

// ExtractIDFromPath extracts a resource ID from a GTM path.
// resourceType should be "account", "container", "workspace", "tag", "trigger", "variable", or "version".
func ExtractIDFromPath(path, resourceType string) (string, error) {
	parts := strings.Split(path, "/")
	plural := resourceType + "s"

	for i, part := range parts {
		if part == plural && i+1 < len(parts) {
			return parts[i+1], nil
		}
	}

	return "", fmt.Errorf("could not extract %s ID from path: %s", resourceType, path)
}

// ParseWorkspacePath parses a workspace path into its components.
func ParseWorkspacePath(path string) (accountID, containerID, workspaceID string, err error) {
	parts := strings.Split(path, "/")
	if len(parts) != 6 {
		return "", "", "", fmt.Errorf("invalid workspace path format: expected accounts/{id}/containers/{id}/workspaces/{id}")
	}

	if parts[0] != "accounts" || parts[2] != "containers" || parts[4] != "workspaces" {
		return "", "", "", fmt.Errorf("invalid workspace path format: %s", path)
	}

	return parts[1], parts[3], parts[5], nil
}

// ParseContainerPath parses a container path into its components.
func ParseContainerPath(path string) (accountID, containerID string, err error) {
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		return "", "", fmt.Errorf("invalid container path format: expected accounts/{id}/containers/{id}")
	}

	if parts[0] != "accounts" || parts[2] != "containers" {
		return "", "", fmt.Errorf("invalid container path format: %s", path)
	}

	return parts[1], parts[3], nil
}

// GA4-specific helper functions

// BuildGA4ConfigTag creates the parameters for a GA4 Configuration tag.
func BuildGA4ConfigTag(measurementID string, sendPageView bool) []Parameter {
	return []Parameter{
		BuildTemplateParameter("measurementId", measurementID),
		BuildBooleanParameter("sendPageView", sendPageView),
	}
}

// BuildGA4EventTag creates the parameters for a GA4 Event tag.
func BuildGA4EventTag(configTagName, eventName string, params map[string]string, sendEcommerce bool) []Parameter {
	result := []Parameter{
		BuildTagReferenceParameter("measurementId", configTagName),
		BuildTemplateParameter("eventName", eventName),
	}

	if len(params) > 0 {
		eventParams := BuildEventParametersList(params)
		result = append(result, BuildListParameter("eventParameters", eventParams))
	}

	if sendEcommerce {
		result = append(result, BuildBooleanParameter("sendEcommerceData", true))
	}

	return result
}

// MergeParameters merges multiple parameter slices, with later slices overriding earlier ones.
// Parameters are matched by key.
func MergeParameters(paramLists ...[]Parameter) []Parameter {
	merged := make(map[string]Parameter)
	var noKeyParams []Parameter

	for _, params := range paramLists {
		for _, param := range params {
			if param.Key != "" {
				merged[param.Key] = param
			} else {
				noKeyParams = append(noKeyParams, param)
			}
		}
	}

	result := make([]Parameter, 0, len(merged)+len(noKeyParams))
	for _, param := range merged {
		result = append(result, param)
	}
	result = append(result, noKeyParams...)

	return result
}

// IsVariableReference checks if a string is a GTM variable reference.
// Variable references are in the format {{Variable Name}}.
func IsVariableReference(s string) bool {
	matched, _ := regexp.MatchString(`^\{\{.+\}\}$`, s)
	return matched
}

// WrapVariableReference wraps a variable name in GTM variable reference syntax.
func WrapVariableReference(variableName string) string {
	if IsVariableReference(variableName) {
		return variableName
	}
	return "{{" + variableName + "}}"
}

// UnwrapVariableReference removes GTM variable reference syntax from a string.
func UnwrapVariableReference(ref string) string {
	if !IsVariableReference(ref) {
		return ref
	}
	return strings.TrimSuffix(strings.TrimPrefix(ref, "{{"), "}}")
}
