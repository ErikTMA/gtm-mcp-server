package gtm

import (
	"context"
	"fmt"

	tagmanager "google.golang.org/api/tagmanager/v2"
)

// Variable is a simplified representation of a GTM variable.
type Variable struct {
	VariableID string `json:"variableId"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Path       string `json:"path"`
}

// ListVariables returns all variables in a workspace.
func (c *Client) ListVariables(ctx context.Context, accountID, containerID, workspaceID string) ([]Variable, error) {
	parent := fmt.Sprintf("accounts/%s/containers/%s/workspaces/%s", accountID, containerID, workspaceID)

	resp, err := c.Service.Accounts.Containers.Workspaces.Variables.List(parent).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	return toVariables(resp.Variable), nil
}

// GetVariable returns a specific variable by ID with full configuration details.
func (c *Client) GetVariable(ctx context.Context, accountID, containerID, workspaceID, variableID string) (*VariableDetail, error) {
	path := fmt.Sprintf("accounts/%s/containers/%s/workspaces/%s/variables/%s",
		accountID, containerID, workspaceID, variableID)

	variable, err := c.Service.Accounts.Containers.Workspaces.Variables.Get(path).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	result := toVariableDetail(variable)
	return &result, nil
}

func toVariables(variables []*tagmanager.Variable) []Variable {
	result := make([]Variable, 0, len(variables))
	for _, v := range variables {
		result = append(result, Variable{
			VariableID: v.VariableId,
			Name:       v.Name,
			Type:       v.Type,
			Path:       v.Path,
		})
	}
	return result
}

// VariableDetail provides full configuration details for a variable.
// Note: Parameter and FormatValue use 'any' to avoid JSON schema cycle detection issues
// since Parameter is a recursive type (contains List and Map of Parameters).
type VariableDetail struct {
	VariableID         string   `json:"variableId"`
	Name               string   `json:"name"`
	Type               string   `json:"type"`
	Path               string   `json:"path"`
	Parameter          any      `json:"parameter,omitempty"`
	FormatValue        any      `json:"formatValue,omitempty"`
	DisablingTriggerID []string `json:"disablingTriggerId,omitempty"`
	EnablingTriggerID  []string `json:"enablingTriggerId,omitempty"`
	Notes              string   `json:"notes,omitempty"`
	ScheduleStartMs    int64    `json:"scheduleStartMs,omitempty"`
	ScheduleEndMs      int64    `json:"scheduleEndMs,omitempty"`
	Fingerprint        string   `json:"fingerprint,omitempty"`
	ParentFolderID     string   `json:"parentFolderId,omitempty"`
}

// VariableFormatValue represents formatting options for a variable value.
type VariableFormatValue struct {
	CaseConversionType      string     `json:"caseConversionType,omitempty"`
	ConvertFalseToValue     *Parameter `json:"convertFalseToValue,omitempty"`
	ConvertNullToValue      *Parameter `json:"convertNullToValue,omitempty"`
	ConvertTrueToValue      *Parameter `json:"convertTrueToValue,omitempty"`
	ConvertUndefinedToValue *Parameter `json:"convertUndefinedToValue,omitempty"`
}

func toVariableDetail(v *tagmanager.Variable) VariableDetail {
	detail := VariableDetail{
		VariableID:         v.VariableId,
		Name:               v.Name,
		Type:               v.Type,
		Path:               v.Path,
		DisablingTriggerID: v.DisablingTriggerId,
		EnablingTriggerID:  v.EnablingTriggerId,
		Notes:              v.Notes,
		ScheduleStartMs:    v.ScheduleStartMs,
		ScheduleEndMs:      v.ScheduleEndMs,
		Fingerprint:        v.Fingerprint,
		ParentFolderID:     v.ParentFolderId,
	}

	// Convert parameters
	if v.Parameter != nil {
		detail.Parameter = convertAPIParameters(v.Parameter)
	}

	// Convert format value
	if v.FormatValue != nil {
		fv := convertFormatValue(v.FormatValue)
		detail.FormatValue = &fv
	}

	return detail
}

func convertFormatValue(fv *tagmanager.VariableFormatValue) VariableFormatValue {
	result := VariableFormatValue{
		CaseConversionType: fv.CaseConversionType,
	}

	if fv.ConvertFalseToValue != nil {
		p := convertAPIParameter(fv.ConvertFalseToValue)
		result.ConvertFalseToValue = &p
	}
	if fv.ConvertNullToValue != nil {
		p := convertAPIParameter(fv.ConvertNullToValue)
		result.ConvertNullToValue = &p
	}
	if fv.ConvertTrueToValue != nil {
		p := convertAPIParameter(fv.ConvertTrueToValue)
		result.ConvertTrueToValue = &p
	}
	if fv.ConvertUndefinedToValue != nil {
		p := convertAPIParameter(fv.ConvertUndefinedToValue)
		result.ConvertUndefinedToValue = &p
	}

	return result
}

func convertAPIParameters(params []*tagmanager.Parameter) []Parameter {
	result := make([]Parameter, 0, len(params))
	for _, p := range params {
		result = append(result, convertAPIParameter(p))
	}
	return result
}

func convertAPIParameter(p *tagmanager.Parameter) Parameter {
	param := Parameter{
		Type:  p.Type,
		Key:   p.Key,
		Value: p.Value,
	}

	if p.List != nil {
		param.List = convertAPIParameters(p.List)
	}
	if p.Map != nil {
		param.Map = convertAPIParameters(p.Map)
	}

	return param
}
