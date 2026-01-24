package gtm

import (
	"strings"
	"testing"
)

func TestValidateAccountID(t *testing.T) {
	tests := []struct {
		name      string
		accountID string
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "valid account ID",
			accountID: "1234567890",
			wantErr:   false,
		},
		{
			name:      "valid long account ID",
			accountID: "12345678901234567890",
			wantErr:   false,
		},
		{
			name:      "empty account ID",
			accountID: "",
			wantErr:   true,
			errMsg:    "cannot be empty",
		},
		{
			name:      "too short account ID",
			accountID: "123456789",
			wantErr:   true,
			errMsg:    "at least 10 digits",
		},
		{
			name:      "non-numeric account ID",
			accountID: "123456789a",
			wantErr:   true,
			errMsg:    "only digits",
		},
		{
			name:      "account ID with spaces",
			accountID: "123456 7890",
			wantErr:   true,
			errMsg:    "only digits",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAccountID(tt.accountID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAccountID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" && err != nil {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateAccountID() error = %v, want error containing %q", err, tt.errMsg)
				}
			}
		})
	}
}

func TestValidateContainerID(t *testing.T) {
	tests := []struct {
		name        string
		containerID string
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "valid container ID",
			containerID: "12345678",
			wantErr:     false,
		},
		{
			name:        "single digit container ID",
			containerID: "1",
			wantErr:     false,
		},
		{
			name:        "empty container ID",
			containerID: "",
			wantErr:     true,
			errMsg:      "cannot be empty",
		},
		{
			name:        "non-numeric container ID",
			containerID: "abc123",
			wantErr:     true,
			errMsg:      "only digits",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateContainerID(tt.containerID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateContainerID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" && err != nil {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateContainerID() error = %v, want error containing %q", err, tt.errMsg)
				}
			}
		})
	}
}

