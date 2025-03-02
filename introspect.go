package graphql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Introspect(addr string) (IntrospectionResponse, error) {
	r := map[string]string{
		"operationName": "IntrospectionQuery",
		"query":         query,
	}
	encoded, err := json.Marshal(r)
	if err != nil {
		return IntrospectionResponse{}, err
	}
	req, err := http.NewRequest("POST", addr, bytes.NewBuffer(encoded))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return IntrospectionResponse{}, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return IntrospectionResponse{}, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return IntrospectionResponse{}, err
	}

	cfg := IntrospectionResponse{}
	err = json.Unmarshal(body, &cfg)
	if err != nil {
		return IntrospectionResponse{}, err
	}

	cfg.Data.Schema.Queries = cfg.Data.Schema.GetQueries()

	return cfg, nil
}

var query string = `
query IntrospectionQuery {
      __schema {
        
        queryType { name }
        mutationType { name }
        subscriptionType { name }
        types {
          ...FullType
        }
        directives {
          name
          description
          
          locations
          args {
            ...InputValue
          }
        }
      }
    }

    fragment FullType on __Type {
      kind
      name
      description
      
      fields(includeDeprecated: true) {
        name
        description
        args {
          ...InputValue
        }
        type {
          ...TypeRef
        }
        isDeprecated
        deprecationReason
      }
      inputFields {
        ...InputValue
      }
      interfaces {
        ...TypeRef
      }
      enumValues(includeDeprecated: true) {
        name
        description
        isDeprecated
        deprecationReason
      }
      possibleTypes {
        ...TypeRef
      }
    }

    fragment InputValue on __InputValue {
      name
      description
      type { ...TypeRef }
      defaultValue
      
      
    }

    fragment TypeRef on __Type {
      kind
      name
      ofType {
        kind
        name
        ofType {
          kind
          name
          ofType {
            kind
            name
            ofType {
              kind
              name
              ofType {
                kind
                name
                ofType {
                  kind
                  name
                  ofType {
                    kind
                    name
                  }
                }
              }
            }
          }
        }
      }
    }
`

// GetSchemaMap returns a map of the schema types
// returns a map of the schema types
func GetSchemaMapString(schema Schema) map[string]string {
	types := make(map[string]string)
	for _, v := range schema.Types {
		if v.Name == schema.QueryType.Name {
			for _, f := range v.Fields {
				types["query."+f.Name] = PrettyPrintField(f)
			}
		} else if v.Name == schema.MutationType.Name {
			for _, f := range v.Fields {
				types["mutation."+f.Name] = PrettyPrintField(f)
			}
		} else if v.Name == schema.SubscriptionType.Name {
			for _, f := range v.Fields {
				types["subscription."+f.Name] = PrettyPrintField(f)
			}
		} else if v.Kind == "SCALAR" {
			types["scalar."+v.Name] = PrettyPrintFullType(v)
		} else if v.Kind == "ENUM" {
			types["enum."+v.Name] = PrettyPrintFullType(v)
		} else if v.Kind == "INTERFACE" {
			types["interface."+v.Name] = PrettyPrintFullType(v)
		} else if v.Kind == "INPUT_OBJECT" {
			types["input."+v.Name] = PrettyPrintFullType(v)
		} else {
			types["type."+v.Name] = PrettyPrintFullType(v)
		}

	}
	return types
}

func PrettyPrintFullType(f FullType) string {
	return fmt.Sprintf("%s", f.String())
}

func PrettyPrintField(f Field) string {
	return fmt.Sprintf("%s(%s): %s", f.Name, ArgsToString(f.Args), f.Type.String())
}

func ArgsToString(args []InputValue) string {
	if len(args) == 0 {
		return ""
	}
	var s string
	for i, v := range args {
		s += fmt.Sprintf("%s: %s", v.Name, v.Type.String())
		if i < len(args)-1 {
			s += ", "
		}
	}
	return s
}
