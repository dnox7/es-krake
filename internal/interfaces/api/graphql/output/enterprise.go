package output

import (
	"github.com/dpe27/es-krake/internal/domain/enterprise/entity"
	"github.com/dpe27/es-krake/pkg/utils"

	"github.com/graphql-go/graphql"
)

func EnterpriseOutput() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "enterprise",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: AnyInt,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(entity.Enterprise).ID, nil
				},
			},
			"name": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(entity.Enterprise).Name, nil
				},
			},
			"is_active": &graphql.Field{
				Type: graphql.Boolean,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(entity.Enterprise).IsActive, nil
				},
			},
			"address": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(entity.Enterprise).Address, nil
				},
			},
			"phone": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(entity.Enterprise).Phone, nil
				},
			},
			"fax": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(entity.Enterprise).Fax, nil
				},
			},
			"mail_address": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(entity.Enterprise).MailAddress, nil
				},
			},
			"mail_signature": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(entity.Enterprise).MailSignature, nil
				},
			},
			"thumbnail_path": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(entity.Enterprise).ThumbnailPath, nil
				},
			},
			"website_url": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(entity.Enterprise).WebsiteURL, nil
				},
			},
			"notes": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return params.Source.(entity.Enterprise).Notes, nil
				},
			},
			"created_at": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					createdAt := params.Source.(entity.Enterprise).CreatedAt
					return utils.ToDateTimeSQL(&createdAt), nil
				},
			},
			"updated_at": &graphql.Field{
				Type: graphql.String,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					updatedAt := params.Source.(entity.Enterprise).UpdatedAt
					return utils.ToDateTimeSQL(&updatedAt), nil
				},
			},
		},
	})
}
