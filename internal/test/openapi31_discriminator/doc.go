// Package openapi31_discriminator verifies that a discriminated `oneOf`
// with inline (non-$ref) branches, nested one level inside an outer
// `anyOf: [<oneOf-with-discriminator>, {type: "null"}]`, generates working
// union accessors. The generated types are emitted in the spec_3_1/
// subpackage and exercised by the tests in this directory.
package openapi31_discriminator

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=config_3_1.yaml spec_3_1.yaml
