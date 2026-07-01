// Package openapi31_recursive_ref verifies that an inline `oneOf` branch
// using the JSON Schema Draft 2019-09 `$recursiveRef: '#'` keyword (a
// self-reference back to the enclosing root schema) is treated as
// equivalent to a `$ref` back to that schema, both for discriminator
// mapping-key derivation and for general Go type generation of recursive
// types. The generated types are emitted in the spec_3_1/ subpackage and
// exercised by the tests in this directory.
package openapi31_recursive_ref

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=config_3_1.yaml spec_3_1.yaml
