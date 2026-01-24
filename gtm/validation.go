package gtm

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// Validation limits
const (
	GTMNameMaxLength          = 256
	GTMNotesMaxLength         = 5000
	GA4EventNameMaxLength     = 40
	GA4ParameterNameMaxLength = 40
	AccountIDMinLength        = 10
)

// ValidateTagInput validates tag creation/update inputs.
func ValidateTagInput(name, tagType string, firingTriggerIDs []string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("tag name is required")
	}
	if len(name) > 256 {
		return fmt.Errorf("tag name must be 256 characters or less")
	}
	if strings.TrimSpace(tagType) == "" {
		return fmt.Errorf("tag type is required")
	}
	if len(firingTriggerIDs) == 0 {
		return fmt.Errorf("at least one firing trigger ID is required")
	}
	for _, id := range firingTriggerIDs {
		if strings.TrimSpace(id) == "" {
			return fmt.Errorf("firing trigger ID cannot be empty")
		}
	}
	return nil
}

// ValidateTriggerInput validates trigger creation inputs.
func ValidateTriggerInput(name, triggerType string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("trigger name is required")
	}
	if len(name) > 256 {
		return fmt.Errorf("trigger name must be 256 characters or less")
	}
	if strings.TrimSpace(triggerType) == "" {
		return fmt.Errorf("trigger type is required")
	}
	return nil
}

// ValidateVariableInput validates variable creation inputs.
func ValidateVariableInput(name, varType string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("variable name is required")
	}
	if len(name) > 256 {
		return fmt.Errorf("variable name must be 256 characters or less")
	}
	if strings.TrimSpace(varType) == "" {
		return fmt.Errorf("variable type is required")
	}
	return nil
}

// ValidateWorkspacePath validates workspace path components.
func ValidateWorkspacePath(accountID, containerID, workspaceID string) error {
	if strings.TrimSpace(accountID) == "" {
		return fmt.Errorf("account ID is required")
	}
	if strings.TrimSpace(containerID) == "" {
		return fmt.Errorf("container ID is required")
	}
	if strings.TrimSpace(workspaceID) == "" {
		return fmt.Errorf("workspace ID is required")
	}
	return nil
}

// BuildWorkspacePath constructs a workspace path from IDs.
func BuildWorkspacePath(accountID, containerID, workspaceID string) string {
	return fmt.Sprintf("accounts/%s/containers/%s/workspaces/%s",
		accountID, containerID, workspaceID)
}

// ValidateContainerPath validates container path components.
func ValidateContainerPath(accountID, containerID string) error {
	if strings.TrimSpace(accountID) == "" {
		return fmt.Errorf("account ID is required")
	}
	if strings.TrimSpace(containerID) == "" {
		return fmt.Errorf("container ID is required")
	}
	return nil
}

// BuildContainerPath constructs a container path from IDs.
func BuildContainerPath(accountID, containerID string) string {
	return fmt.Sprintf("accounts/%s/containers/%s", accountID, containerID)
}

// ValidateAccountID validates GTM account ID format.
// Account IDs must be numeric and at least 10 digits.
func ValidateAccountID(accountID string) error {
	if accountID == "" {
		return fmt.Errorf("account ID cannot be empty")
	}

	// Must be numeric
	for _, c := range accountID {
		if c < '0' || c > '9' {
			return fmt.Errorf("account ID must contain only digits, got: %s", accountID)
		}
	}

	// Must be at least 10 digits
	if len(accountID) < AccountIDMinLength {
		return fmt.Errorf("account ID must be at least %d digits, got: %s", AccountIDMinLength, accountID)
	}

	return nil
}

// ValidateContainerID validates GTM container ID format.
// Container IDs must be numeric.
func ValidateContainerID(containerID string) error {
	if containerID == "" {
		return fmt.Errorf("container ID cannot be empty")
	}

	for _, c := range containerID {
		if c < '0' || c > '9' {
			return fmt.Errorf("container ID must contain only digits, got: %s", containerID)
		}
	}

	return nil
}

// ValidateWorkspaceID validates GTM workspace ID format.
// Workspace IDs must be numeric.
func ValidateWorkspaceID(workspaceID string) error {
	if workspaceID == "" {
		return fmt.Errorf("workspace ID cannot be empty")
	}

	for _, c := range workspaceID {
		if c < '0' || c > '9' {
			return fmt.Errorf("workspace ID must contain only digits, got: %s", workspaceID)
		}
	}

	return nil
}

// ValidateVariableID validates GTM variable ID format.
// Variable IDs must be numeric.
func ValidateVariableID(variableID string) error {
	if variableID == "" {
		return fmt.Errorf("variable ID cannot be empty")
	}

	for _, c := range variableID {
		if c < '0' || c > '9' {
			return fmt.Errorf("variable ID must contain only digits, got: %s", variableID)
		}
	}

	return nil
}

// ValidateName validates an entity name.
// Names cannot be empty and must not exceed the max length.
func ValidateName(name string, fieldName string) error {
	if fieldName == "" {
		fieldName = "name"
	}

	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("%s cannot be empty", fieldName)
	}

	if len(name) > GTMNameMaxLength {
		return fmt.Errorf("%s exceeds maximum length of %d characters", fieldName, GTMNameMaxLength)
	}

	return nil
}

