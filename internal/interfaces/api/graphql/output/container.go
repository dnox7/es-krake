package output

import (
	"github.com/dpe27/es-krake/internal/usecase"
	"github.com/graphql-go/graphql"
)

func NewOutputTypes(usecases *usecase.UsecasesContainer) map[string]*graphql.Object {
	outputTypes := make(map[string]*graphql.Object)

	for _, graphqlType := range []*graphql.Object{
		RoleTypeOutput(),
		AccessOperationOutput(),
		PlatformAccountOutput(outputTypes, usecases.AuthUsecase, usecases.PlatformUsecase),
		PermissionOutput(outputTypes, usecases.AuthUsecase),
		RoleOutput(outputTypes, usecases.AuthUsecase),
		DepartmentOutput(),
		LoginHistoryOutput(),
	} {
		outputTypes[graphqlType.Name()] = graphqlType
	}

	return outputTypes
}
