package main

import (
	"fmt"

	"github.com/wricardo/graphql"
)

func main() {
	res, err := graphql.Introspect("https://countries.trevorblades.com/graphql")
	if err != nil {
		fmt.Println(err)
		return
	}
	smap := graphql.GetSchemaMap(res.Data.Schema)
	for k, _ := range smap {
		fmt.Println(k)
		/*
			enum.__DirectiveLocation
			enum.__TypeKind
			input.ContinentFilterInput
			input.CountryFilterInput
			input.LanguageFilterInput
			input.StringQueryOperatorInput
			query.continent
			query.continents
			query.countries
			query.country
			query.language
			query.languages
			scalar.Boolean
			scalar.Float
			scalar.ID
			scalar.Int
			scalar.String
			type.Continent
			type.Country
			type.Language
			type.State
			type.Subdivision
			type.__Schema
			type.__Directive
			type.__EnumValue
			type.__Field
			type.__InputValue
			type.__Type
		*/
	}
}