// ValidateNotes validates a notes/description field.
func ValidateNotes(notes string) error {
	if len(notes) > GTMNotesMaxLength {
		return fmt.Errorf("notes exceed maximum length of %d characters", GTMNotesMaxLength)
	}
	return nil
}

// ValidateTriggerType validates a trigger type.
// Returns an error if the trigger type is empty.
func ValidateTriggerType(triggerType string) error {
	if strings.TrimSpace(triggerType) == "" {
		return fmt.Errorf("trigger type cannot be empty")
	}
	return nil
}

// ValidateGA4EventName validates a GA4 event name.
// Event names must:
// - Start with a letter
// - Contain only letters, numbers, and underscores
// - Be 40 characters or less
func ValidateGA4EventName(eventName string) error {
	if eventName == "" {
		return fmt.Errorf("event name cannot be empty")
	}

	if len(eventName) > GA4EventNameMaxLength {
		return fmt.Errorf("event name exceeds maximum length of %d characters", GA4EventNameMaxLength)
	}

	// Must start with a letter
	if !isLetter(rune(eventName[0])) {
		return fmt.Errorf("event name must start with a letter, got: %s", eventName)
	}

	// Can only contain letters, numbers, and underscores
	matched, _ := regexp.MatchString(`^[a-zA-Z][a-zA-Z0-9_]*$`, eventName)
	if !matched {
		return fmt.Errorf("event name can only contain letters, numbers, and underscores, got: %s", eventName)
	}

	return nil
}

// ValidateGA4ParamName validates a GA4 event parameter name.
// Parameter names follow the same rules as event names.
func ValidateGA4ParamName(paramName string) error {
	if paramName == "" {
		return fmt.Errorf("parameter name cannot be empty")
	}

	if len(paramName) > GA4ParameterNameMaxLength {
		return fmt.Errorf("parameter name exceeds maximum length of %d characters", GA4ParameterNameMaxLength)
	}

	// Must start with a letter
	if !isLetter(rune(paramName[0])) {
		return fmt.Errorf("parameter name must start with a letter, got: %s", paramName)
	}

	// Can only contain letters, numbers, and underscores
	matched, _ := regexp.MatchString(`^[a-zA-Z][a-zA-Z0-9_]*$`, paramName)
	if !matched {
		return fmt.Errorf("parameter name can only contain letters, numbers, and underscores, got: %s", paramName)
	}

	return nil
}

// ValidateCSSSelector validates a CSS selector.
// The selector must not be empty or contain only whitespace.
func ValidateCSSSelector(selector string) error {
	if selector == "" {
		return fmt.Errorf("CSS selector cannot be empty")
	}

	trimmed := strings.TrimSpace(selector)
	if trimmed == "" {
		return fmt.Errorf("CSS selector cannot be empty or whitespace only")
	}

	if selector != trimmed {
		return fmt.Errorf("CSS selector cannot start or end with whitespace")
	}

	return nil
}

// ValidateScrollPercentages validates scroll depth percentages.
// Returns deduplicated and sorted percentages, or an error if invalid.
func ValidateScrollPercentages(percentages []int) ([]int, error) {
	if len(percentages) == 0 {
		return nil, fmt.Errorf("scroll percentages cannot be empty")
	}

	// Validate each percentage
	for i, pct := range percentages {
		if pct < 0 || pct > 100 {
			return nil, fmt.Errorf("percentage at index %d must be between 0 and 100, got: %d", i, pct)
		}
	}

	// Deduplicate
	seen := make(map[int]bool)
	unique := make([]int, 0, len(percentages))
	for _, pct := range percentages {
		if !seen[pct] {
			seen[pct] = true
			unique = append(unique, pct)
		}
	}

	// Sort
	sort.Ints(unique)

	return unique, nil
}

// ValidateTriggerIDs validates a list of trigger IDs.
func ValidateTriggerIDs(triggerIDs []string) error {
	if len(triggerIDs) == 0 {
		return fmt.Errorf("trigger IDs list cannot be empty")
	}

	for i, id := range triggerIDs {
		if strings.TrimSpace(id) == "" {
			return fmt.Errorf("trigger ID at index %d cannot be empty", i)
		}
	}

	return nil
}

// ValidatePositiveInteger validates a positive integer value.
func ValidatePositiveInteger(value int, fieldName string, minValue, maxValue int) error {
	if value < minValue {
		return fmt.Errorf("%s must be >= %d, got: %d", fieldName, minValue, value)
	}

	if maxValue > 0 && value > maxValue {
		return fmt.Errorf("%s must be <= %d, got: %d", fieldName, maxValue, value)
	}

	return nil
}

// ValidateGTMPath validates a GTM resource path format.
func ValidateGTMPath(path string) error {
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	if !strings.HasPrefix(path, "accounts/") {
		return fmt.Errorf("path must start with 'accounts/', got: %s", path)
	}

	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		return fmt.Errorf("invalid path format: %s", path)
	}

	// Validate account ID is numeric
	for _, c := range parts[1] {
		if c < '0' || c > '9' {
			return fmt.Errorf("account ID in path must be numeric: %s", path)
		}
	}

	return nil
}

// isLetter checks if a rune is a letter (a-z or A-Z).
func isLetter(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}
