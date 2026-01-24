package gtm

import (
	"testing"
)

func TestBuildTemplateParameter(t *testing.T) {
	param := BuildTemplateParameter("testKey", "testValue")

	if param.Type != ParameterTypeTemplate {
		t.Errorf("BuildTemplateParameter() Type = %v, want %v", param.Type, ParameterTypeTemplate)
	}
	if param.Key != "testKey" {
		t.Errorf("BuildTemplateParameter() Key = %v, want %v", param.Key, "testKey")
	}
	if param.Value != "testValue" {
		t.Errorf("BuildTemplateParameter() Value = %v, want %v", param.Value, "testValue")
	}
}

func TestBuildBooleanParameter(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    bool
		expected string
	}{
		{"true value", "enabled", true, "true"},
		{"false value", "disabled", false, "false"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			param := BuildBooleanParameter(tt.key, tt.value)

			if param.Type != ParameterTypeBoolean {
				t.Errorf("BuildBooleanParameter() Type = %v, want %v", param.Type, ParameterTypeBoolean)
			}
			if param.Key != tt.key {
				t.Errorf("BuildBooleanParameter() Key = %v, want %v", param.Key, tt.key)
			}
			if param.Value != tt.expected {
				t.Errorf("BuildBooleanParameter() Value = %v, want %v", param.Value, tt.expected)
			}
		})
	}
}

func TestBuildIntegerParameter(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    int
		expected string
	}{
		{"positive value", "count", 42, "42"},
		{"zero value", "index", 0, "0"},
		{"negative value", "offset", -10, "-10"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			param := BuildIntegerParameter(tt.key, tt.value)

			if param.Type != ParameterTypeInteger {
				t.Errorf("BuildIntegerParameter() Type = %v, want %v", param.Type, ParameterTypeInteger)
			}
			if param.Key != tt.key {
				t.Errorf("BuildIntegerParameter() Key = %v, want %v", param.Key, tt.key)
			}
			if param.Value != tt.expected {
				t.Errorf("BuildIntegerParameter() Value = %v, want %v", param.Value, tt.expected)
			}
		})
	}
}

func TestBuildListParameter(t *testing.T) {
	items := []Parameter{
		BuildTemplateParameter("item1", "value1"),
		BuildTemplateParameter("item2", "value2"),
	}
	param := BuildListParameter("myList", items)

	if param.Type != ParameterTypeList {
		t.Errorf("BuildListParameter() Type = %v, want %v", param.Type, ParameterTypeList)
	}
	if param.Key != "myList" {
		t.Errorf("BuildListParameter() Key = %v, want %v", param.Key, "myList")
	}
	if len(param.List) != 2 {
		t.Errorf("BuildListParameter() List length = %v, want %v", len(param.List), 2)
	}
}

func TestBuildMapParameter(t *testing.T) {
	pairs := []Parameter{
		BuildTemplateParameter("key1", "value1"),
		BuildTemplateParameter("key2", "value2"),
	}
	param := BuildMapParameter(pairs)

	if param.Type != ParameterTypeMap {
		t.Errorf("BuildMapParameter() Type = %v, want %v", param.Type, ParameterTypeMap)
	}
	if len(param.Map) != 2 {
		t.Errorf("BuildMapParameter() Map length = %v, want %v", len(param.Map), 2)
	}
}

func TestBuildTagReferenceParameter(t *testing.T) {
	param := BuildTagReferenceParameter("measurementId", "GA4 Config")

	if param.Type != ParameterTypeTagReference {
		t.Errorf("BuildTagReferenceParameter() Type = %v, want %v", param.Type, ParameterTypeTagReference)
	}
	if param.Key != "measurementId" {
		t.Errorf("BuildTagReferenceParameter() Key = %v, want %v", param.Key, "measurementId")
	}
	if param.Value != "GA4 Config" {
		t.Errorf("BuildTagReferenceParameter() Value = %v, want %v", param.Value, "GA4 Config")
	}
}

func TestBuildTriggerReferenceParameter(t *testing.T) {
	param := BuildTriggerReferenceParameter("123")

	if param.Type != ParameterTypeTriggerReference {
		t.Errorf("BuildTriggerReferenceParameter() Type = %v, want %v", param.Type, ParameterTypeTriggerReference)
	}
	if param.Value != "123" {
		t.Errorf("BuildTriggerReferenceParameter() Value = %v, want %v", param.Value, "123")
	}
}