func TestValidateWorkspaceID(t *testing.T) {
	tests := []struct {
		name        string
		workspaceID string
		wantErr     bool
	}{
		{"valid workspace ID", "1", false},
		{"valid larger workspace ID", "123", false},
		{"empty workspace ID", "", true},
		{"non-numeric workspace ID", "abc", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateWorkspaceID(tt.workspaceID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateWorkspaceID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateVariableID(t *testing.T) {
	tests := []struct {
		name       string
		variableID string
		wantErr    bool
	}{
		{"valid variable ID", "123", false},
		{"empty variable ID", "", true},
		{"non-numeric variable ID", "var1", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateVariableID(tt.variableID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateVariableID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateName(t *testing.T) {
	tests := []struct {
		name      string
		inputName string
		fieldName string
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "valid name",
			inputName: "My Tag",
			fieldName: "tag name",
			wantErr:   false,
		},
		{
			name:      "empty name",
			inputName: "",
			fieldName: "tag name",
			wantErr:   true,
			errMsg:    "cannot be empty",
		},
		{
			name:      "whitespace only name",
			inputName: "   ",
			fieldName: "tag name",
			wantErr:   true,
			errMsg:    "cannot be empty",
		},
		{
			name:      "name at max length",
			inputName: strings.Repeat("a", GTMNameMaxLength),
			fieldName: "tag name",
			wantErr:   false,
		},
		{
			name:      "name exceeds max length",
			inputName: strings.Repeat("a", GTMNameMaxLength+1),
			fieldName: "tag name",
			wantErr:   true,
			errMsg:    "exceeds maximum length",
		},
		{
			name:      "empty field name uses default",
			inputName: "",
			fieldName: "",
			wantErr:   true,
			errMsg:    "name cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateName(tt.inputName, tt.fieldName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" && err != nil {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateName() error = %v, want error containing %q", err, tt.errMsg)
				}
			}
		})
	}
}

func TestValidateNotes(t *testing.T) {
	tests := []struct {
		name    string
		notes   string
		wantErr bool
	}{
		{"empty notes", "", false},
		{"valid notes", "Some description", false},
		{"notes at max length", strings.Repeat("a", GTMNotesMaxLength), false},
		{"notes exceed max length", strings.Repeat("a", GTMNotesMaxLength+1), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNotes(tt.notes)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateNotes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGA4EventName(t *testing.T) {
	tests := []struct {
		name      string
		eventName string
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "valid event name",
			eventName: "purchase",
			wantErr:   false,
		},
		{
			name:      "valid event name with underscore",
			eventName: "add_to_cart",
			wantErr:   false,
		},
		{
			name:      "valid event name with numbers",
			eventName: "event123",
			wantErr:   false,
		},
		{
			name:      "empty event name",
			eventName: "",
			wantErr:   true,
			errMsg:    "cannot be empty",
		},
		{
			name:      "event name starting with number",
			eventName: "123event",
			wantErr:   true,
			errMsg:    "must start with a letter",
		},
		{
			name:      "event name with spaces",
			eventName: "my event",
			wantErr:   true,
			errMsg:    "only contain letters, numbers, and underscores",
		},
		{
			name:      "event name with hyphen",
			eventName: "my-event",
			wantErr:   true,
			errMsg:    "only contain letters, numbers, and underscores",
		},
		{
			name:      "event name too long",
			eventName: strings.Repeat("a", GA4EventNameMaxLength+1),
			wantErr:   true,
			errMsg:    "exceeds maximum length",
		},
		{
			name:      "event name at max length",
			eventName: "a" + strings.Repeat("b", GA4EventNameMaxLength-1),
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGA4EventName(tt.eventName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGA4EventName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" && err != nil {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateGA4EventName() error = %v, want error containing %q", err, tt.errMsg)
				}
			}
		})
	}
}

func TestValidateGA4ParamName(t *testing.T) {
	tests := []struct {
		name      string
		paramName string
		wantErr   bool
	}{
		{"valid param name", "item_id", false},
		{"valid param name", "value", false},
		{"empty param name", "", true},
		{"param name starting with number", "1param", true},
		{"param name with special chars", "param@name", true},
		{"param name too long", strings.Repeat("a", GA4ParameterNameMaxLength+1), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGA4ParamName(tt.paramName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGA4ParamName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateCSSSelector(t *testing.T) {
	tests := []struct {
		name     string
		selector string
		wantErr  bool
		errMsg   string
	}{
		{"valid ID selector", "#myId", false, ""},
		{"valid class selector", ".myClass", false, ""},
		{"valid tag selector", "div", false, ""},
		{"valid complex selector", "div.class > span#id", false, ""},
		{"empty selector", "", true, "cannot be empty"},
		{"whitespace only selector", "   ", true, "cannot be empty"},
		{"selector with leading space", " .class", true, "cannot start or end with whitespace"},
		{"selector with trailing space", ".class ", true, "cannot start or end with whitespace"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCSSSelector(tt.selector)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCSSSelector() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" && err != nil {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateCSSSelector() error = %v, want error containing %q", err, tt.errMsg)
				}
			}
		})
	}
}

func TestValidateScrollPercentages(t *testing.T) {
	tests := []struct {
		name        string
		percentages []int
		want        []int
		wantErr     bool
	}{
		{
			name:        "valid percentages",
			percentages: []int{25, 50, 75, 100},
			want:        []int{25, 50, 75, 100},
			wantErr:     false,
		},
		{
			name:        "unsorted percentages",
			percentages: []int{100, 25, 75, 50},
			want:        []int{25, 50, 75, 100},
			wantErr:     false,
		},
		{
			name:        "duplicate percentages",
			percentages: []int{25, 50, 25, 75},
			want:        []int{25, 50, 75},
			wantErr:     false,
		},
		{
			name:        "boundary values",
			percentages: []int{0, 100},
			want:        []int{0, 100},
			wantErr:     false,
		},
		{
			name:        "empty percentages",
			percentages: []int{},
			want:        nil,
			wantErr:     true,
		},
		{
			name:        "negative percentage",
			percentages: []int{25, -10, 50},
			want:        nil,
			wantErr:     true,
		},
		{
			name:        "percentage over 100",
			percentages: []int{25, 50, 150},
			want:        nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateScrollPercentages(tt.percentages)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateScrollPercentages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("ValidateScrollPercentages() = %v, want %v", got, tt.want)
					return
				}
				for i, v := range got {
					if v != tt.want[i] {
						t.Errorf("ValidateScrollPercentages() = %v, want %v", got, tt.want)
						return
					}
				}
			}
		})
	}
}

func TestValidateTriggerIDs(t *testing.T) {
	tests := []struct {
		name       string
		triggerIDs []string
		wantErr    bool
	}{
		{"valid trigger IDs", []string{"1", "2", "3"}, false},
		{"single trigger ID", []string{"1"}, false},
		{"empty list", []string{}, true},
		{"empty ID in list", []string{"1", "", "3"}, true},
		{"whitespace ID in list", []string{"1", "   ", "3"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTriggerIDs(tt.triggerIDs)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTriggerIDs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePositiveInteger(t *testing.T) {
	tests := []struct {
		name     string
		value    int
		field    string
		minValue int
		maxValue int
		wantErr  bool
	}{
		{"value in range", 50, "count", 0, 100, false},
		{"value at min", 0, "count", 0, 100, false},
		{"value at max", 100, "count", 0, 100, false},
		{"value below min", -1, "count", 0, 100, true},
		{"value above max", 101, "count", 0, 100, true},
		{"no max constraint", 1000, "count", 0, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePositiveInteger(tt.value, tt.field, tt.minValue, tt.maxValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePositiveInteger() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateGTMPath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid account path",
			path:    "accounts/1234567890",
			wantErr: false,
		},
		{
			name:    "valid container path",
			path:    "accounts/1234567890/containers/12345678",
			wantErr: false,
		},
		{
			name:    "valid workspace path",
			path:    "accounts/1234567890/containers/12345678/workspaces/1",
			wantErr: false,
		},
		{
			name:    "empty path",
			path:    "",
			wantErr: true,
			errMsg:  "cannot be empty",
		},
		{
			name:    "path not starting with accounts",
			path:    "containers/12345678",
			wantErr: true,
			errMsg:  "must start with 'accounts/'",
		},
		{
			name:    "non-numeric account ID in path",
			path:    "accounts/abc123",
			wantErr: true,
			errMsg:  "must be numeric",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGTMPath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGTMPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" && err != nil {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateGTMPath() error = %v, want error containing %q", err, tt.errMsg)
				}
			}
		})
	}
}

func TestValidateTagInput(t *testing.T) {
	tests := []struct {
		name              string
		tagName           string
		tagType           string
		firingTriggerIDs  []string
		wantErr           bool
	}{
		{"valid input", "My Tag", "html", []string{"1"}, false},
		{"empty name", "", "html", []string{"1"}, true},
		{"long name", strings.Repeat("a", 257), "html", []string{"1"}, true},
		{"empty type", "My Tag", "", []string{"1"}, true},
		{"empty trigger list", "My Tag", "html", []string{}, true},
		{"empty trigger ID", "My Tag", "html", []string{"1", ""}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTagInput(tt.tagName, tt.tagType, tt.firingTriggerIDs)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTagInput() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateTriggerInput(t *testing.T) {
	tests := []struct {
		name        string
		triggerName string
		triggerType string
		wantErr     bool
	}{
		{"valid input", "My Trigger", "PAGEVIEW", false},
		{"empty name", "", "PAGEVIEW", true},
		{"long name", strings.Repeat("a", 257), "PAGEVIEW", true},
		{"empty type", "My Trigger", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTriggerInput(tt.triggerName, tt.triggerType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTriggerInput() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateVariableInput(t *testing.T) {
	tests := []struct {
		name     string
		varName  string
		varType  string
		wantErr  bool
	}{
		{"valid input", "My Variable", "c", false},
		{"empty name", "", "c", true},
		{"long name", strings.Repeat("a", 257), "c", true},
		{"empty type", "My Variable", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateVariableInput(tt.varName, tt.varType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateVariableInput() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateWorkspacePath(t *testing.T) {
	tests := []struct {
		name        string
		accountID   string
		containerID string
		workspaceID string
		wantErr     bool
	}{
		{"valid input", "1234567890", "12345678", "1", false},
		{"empty account ID", "", "12345678", "1", true},
		{"empty container ID", "1234567890", "", "1", true},
		{"empty workspace ID", "1234567890", "12345678", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateWorkspacePath(tt.accountID, tt.containerID, tt.workspaceID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateWorkspacePath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateContainerPath(t *testing.T) {
	tests := []struct {
		name        string
		accountID   string
		containerID string
		wantErr     bool
	}{
		{"valid input", "1234567890", "12345678", false},
		{"empty account ID", "", "12345678", true},
		{"empty container ID", "1234567890", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateContainerPath(tt.accountID, tt.containerID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateContainerPath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateTriggerType(t *testing.T) {
	tests := []struct {
		name        string
		triggerType string
		wantErr     bool
	}{
		{"valid type", "PAGEVIEW", false},
		{"empty type", "", true},
		{"whitespace type", "   ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTriggerType(tt.triggerType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTriggerType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
