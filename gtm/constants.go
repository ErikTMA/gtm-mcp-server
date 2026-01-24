// Package gtm provides constants and enumerations for GTM API.
//
// This module defines type-safe constants for GTM trigger types, tag types,
// variable types, filter operators, and parameter types based on the official
// Google Tag Manager API specification.
package gtm

// Trigger types for GTM API.
// Reference: https://developers.google.com/tag-platform/tag-manager/api/v2/reference/accounts/containers/workspaces/triggers
const (
	// Web container triggers
	TriggerTypePageview          = "PAGEVIEW"
	TriggerTypeDOMReady          = "DOM_READY"
	TriggerTypeWindowLoaded      = "WINDOW_LOADED"
	TriggerTypeCustomEvent       = "CUSTOM_EVENT"
	TriggerTypeTriggerGroup      = "TRIGGER_GROUP"
	TriggerTypeFormSubmission    = "FORM_SUBMISSION"
	TriggerTypeClick             = "CLICK"
	TriggerTypeLinkClick         = "LINK_CLICK"
	TriggerTypeJSError           = "JS_ERROR"
	TriggerTypeHistoryChange     = "HISTORY_CHANGE"
	TriggerTypeTimer             = "TIMER"
	TriggerTypeScrollDepth       = "SCROLL_DEPTH"
	TriggerTypeElementVisibility = "ELEMENT_VISIBILITY"
	TriggerTypeYouTubeVideo      = "YOUTUBE_VIDEO"

	// Server container triggers
	TriggerTypeServerPageview = "SERVER_PAGEVIEW"
	TriggerTypeAlways         = "ALWAYS"
	TriggerTypeConsentInit    = "CONSENT_INIT"
	TriggerTypeInit           = "INIT"

	// Mobile/Firebase container triggers
	TriggerTypeFirebaseAppException         = "FIREBASE_APP_EXCEPTION"
	TriggerTypeFirebaseAppUpdate            = "FIREBASE_APP_UPDATE"
	TriggerTypeFirebaseCampaign             = "FIREBASE_CAMPAIGN"
	TriggerTypeFirebaseFirstOpen            = "FIREBASE_FIRST_OPEN"
	TriggerTypeFirebaseInAppPurchase        = "FIREBASE_IN_APP_PURCHASE"
	TriggerTypeFirebaseNotificationDismiss  = "FIREBASE_NOTIFICATION_DISMISS"
	TriggerTypeFirebaseNotificationForeground = "FIREBASE_NOTIFICATION_FOREGROUND"
	TriggerTypeFirebaseNotificationOpen     = "FIREBASE_NOTIFICATION_OPEN"
	TriggerTypeFirebaseNotificationReceive  = "FIREBASE_NOTIFICATION_RECEIVE"
	TriggerTypeFirebaseOSUpdate             = "FIREBASE_OS_UPDATE"
	TriggerTypeFirebaseSessionStart         = "FIREBASE_SESSION_START"
	TriggerTypeFirebaseUserEngagement       = "FIREBASE_USER_ENGAGEMENT"

	// AMP container triggers
	TriggerTypeAMPClick      = "AMP_CLICK"
	TriggerTypeAMPTimer      = "AMP_TIMER"
	TriggerTypeAMPScroll     = "AMP_SCROLL"
	TriggerTypeAMPVisibility = "AMP_VISIBILITY"
)

// Tag types for GTM API.
// Common tag types used in web containers. For a complete list,
// see the GTM Tag Gallery documentation.
const (
	// Google Analytics
	TagTypeGA4Config = "gaawc" // Google Analytics 4 - Configuration
	TagTypeGA4Event  = "gaawe" // Google Analytics 4 - Event
	TagTypeUA        = "ua"    // Universal Analytics (deprecated)

	// Google Ads
	TagTypeGoogleAdsConversion  = "awct"
	TagTypeGoogleAdsRemarketing = "sp"

	// Custom
	TagTypeCustomHTML  = "html"
	TagTypeCustomImage = "img"

	// Floodlight
	TagTypeFloodlightCounter = "flc"
	TagTypeFloodlightSales   = "fls"

	// Third-party
	TagTypeFacebookPixel     = "baut"
	TagTypeLinkedInInsight   = "linkedin"
	TagTypeTwitterConversion = "twitter_website_tag"
)

// Variable types for GTM API.
// Reference: https://developers.google.com/tag-platform/tag-manager/api/v2/reference/accounts/containers/workspaces/variables
const (
	// Basic variables
	VariableTypeConstant         = "c"
	VariableTypeCustomJavaScript = "jsm"
	VariableTypeDataLayer        = "v"
	VariableTypeURL              = "u"
	VariableTypeFirstPartyCookie = "k"
	VariableTypeLookupTable      = "smm"
	VariableTypeRegexTable       = "remm"
	VariableTypeRandomNumber     = "r"
	VariableTypeJavaScript       = "j"

	// Google Analytics 4
	VariableTypeGA4EventSettings  = "gtes"
	VariableTypeGA4ConfigSettings = "gas"

	// Enhanced Conversions
	VariableTypeUserProvidedData = "awec"

	// Container variables
	VariableTypeContainerVersion = "ctv"
	VariableTypeDebugMode        = "dbg"
	VariableTypeEnvironmentName  = "env"
)

