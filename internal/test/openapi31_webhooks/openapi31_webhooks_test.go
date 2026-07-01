// Package openapi31_webhooks is the regression test for a component
// schema reachable only via a top-level `webhooks:` entry (OpenAPI 3.1+),
// asserting it still gets a real, usable Go type even though the default
// pruning pass runs (skip-prune is not set in config_3_1.yaml). The test
// is structural -- it constructs and round-trips the generated type --
// rather than string-matching the generated source.
package openapi31_webhooks

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	spec31 "github.com/oapi-codegen/oapi-codegen/v2/internal/test/openapi31_webhooks/spec_3_1"
)

// TestWebhookOnlySchema_TypeIsGenerated is the core regression check:
// before the fix, WidgetCreated -- referenced only from
// webhooks.widget\.created.post.requestBody, and not from any `paths:`
// operation nor any other `components/` schema -- looked orphaned to
// pruneUnusedComponents's reachability walk (which never visited
// swagger.Webhooks) and was deleted from the spec before
// WebhookOperationDefinitions ran. oapi-codegen exited 0, but the
// generated file referenced the now-missing type in
// WidgetCreatedJSONRequestBody, so it failed to compile with "undefined:
// WidgetCreated". This test just needs to compile and construct the type
// to prove the type declaration exists and is usable.
func TestWebhookOnlySchema_TypeIsGenerated(t *testing.T) {
	w := spec31.WidgetCreated{
		Id:   "widget_123",
		Name: "gizmo",
		Type: spec31.WidgetCreatedTypeWidgetCreated,
	}

	encoded, err := json.Marshal(w)
	require.NoError(t, err)
	assert.JSONEq(t, `{"id":"widget_123","name":"gizmo","type":"widget.created"}`, string(encoded))

	var decoded spec31.WidgetCreated
	require.NoError(t, json.Unmarshal(encoded, &decoded))
	assert.Equal(t, w, decoded)
	assert.True(t, decoded.Type.Valid())

	// The webhook's requestBody alias should also reference the same
	// generated type, confirming WebhookOperationDefinitions and the
	// general component-schema generation pass agree on this type.
	var viaRequestBody spec31.WidgetCreatedJSONRequestBody = w
	assert.Equal(t, w, spec31.WidgetCreated(viaRequestBody))
}
