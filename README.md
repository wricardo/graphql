# GraphQL Utility Library

This library is a collection of utilities designed to assist in GraphQL development. It provides tools for introspection, schema mapping, and more, making it easier to work with GraphQL APIs.

## Features

- **Introspection**: Perform introspection queries to retrieve schema details from a GraphQL server.
- **Schema Mapping**: Generate a map of schema types and their corresponding fields.
- **Pretty Printing**: Format schema types and fields for easy readability.

## Installation

To use this library, you need to have Go installed on your system. You can install the library by running:

```bash
go get github.com/yourusername/graphql
```

## Usage

### Introspection

You can perform an introspection query to get the schema details of a GraphQL server:

```go
package main

import (
	"fmt"
	"github.com/yourusername/graphql"
)

func main() {
	addr := "http://yourgraphqlserver.com/graphql"
	response, err := graphql.Introspect(addr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Schema:", response)
}
```

### Schema Mapping

Generate a map of schema types and fields:

```go
schemaMap := graphql.GetSchemaMap(yourSchema)
for key, value := range schemaMap {
	fmt.Println(key, ":", value)
}
```

## Examples

You can find example usage in the `examples/` directory, which includes:

- `introspection/`: Demonstrates how to perform introspection queries.