func TestBuildEventParameter(t *testing.T) {
	param := BuildEventParameter("item_id", "SKU123")

	if param.Type != ParameterTypeMap {
		t.Errorf("BuildEventParameter() Type = %v, want %v", param.Type, ParameterTypeMap)
	}
	if len(param.Map) != 2 {
		t.Errorf("BuildEventParameter() Map length = %v, want %v", len(param.Map), 2)
	}

	// Check that name and value are present
	var hasName, hasValue bool
	for _, p := range param.Map {
		if p.Key == "name" && p.Value == "item_id" {
			hasName = true
		}
		if p.Key == "value" && p.Value == "SKU123" {
			hasValue = true
		}
	}
	if !hasName {
		t.Error("BuildEventParameter() missing name parameter")
	}
	if !hasValue {
		t.Error("BuildEventParameter() missing value parameter")
	}
}

func TestBuildEventParametersList(t *testing.T) {
	params := map[string]string{
		"item_id": "SKU123",
		"price":   "29.99",
	}
	result := BuildEventParametersList(params)

	if len(result) != 2 {
		t.Errorf("BuildEventParametersList() length = %v, want %v", len(result), 2)
	}

	// Each item should be a map type
	for _, item := range result {
		if item.Type != ParameterTypeMap {
			t.Errorf("BuildEventParametersList() item Type = %v, want %v", item.Type, ParameterTypeMap)
		}
	}
}

func TestBuildScrollPercentageList(t *testing.T) {
	percentages := []int{25, 50, 75, 100}
	param := BuildScrollPercentageList(percentages)

	if param.Type != ParameterTypeList {
		t.Errorf("BuildScrollPercentageList() Type = %v, want %v", param.Type, ParameterTypeList)
	}
	if param.Key != "verticalScrollPercentageList" {
		t.Errorf("BuildScrollPercentageList() Key = %v, want %v", param.Key, "verticalScrollPercentageList")
	}
	if len(param.List) != 4 {
		t.Errorf("BuildScrollPercentageList() List length = %v, want %v", len(param.List), 4)
	}

	// Check values
	expectedValues := []string{"25", "50", "75", "100"}
	for i, item := range param.List {
		if item.Value != expectedValues[i] {
			t.Errorf("BuildScrollPercentageList() item[%d].Value = %v, want %v", i, item.Value, expectedValues[i])
		}
		if item.Type != ParameterTypeTemplate {
			t.Errorf("BuildScrollPercentageList() item[%d].Type = %v, want %v", i, item.Type, ParameterTypeTemplate)
		}
	}
}

func TestBuildCustomEventFilter(t *testing.T) {
	tests := []struct {
		name          string
		eventName     string
		matchType     string
		wantMatchType string
	}{
		{"default match type", "purchase", "", "EQUALS"},
		{"explicit equals", "purchase", "EQUALS", "EQUALS"},
		{"contains match", "purchase", "CONTAINS", "CONTAINS"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conditions := BuildCustomEventFilter(tt.eventName, tt.matchType)

			if len(conditions) != 1 {
				t.Errorf("BuildCustomEventFilter() length = %v, want %v", len(conditions), 1)
				return
			}
			if conditions[0].Type != tt.wantMatchType {
				t.Errorf("BuildCustomEventFilter() Type = %v, want %v", conditions[0].Type, tt.wantMatchType)
			}
			if len(conditions[0].Parameter) != 2 {
				t.Errorf("BuildCustomEventFilter() Parameter length = %v, want %v", len(conditions[0].Parameter), 2)
			}
		})
	}
}

func TestBuildURLFilter(t *testing.T) {
	conditions := BuildURLFilter("{{Page URL}}", "CONTAINS", "/products")

	if len(conditions) != 1 {
		t.Errorf("BuildURLFilter() length = %v, want %v", len(conditions), 1)
		return
	}
	if conditions[0].Type != "CONTAINS" {
		t.Errorf("BuildURLFilter() Type = %v, want %v", conditions[0].Type, "CONTAINS")
	}
	if len(conditions[0].Parameter) != 2 {
		t.Errorf("BuildURLFilter() Parameter length = %v, want %v", len(conditions[0].Parameter), 2)
	}
}

func TestBuildPageViewFilter(t *testing.T) {
	conditions := BuildPageViewFilter("CONTAINS", "/checkout")

	if len(conditions) != 1 {
		t.Errorf("BuildPageViewFilter() length = %v, want %v", len(conditions), 1)
		return
	}
	// Should use {{Page URL}} as the variable
	if conditions[0].Parameter[0].Value != "{{Page URL}}" {
		t.Errorf("BuildPageViewFilter() variable = %v, want %v", conditions[0].Parameter[0].Value, "{{Page URL}}")
	}
}

