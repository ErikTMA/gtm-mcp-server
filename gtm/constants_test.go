package gtm

import (
	"testing"
)

func TestIsValidTriggerType(t *testing.T) {
	tests := []struct {
		name        string
		triggerType string
		want        bool
	}{
		{"valid PAGEVIEW", TriggerTypePageview, true},
		{"valid CUSTOM_EVENT", TriggerTypeCustomEvent, true},
		{"valid SCROLL_DEPTH", TriggerTypeScrollDepth, true},
		{"valid DOM_READY", TriggerTypeDOMReady, true},
		{"valid WINDOW_LOADED", TriggerTypeWindowLoaded, true},
		{"valid YOUTUBE_VIDEO", TriggerTypeYouTubeVideo, true},
		{"valid server trigger", TriggerTypeServerPageview, true},
		{"valid firebase trigger", TriggerTypeFirebaseAppException, true},
		{"valid AMP trigger", TriggerTypeAMPClick, true},
		{"invalid type", "INVALID_TYPE", false},
		{"empty type", "", false},
		{"lowercase pageview", "pageview", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidTriggerType(tt.triggerType)
			if got != tt.want {
				t.Errorf("IsValidTriggerType(%q) = %v, want %v", tt.triggerType, got, tt.want)
			}
		})
	}
}

func TestIsValidTagType(t *testing.T) {
	tests := []struct {
		name    string
		tagType string
		want    bool
	}{
		{"valid GA4 config", TagTypeGA4Config, true},
		{"valid GA4 event", TagTypeGA4Event, true},
		{"valid custom HTML", TagTypeCustomHTML, true},
		{"valid UA", TagTypeUA, true},
		{"valid Google Ads conversion", TagTypeGoogleAdsConversion, true},
		{"invalid type", "invalid_tag", false},
		{"empty type", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidTagType(tt.tagType)
			if got != tt.want {
				t.Errorf("IsValidTagType(%q) = %v, want %v", tt.tagType, got, tt.want)
			}
		})
	}
}

func TestIsValidVariableType(t *testing.T) {
	tests := []struct {
		name         string
		variableType string
		want         bool
	}{
		{"valid constant", VariableTypeConstant, true},
		{"valid custom JS", VariableTypeCustomJavaScript, true},
		{"valid data layer", VariableTypeDataLayer, true},
		{"valid URL", VariableTypeURL, true},
		{"valid cookie", VariableTypeFirstPartyCookie, true},
		{"valid GA4 event settings", VariableTypeGA4EventSettings, true},
		{"invalid type", "invalid_var", false},
		{"empty type", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidVariableType(tt.variableType)
			if got != tt.want {
				t.Errorf("IsValidVariableType(%q) = %v, want %v", tt.variableType, got, tt.want)
			}
		})
	}
}

func TestIsValidFilterType(t *testing.T) {
	tests := []struct {
		name       string
		filterType string
		want       bool
	}{
		{"valid EQUALS", FilterTypeEquals, true},
		{"valid CONTAINS", FilterTypeContains, true},
		{"valid STARTS_WITH", FilterTypeStartsWith, true},
		{"valid ENDS_WITH", FilterTypeEndsWith, true},
		{"valid MATCHES_REGEX", FilterTypeMatchesRegex, true},
		{"valid CSS_SELECTOR", FilterTypeCSSSelector, true},
		{"invalid type", "INVALID", false},
		{"empty type", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidFilterType(tt.filterType)
			if got != tt.want {
				t.Errorf("IsValidFilterType(%q) = %v, want %v", tt.filterType, got, tt.want)
			}
		})
	}
}

func TestIsBuiltInVariable(t *testing.T) {
	tests := []struct {
		name string
		varName string
		want bool
	}{
		{"PAGE_URL", "PAGE_URL", true},
		{"CLICK_URL", "CLICK_URL", true},
		{"EVENT", "EVENT", true},
		{"REFERRER", "REFERRER", true},
		{"VIDEO_PROVIDER", "VIDEO_PROVIDER", true},
		{"SCROLL_DEPTH_THRESHOLD", "SCROLL_DEPTH_THRESHOLD", true},
		{"invalid variable", "NOT_A_BUILTIN", false},
		{"lowercase", "page_url", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsBuiltInVariable(tt.varName)
			if got != tt.want {
				t.Errorf("IsBuiltInVariable(%q) = %v, want %v", tt.varName, got, tt.want)
			}
		})
	}
}

func TestBuiltInVariablesCount(t *testing.T) {
	// Ensure we have a reasonable number of built-in variables
	if len(BuiltInVariables) < 30 {
		t.Errorf("BuiltInVariables count = %d, expected at least 30", len(BuiltInVariables))
	}
}

func TestDefaultScrollPercentages(t *testing.T) {
	expected := []int{10, 25, 50, 75, 90, 100}

	if len(DefaultScrollPercentages) != len(expected) {
		t.Errorf("DefaultScrollPercentages count = %d, expected %d", len(DefaultScrollPercentages), len(expected))
		return
	}

	for i, v := range DefaultScrollPercentages {
		if v != expected[i] {
			t.Errorf("DefaultScrollPercentages[%d] = %d, expected %d", i, v, expected[i])
		}
	}
}

