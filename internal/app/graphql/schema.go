package graphql

import "github.com/graphql-go/graphql"

func NewSchema() (*graphql.Schema, error) {
	q := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Name: "name",
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return "Kevin", nil
				},
			},
		},
	})

	s, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:        q,
		Mutation:     nil,
		Subscription: nil,
		Types:        nil,
		Directives:   nil,
		Extensions:   nil,
	})
	if err != nil {
		return nil, err
	}
	return &s, nil
}
