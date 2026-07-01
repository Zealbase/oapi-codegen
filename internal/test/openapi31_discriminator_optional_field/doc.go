// Package openapi31_discriminator_optional_field verifies that the
// generated From*/Merge* union discriminator-setter helpers compile and
// behave correctly when a union member's discriminator property is
// optional (pointer-typed) or when a union member is itself another
// discriminated union with no flat discriminator field of its own. The
// generated types are emitted in the spec_3_1/ subpackage and exercised by
// the tests in this directory.
package openapi31_discriminator_optional_field

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=config_3_1.yaml spec_3_1.yaml