func TestScopeConstants(t *testing.T) {
	// Verify scope URLs are properly formatted
	scopes := []string{
		ScopeDeleteContainers,
		ScopeEditContainers,
		ScopeEditContainerVersions,
		ScopeManageAccounts,
		ScopeManageUsers,
		ScopePublish,
		ScopeReadonly,
	}

	for _, scope := range scopes {
		if scope == "" {
			t.Error("Found empty scope constant")
		}
		if len(scope) < 10 {
			t.Errorf("Scope %q seems invalid (too short)", scope)
		}
	}

	// Verify AllScopes contains all individual scopes
	if len(AllScopes) != 7 {
		t.Errorf("AllScopes count = %d, expected 7", len(AllScopes))
	}
}

func TestTriggerTypeConstants(t *testing.T) {
	// Web trigger types
	webTriggers := []string{
		TriggerTypePageview,
		TriggerTypeDOMReady,
		TriggerTypeWindowLoaded,
		TriggerTypeCustomEvent,
		TriggerTypeClick,
		TriggerTypeLinkClick,
		TriggerTypeFormSubmission,
		TriggerTypeScrollDepth,
		TriggerTypeYouTubeVideo,
	}

	for _, trigger := range webTriggers {
		if !IsValidTriggerType(trigger) {
			t.Errorf("Web trigger %q not in ValidTriggerTypes", trigger)
		}
	}

	// Server trigger types
	serverTriggers := []string{
		TriggerTypeServerPageview,
		TriggerTypeAlways,
		TriggerTypeConsentInit,
		TriggerTypeInit,
	}

	for _, trigger := range serverTriggers {
		if !IsValidTriggerType(trigger) {
			t.Errorf("Server trigger %q not in ValidTriggerTypes", trigger)
		}
	}
}

func TestTagTypeConstants(t *testing.T) {
	// GA4 tag types
	if TagTypeGA4Config != "gaawc" {
		t.Errorf("TagTypeGA4Config = %q, expected %q", TagTypeGA4Config, "gaawc")
	}
	if TagTypeGA4Event != "gaawe" {
		t.Errorf("TagTypeGA4Event = %q, expected %q", TagTypeGA4Event, "gaawe")
	}

	// Custom tag types
	if TagTypeCustomHTML != "html" {
		t.Errorf("TagTypeCustomHTML = %q, expected %q", TagTypeCustomHTML, "html")
	}
}

func TestVariableTypeConstants(t *testing.T) {
	// Common variable types
	if VariableTypeConstant != "c" {
		t.Errorf("VariableTypeConstant = %q, expected %q", VariableTypeConstant, "c")
	}
	if VariableTypeCustomJavaScript != "jsm" {
		t.Errorf("VariableTypeCustomJavaScript = %q, expected %q", VariableTypeCustomJavaScript, "jsm")
	}
	if VariableTypeDataLayer != "v" {
		t.Errorf("VariableTypeDataLayer = %q, expected %q", VariableTypeDataLayer, "v")
	}
}

func TestFilterTypeConstants(t *testing.T) {
	// All filter types should be uppercase
	filterTypes := []string{
		FilterTypeEquals,
		FilterTypeContains,
		FilterTypeStartsWith,
		FilterTypeEndsWith,
		FilterTypeMatchesRegex,
		FilterTypeGreaterThan,
		FilterTypeGreaterOrEquals,
		FilterTypeLessThan,
		FilterTypeLessOrEquals,
		FilterTypeCSSSelector,
		FilterTypeMatchRegex,
	}

	for _, ft := range filterTypes {
		if ft == "" {
			t.Error("Found empty filter type constant")
		}
		// Check they're uppercase (GTM convention)
		if ft != "" && ft[0] >= 'a' && ft[0] <= 'z' {
			t.Errorf("Filter type %q should be uppercase", ft)
		}
	}
}

func TestDefaultWorkspaceID(t *testing.T) {
	if DefaultWorkspaceID != "1" {
		t.Errorf("DefaultWorkspaceID = %q, expected %q", DefaultWorkspaceID, "1")
	}
}

func TestGA4ParameterValueMaxLength(t *testing.T) {
	if GA4ParameterValueMaxLength != 100 {
		t.Errorf("GA4ParameterValueMaxLength = %d, expected %d", GA4ParameterValueMaxLength, 100)
	}
}

func TestTagFiringOptions(t *testing.T) {
	options := []string{
		TagFiringUnlimited,
		TagFiringOncePerEvent,
		TagFiringOncePerLoad,
	}

	for _, opt := range options {
		if opt == "" {
			t.Error("Found empty tag firing option")
		}
	}

	if TagFiringUnlimited != "UNLIMITED" {
		t.Errorf("TagFiringUnlimited = %q, expected %q", TagFiringUnlimited, "UNLIMITED")
	}
}
