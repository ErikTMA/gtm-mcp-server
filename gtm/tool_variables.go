package gtm

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ListVariablesInput struct {
	AccountID   string `json:"accountId" jsonschema:"description:The GTM account ID"`
	ContainerID string `json:"containerId" jsonschema:"description:The GTM container ID"`
	WorkspaceID string `json:"workspaceId" jsonschema:"description:The GTM workspace ID"`
}
type ListVariablesOutput struct {
	Variables []Variable `json:"variables"`
}

func registerListVariables(server *mcp.Server) {
	handler := func(ctx context.Context, req *mcp.CallToolRequest, input ListVariablesInput) (*mcp.CallToolResult, ListVariablesOutput, error) {
		client, err := getClient(ctx)
		if err != nil {
			return nil, ListVariablesOutput{}, err
		}

		// Validate account access
		if err := client.ValidateAccess(input.AccountID, input.ContainerID); err != nil {
			return nil, ListVariablesOutput{}, err
		}

		variables, err := client.ListVariables(ctx, input.AccountID, input.ContainerID, input.WorkspaceID)
		if err != nil {
			return nil, ListVariablesOutput{}, err
		}

		return nil, ListVariablesOutput{Variables: variables}, nil
	}

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_variables",
		Description: "List all variables in a GTM workspace",
	}, handler)
}

type GetVariableInput struct {
	AccountID   string `json:"accountId" jsonschema:"description:The GTM account ID"`
	ContainerID string `json:"containerId" jsonschema:"description:The GTM container ID"`
	WorkspaceID string `json:"workspaceId" jsonschema:"description:The GTM workspace ID"`
	VariableID  string `json:"variableId" jsonschema:"description:The variable ID to retrieve"`
}

type GetVariableOutput struct {
	Variable VariableDetail `json:"variable"`
}

func registerGetVariable(server *mcp.Server) {
	handler := func(ctx context.Context, req *mcp.CallToolRequest, input GetVariableInput) (*mcp.CallToolResult, GetVariableOutput, error) {
		client, err := getClient(ctx)
		if err != nil {
			return nil, GetVariableOutput{}, err
		}

		// Validate account access
		if err := client.ValidateAccess(input.AccountID, input.ContainerID); err != nil {
			return nil, GetVariableOutput{}, err
		}

		variable, err := client.GetVariable(ctx, input.AccountID, input.ContainerID, input.WorkspaceID, input.VariableID)
		if err != nil {
			return nil, GetVariableOutput{}, err
		}

		return nil, GetVariableOutput{Variable: *variable}, nil
	}

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_variable",
		Description: "Get a specific variable by ID with full configuration details including parameters",
	}, handler)
}