func TestBuildClickFilter(t *testing.T) {
	tests := []struct {
		name          string
		matchType     string
		pattern       string
		clickProperty string
		wantVariable  string
	}{
		{"default click URL", "CONTAINS", "/link", "", "{{Click URL}}"},
		{"explicit click URL", "CONTAINS", "/link", "{{Click URL}}", "{{Click URL}}"},
		{"click text", "EQUALS", "Buy Now", "{{Click Text}}", "{{Click Text}}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conditions := BuildClickFilter(tt.matchType, tt.pattern, tt.clickProperty)

			if len(conditions) != 1 {
				t.Errorf("BuildClickFilter() length = %v, want %v", len(conditions), 1)
				return
			}
			if conditions[0].Parameter[0].Value != tt.wantVariable {
				t.Errorf("BuildClickFilter() variable = %v, want %v", conditions[0].Parameter[0].Value, tt.wantVariable)
			}
		})
	}
}

func TestBuildVersionPath(t *testing.T) {
	path := BuildVersionPath("1234567890", "12345678", "5")
	expected := "accounts/1234567890/containers/12345678/versions/5"

	if path != expected {
		t.Errorf("BuildVersionPath() = %v, want %v", path, expected)
	}
}

func TestBuildWorkspacePath(t *testing.T) {
	path := BuildWorkspacePath("1234567890", "12345678", "1")
	expected := "accounts/1234567890/containers/12345678/workspaces/1"

	if path != expected {
		t.Errorf("BuildWorkspacePath() = %v, want %v", path, expected)
	}
}

func TestBuildContainerPath(t *testing.T) {
	path := BuildContainerPath("1234567890", "12345678")
	expected := "accounts/1234567890/containers/12345678"

	if path != expected {
		t.Errorf("BuildContainerPath() = %v, want %v", path, expected)
	}
}

func TestExtractIDFromPath(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		resourceType string
		wantID       string
		wantErr      bool
	}{
		{
			name:         "extract account ID",
			path:         "accounts/1234567890/containers/12345678",
			resourceType: "account",
			wantID:       "1234567890",
			wantErr:      false,
		},
		{
			name:         "extract container ID",
			path:         "accounts/1234567890/containers/12345678",
			resourceType: "container",
			wantID:       "12345678",
			wantErr:      false,
		},
		{
			name:         "extract workspace ID",
			path:         "accounts/1234567890/containers/12345678/workspaces/1",
			resourceType: "workspace",
			wantID:       "1",
			wantErr:      false,
		},
		{
			name:         "extract tag ID",
			path:         "accounts/1234567890/containers/12345678/workspaces/1/tags/456",
			resourceType: "tag",
			wantID:       "456",
			wantErr:      false,
		},
		{
			name:         "resource not found",
			path:         "accounts/1234567890/containers/12345678",
			resourceType: "workspace",
			wantID:       "",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractIDFromPath(tt.path, tt.resourceType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractIDFromPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantID {
				t.Errorf("ExtractIDFromPath() = %v, want %v", got, tt.wantID)
			}
		})
	}
}

func TestParseWorkspacePath(t *testing.T) {
	tests := []struct {
		name            string
		path            string
		wantAccountID   string
		wantContainerID string
		wantWorkspaceID string
		wantErr         bool
	}{
		{
			name:            "valid path",
			path:            "accounts/1234567890/containers/12345678/workspaces/1",
			wantAccountID:   "1234567890",
			wantContainerID: "12345678",
			wantWorkspaceID: "1",
			wantErr:         false,
		},
		{
			name:    "invalid format - too few parts",
			path:    "accounts/1234567890/containers",
			wantErr: true,
		},
		{
			name:    "invalid format - wrong structure",
			path:    "containers/1234567890/accounts/12345678/workspaces/1",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accountID, containerID, workspaceID, err := ParseWorkspacePath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseWorkspacePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if accountID != tt.wantAccountID {
					t.Errorf("ParseWorkspacePath() accountID = %v, want %v", accountID, tt.wantAccountID)
				}
				if containerID != tt.wantContainerID {
					t.Errorf("ParseWorkspacePath() containerID = %v, want %v", containerID, tt.wantContainerID)
				}
				if workspaceID != tt.wantWorkspaceID {
					t.Errorf("ParseWorkspacePath() workspaceID = %v, want %v", workspaceID, tt.wantWorkspaceID)
				}
			}
		})
	}
}

func TestParseContainerPath(t *testing.T) {
	tests := []struct {
		name            string
		path            string
		wantAccountID   string
		wantContainerID string
		wantErr         bool
	}{
		{
			name:            "valid path",
			path:            "accounts/1234567890/containers/12345678",
			wantAccountID:   "1234567890",
			wantContainerID: "12345678",
			wantErr:         false,
		},
		{
			name:            "valid path with extra parts",
			path:            "accounts/1234567890/containers/12345678/workspaces/1",
			wantAccountID:   "1234567890",
			wantContainerID: "12345678",
			wantErr:         false,
		},
		{
			name:    "invalid format - too few parts",
			path:    "accounts/1234567890",
			wantErr: true,
		},
		{
			name:    "invalid format - wrong structure",
			path:    "containers/1234567890/accounts/12345678",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accountID, containerID, err := ParseContainerPath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseContainerPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if accountID != tt.wantAccountID {
					t.Errorf("ParseContainerPath() accountID = %v, want %v", accountID, tt.wantAccountID)
				}
				if containerID != tt.wantContainerID {
					t.Errorf("ParseContainerPath() containerID = %v, want %v", containerID, tt.wantContainerID)
				}
			}
		})
	}
}

