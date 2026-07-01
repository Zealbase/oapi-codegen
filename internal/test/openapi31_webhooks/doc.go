// Package openapi31_webhooks verifies that a component schema referenced
// only from a top-level OpenAPI 3.1 `webhooks:` entry (and not from any
// `paths:` operation or any other `components/` cross-reference) still
// gets its Go type generated, even though the default (skip-prune: false)
// pruning pass runs. The generated types are emitted in the spec_3_1/
// subpackage and exercised by the tests in this directory.
package openapi31_webhooks

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=config_3_1.yaml spec_3_1.yaml