// Filter/comparison types for triggers.
// Used in trigger filters, custom event filters, and auto-event filters.
const (
	FilterTypeEquals         = "EQUALS"
	FilterTypeContains       = "CONTAINS"
	FilterTypeStartsWith     = "STARTS_WITH"
	FilterTypeEndsWith       = "ENDS_WITH"
	FilterTypeMatchesRegex   = "MATCHES_REGEX"
	FilterTypeGreaterThan    = "GREATER_THAN"
	FilterTypeGreaterOrEquals = "GREATER_OR_EQUALS"
	FilterTypeLessThan       = "LESS_THAN"
	FilterTypeLessOrEquals   = "LESS_OR_EQUALS"
	FilterTypeCSSSelector    = "CSS_SELECTOR"
	FilterTypeMatchRegex     = "MATCH_REGEX"
)

// Parameter type constants are defined in helpers.go:
// - ParameterTypeTemplate
// - ParameterTypeBoolean
// - ParameterTypeInteger
// - ParameterTypeList
// - ParameterTypeMap
// - ParameterTypeTagReference
// - ParameterTypeTriggerReference

// Tag firing options
const (
	TagFiringUnlimited    = "UNLIMITED"      // Fire every time the trigger fires
	TagFiringOncePerEvent = "ONCE_PER_EVENT" // Fire once per event
	TagFiringOncePerLoad  = "ONCE_PER_LOAD"  // Fire once per page load
)

// Built-in variable names available in GTM.
var BuiltInVariables = []string{
	"PAGE_URL",
	"PAGE_HOSTNAME",
	"PAGE_PATH",
	"REFERRER",
	"EVENT",
	"CLICK_ELEMENT",
	"CLICK_CLASSES",
	"CLICK_ID",
	"CLICK_TARGET",
	"CLICK_URL",
	"CLICK_TEXT",
	"FORM_ELEMENT",
	"FORM_CLASSES",
	"FORM_ID",
	"FORM_TARGET",
	"FORM_URL",
	"FORM_TEXT",
	"ERROR_MESSAGE",
	"ERROR_URL",
	"ERROR_LINE",
	"NEW_HISTORY_FRAGMENT",
	"OLD_HISTORY_FRAGMENT",
	"NEW_HISTORY_STATE",
	"OLD_HISTORY_STATE",
	"HISTORY_SOURCE",
	"VIDEO_PROVIDER",
	"VIDEO_URL",
	"VIDEO_TITLE",
	"VIDEO_DURATION",
	"VIDEO_PERCENT",
	"VIDEO_VISIBLE",
	"VIDEO_STATUS",
	"VIDEO_CURRENT_TIME",
	"SCROLL_DEPTH_THRESHOLD",
	"SCROLL_DEPTH_UNITS",
	"SCROLL_DIRECTION",
	"ELEMENT_VISIBILITY_RATIO",
	"ELEMENT_VISIBILITY_TIME",
	"ELEMENT_VISIBILITY_FIRST_TIME",
	"ELEMENT_VISIBILITY_RECENT_TIME",
}

// Common scroll depth percentages
var DefaultScrollPercentages = []int{10, 25, 50, 75, 90, 100}

// Workspace IDs
const DefaultWorkspaceID = "1"

// Tag Manager OAuth scopes
const (
	ScopeDeleteContainers      = "https://www.googleapis.com/auth/tagmanager.delete.containers"
	ScopeEditContainers        = "https://www.googleapis.com/auth/tagmanager.edit.containers"
	ScopeEditContainerVersions = "https://www.googleapis.com/auth/tagmanager.edit.containerversions"
	ScopeManageAccounts        = "https://www.googleapis.com/auth/tagmanager.manage.accounts"
	ScopeManageUsers           = "https://www.googleapis.com/auth/tagmanager.manage.users"
	ScopePublish               = "https://www.googleapis.com/auth/tagmanager.publish"
	ScopeReadonly              = "https://www.googleapis.com/auth/tagmanager.readonly"
)

// AllScopes contains all available Tag Manager scopes.
var AllScopes = []string{
	ScopeDeleteContainers,
	ScopeEditContainers,
	ScopeEditContainerVersions,
	ScopeManageAccounts,
	ScopeManageUsers,
	ScopePublish,
	ScopeReadonly,
}

// MinimumReadScopes contains minimum scopes for read operations.
var MinimumReadScopes = []string{
	ScopeReadonly,
}

// MinimumWriteScopes contains minimum scopes for write operations.
var MinimumWriteScopes = []string{
	ScopeEditContainers,
	ScopeEditContainerVersions,
}

// MinimumPublishScopes contains minimum scopes for publish operations.
var MinimumPublishScopes = []string{
	ScopePublish,
}

