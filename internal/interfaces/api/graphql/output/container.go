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
		LoginHistoryOutput(),
		PermissionOutput(outputTypes, usecases.AuthUsecase),
		RoleOutput(outputTypes, usecases.AuthUsecase),

		DepartmentOutput(),
		PlatformAccountOutput(outputTypes, usecases.AuthUsecase, usecases.PlatformUsecase),

		EnterpriseOutput(),
		EnterpriseAccountOutput(outputTypes, usecases.AuthUsecase),
	} {
		outputTypes[graphqlType.Name()] = graphqlType
	}

	return outputTypes
}