func TestBuildGA4ConfigTag(t *testing.T) {
	params := BuildGA4ConfigTag("G-XXXXXXX", true)

	if len(params) != 2 {
		t.Errorf("BuildGA4ConfigTag() length = %v, want %v", len(params), 2)
	}

	// Check measurement ID
	var hasMeasurementID, hasSendPageView bool
	for _, p := range params {
		if p.Key == "measurementId" && p.Value == "G-XXXXXXX" {
			hasMeasurementID = true
		}
		if p.Key == "sendPageView" && p.Value == "true" {
			hasSendPageView = true
		}
	}
	if !hasMeasurementID {
		t.Error("BuildGA4ConfigTag() missing measurementId parameter")
	}
	if !hasSendPageView {
		t.Error("BuildGA4ConfigTag() missing sendPageView parameter")
	}
}

func TestBuildGA4EventTag(t *testing.T) {
	tests := []struct {
		name          string
		configTagName string
		eventName     string
		params        map[string]string
		sendEcommerce bool
		wantMinParams int
	}{
		{"basic event", "GA4 Config", "purchase", nil, false, 2},
		{"event with params", "GA4 Config", "purchase", map[string]string{"item_id": "SKU123"}, false, 3},
		{"event with ecommerce", "GA4 Config", "purchase", nil, true, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := BuildGA4EventTag(tt.configTagName, tt.eventName, tt.params, tt.sendEcommerce)

			if len(params) < tt.wantMinParams {
				t.Errorf("BuildGA4EventTag() length = %v, want at least %v", len(params), tt.wantMinParams)
			}

			// Check required parameters exist
			var hasConfigRef, hasEventName bool
			for _, p := range params {
				if p.Key == "measurementId" && p.Type == ParameterTypeTagReference {
					hasConfigRef = true
				}
				if p.Key == "eventName" && p.Value == tt.eventName {
					hasEventName = true
				}
			}
			if !hasConfigRef {
				t.Error("BuildGA4EventTag() missing measurementId tag reference")
			}
			if !hasEventName {
				t.Error("BuildGA4EventTag() missing eventName parameter")
			}
		})
	}
}

func TestMergeParameters(t *testing.T) {
	list1 := []Parameter{
		BuildTemplateParameter("key1", "value1"),
		BuildTemplateParameter("key2", "value2"),
	}
	list2 := []Parameter{
		BuildTemplateParameter("key2", "overridden"),
		BuildTemplateParameter("key3", "value3"),
	}

	merged := MergeParameters(list1, list2)

	// Should have 3 unique keys
	if len(merged) != 3 {
		t.Errorf("MergeParameters() length = %v, want %v", len(merged), 3)
	}

	// Check that key2 was overridden
	for _, p := range merged {
		if p.Key == "key2" && p.Value != "overridden" {
			t.Errorf("MergeParameters() key2 Value = %v, want %v", p.Value, "overridden")
		}
	}
}

func TestIsVariableReference(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"valid reference", "{{Page URL}}", true},
		{"valid reference with spaces", "{{ Click URL }}", true},
		{"valid reference", "{{myVar}}", true},
		{"not a reference - no braces", "Page URL", false},
		{"not a reference - partial braces", "{Page URL}", false},
		{"not a reference - empty", "", false},
		{"not a reference - only open", "{{test", false},
		{"not a reference - only close", "test}}", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsVariableReference(tt.input)
			if got != tt.want {
				t.Errorf("IsVariableReference(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestWrapVariableReference(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"wrap plain name", "Page URL", "{{Page URL}}"},
		{"already wrapped", "{{Page URL}}", "{{Page URL}}"},
		{"wrap simple name", "myVar", "{{myVar}}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WrapVariableReference(tt.input)
			if got != tt.want {
				t.Errorf("WrapVariableReference(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestUnwrapVariableReference(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"unwrap reference", "{{Page URL}}", "Page URL"},
		{"plain name unchanged", "Page URL", "Page URL"},
		{"unwrap simple", "{{myVar}}", "myVar"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UnwrapVariableReference(tt.input)
			if got != tt.want {
				t.Errorf("UnwrapVariableReference(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