// GA4 constraints are defined in validation.go:
// - GA4EventNameMaxLength = 40
// - GA4ParameterNameMaxLength = 40
// Additional GA4 constraint
const GA4ParameterValueMaxLength = 100

// GTM constraints are defined in validation.go:
// - GTMNameMaxLength = 256
// - GTMNotesMaxLength = 5000

// ValidTriggerTypes contains all valid trigger type values.
var ValidTriggerTypes = map[string]bool{
	TriggerTypePageview:                     true,
	TriggerTypeDOMReady:                     true,
	TriggerTypeWindowLoaded:                 true,
	TriggerTypeCustomEvent:                  true,
	TriggerTypeTriggerGroup:                 true,
	TriggerTypeFormSubmission:               true,
	TriggerTypeClick:                        true,
	TriggerTypeLinkClick:                    true,
	TriggerTypeJSError:                      true,
	TriggerTypeHistoryChange:                true,
	TriggerTypeTimer:                        true,
	TriggerTypeScrollDepth:                  true,
	TriggerTypeElementVisibility:            true,
	TriggerTypeYouTubeVideo:                 true,
	TriggerTypeServerPageview:               true,
	TriggerTypeAlways:                       true,
	TriggerTypeConsentInit:                  true,
	TriggerTypeInit:                         true,
	TriggerTypeFirebaseAppException:         true,
	TriggerTypeFirebaseAppUpdate:            true,
	TriggerTypeFirebaseCampaign:             true,
	TriggerTypeFirebaseFirstOpen:            true,
	TriggerTypeFirebaseInAppPurchase:        true,
	TriggerTypeFirebaseNotificationDismiss:  true,
	TriggerTypeFirebaseNotificationForeground: true,
	TriggerTypeFirebaseNotificationOpen:     true,
	TriggerTypeFirebaseNotificationReceive:  true,
	TriggerTypeFirebaseOSUpdate:             true,
	TriggerTypeFirebaseSessionStart:         true,
	TriggerTypeFirebaseUserEngagement:       true,
	TriggerTypeAMPClick:                     true,
	TriggerTypeAMPTimer:                     true,
	TriggerTypeAMPScroll:                    true,
	TriggerTypeAMPVisibility:                true,
}

// ValidTagTypes contains all valid tag type values.
var ValidTagTypes = map[string]bool{
	TagTypeGA4Config:            true,
	TagTypeGA4Event:             true,
	TagTypeUA:                   true,
	TagTypeGoogleAdsConversion:  true,
	TagTypeGoogleAdsRemarketing: true,
	TagTypeCustomHTML:           true,
	TagTypeCustomImage:          true,
	TagTypeFloodlightCounter:    true,
	TagTypeFloodlightSales:      true,
	TagTypeFacebookPixel:        true,
	TagTypeLinkedInInsight:      true,
	TagTypeTwitterConversion:    true,
}

// ValidVariableTypes contains all valid variable type values.
var ValidVariableTypes = map[string]bool{
	VariableTypeConstant:         true,
	VariableTypeCustomJavaScript: true,
	VariableTypeDataLayer:        true,
	VariableTypeURL:              true,
	VariableTypeFirstPartyCookie: true,
	VariableTypeLookupTable:      true,
	VariableTypeRegexTable:       true,
	VariableTypeRandomNumber:     true,
	VariableTypeJavaScript:       true,
	VariableTypeGA4EventSettings:  true,
	VariableTypeGA4ConfigSettings: true,
	VariableTypeUserProvidedData: true,
	VariableTypeContainerVersion: true,
	VariableTypeDebugMode:        true,
	VariableTypeEnvironmentName:  true,
}

// ValidFilterTypes contains all valid filter type values.
var ValidFilterTypes = map[string]bool{
	FilterTypeEquals:         true,
	FilterTypeContains:       true,
	FilterTypeStartsWith:     true,
	FilterTypeEndsWith:       true,
	FilterTypeMatchesRegex:   true,
	FilterTypeGreaterThan:    true,
	FilterTypeGreaterOrEquals: true,
	FilterTypeLessThan:       true,
	FilterTypeLessOrEquals:   true,
	FilterTypeCSSSelector:    true,
	FilterTypeMatchRegex:     true,
}

// IsValidTriggerType checks if a trigger type is valid.
func IsValidTriggerType(triggerType string) bool {
	return ValidTriggerTypes[triggerType]
}

// IsValidTagType checks if a tag type is valid.
func IsValidTagType(tagType string) bool {
	return ValidTagTypes[tagType]
}

// IsValidVariableType checks if a variable type is valid.
func IsValidVariableType(variableType string) bool {
	return ValidVariableTypes[variableType]
}

// IsValidFilterType checks if a filter type is valid.
func IsValidFilterType(filterType string) bool {
	return ValidFilterTypes[filterType]
}

// IsBuiltInVariable checks if a variable name is a built-in variable.
func IsBuiltInVariable(name string) bool {
	for _, v := range BuiltInVariables {
		if v == name {
			return true
		}
	}
	return false
}
