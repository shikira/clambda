package usecase

import (
	"context"
	"fmt"

	"github.com/shirasu/delambda/internal/domain/function"
	"github.com/shirasu/delambda/internal/domain/stack"
)

// ListStackFunctionsUseCase handles listing Lambda functions in a CloudFormation stack
type ListStackFunctionsUseCase struct {
	functionRepo function.Repository
	stackRepo    stack.Repository
}

// NewListStackFunctionsUseCase creates a new ListStackFunctionsUseCase
func NewListStackFunctionsUseCase(functionRepo function.Repository, stackRepo stack.Repository) *ListStackFunctionsUseCase {
	return &ListStackFunctionsUseCase{
		functionRepo: functionRepo,
		stackRepo:    stackRepo,
	}
}

// Execute lists all Lambda functions in the specified CloudFormation stack
func (uc *ListStackFunctionsUseCase) Execute(ctx context.Context, stackName string) ([]*function.Function, error) {
	// Get all Lambda function names in the stack
	functionNames, err := uc.stackRepo.ListLambdaFunctions(ctx, stackName)
	if err != nil {
		return nil, fmt.Errorf("failed to list Lambda functions in stack: %w", err)
	}

	if len(functionNames) == 0 {
		return []*function.Function{}, nil
	}

	// Fetch details for each function
	functions := make([]*function.Function, 0, len(functionNames))
	for _, functionName := range functionNames {
		fn, err := uc.functionRepo.FindByName(ctx, functionName)
		if err != nil {
			return nil, fmt.Errorf("failed to get function %s: %w", functionName, err)
		}
		functions = append(functions, fn)
	}

	return functions, nil
}
