// Package openapi31_discriminator verifies that a discriminated `oneOf`
// with inline (non-$ref) branches generates working union accessors even
// when nested one level inside an outer `anyOf: [<oneOf-with-discriminator>,
// {type: "null"}]`. The test is structural -- it exercises the generated
// union accessors and JSON round-trips -- rather than string-matching the
// generated source.
package openapi31_discriminator

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	spec31 "github.com/oapi-codegen/oapi-codegen/v2/internal/test/openapi31_discriminator/spec_3_1"
)

// TestInlineDiscriminatedUnion_ServerVAD asserts that the "server_vad"
// branch -- an inline schema with no $ref -- round-trips through the
// generated union accessors and that the discriminator property is
// populated from the branch's own `const` value on Marshal.
//
// Before the fix, every inline branch's implicit discriminator-mapping key
// was derived from RefPathToObjName(element.Ref), which degenerates to the
// empty string for inline schemas (no $ref to derive a name from). All
// inline branches collided on that single empty key, leaving the mapping
// with only one entry instead of one-per-branch, and generation failed
// with "discriminator: not all schemas were mapped".
func TestInlineDiscriminatedUnion_ServerVAD(t *testing.T) {
	threshold := float32(0.5)
	var td spec31.TurnDetection
	require.NoError(t, td.FromTurnDetection0(spec31.TurnDetection0{
		Threshold: &threshold,
	}))

	encoded, err := json.Marshal(td)
	require.NoError(t, err)
	assert.JSONEq(t, `{"type":"server_vad","threshold":0.5}`, string(encoded))

	disc, err := td.Discriminator()
	require.NoError(t, err)
	assert.Equal(t, "server_vad", disc)

	value, err := td.ValueByDiscriminator()
	require.NoError(t, err)
	variant, ok := value.(spec31.TurnDetection0)
	require.True(t, ok)
	require.NotNil(t, variant.Threshold)
	assert.Equal(t, float32(0.5), *variant.Threshold)
}

// TestInlineDiscriminatedUnion_SemanticVAD is the matching case for the
// second inline branch, confirming it maps to its own distinct
// discriminator value ("semantic_vad") rather than colliding with the
// first branch's mapping entry.
func TestInlineDiscriminatedUnion_SemanticVAD(t *testing.T) {
	eagerness := spec31.High
	var td spec31.TurnDetection
	require.NoError(t, td.FromTurnDetection1(spec31.TurnDetection1{
		Eagerness: &eagerness,
	}))

	encoded, err := json.Marshal(td)
	require.NoError(t, err)
	assert.JSONEq(t, `{"type":"semantic_vad","eagerness":"high"}`, string(encoded))

	disc, err := td.Discriminator()
	require.NoError(t, err)
	assert.Equal(t, "semantic_vad", disc)

	value, err := td.ValueByDiscriminator()
	require.NoError(t, err)
	variant, ok := value.(spec31.TurnDetection1)
	require.True(t, ok)
	require.NotNil(t, variant.Eagerness)
	assert.Equal(t, spec31.High, *variant.Eagerness)
}

// TestInlineDiscriminatedUnion_JSONUnmarshal asserts that decoding raw JSON
// for each branch and reading it back out through the As* accessor works,
// confirming the discriminator switch in ValueByDiscriminator resolves to
// the correct branch for both "server_vad" and "semantic_vad".
func TestInlineDiscriminatedUnion_JSONUnmarshal(t *testing.T) {
	var td spec31.TurnDetection
	require.NoError(t, json.Unmarshal([]byte(`{"type":"server_vad","threshold":0.7}`), &td))
	sv, err := td.AsTurnDetection0()
	require.NoError(t, err)
	require.NotNil(t, sv.Threshold)
	assert.Equal(t, float32(0.7), *sv.Threshold)

	var td2 spec31.TurnDetection
	require.NoError(t, json.Unmarshal([]byte(`{"type":"semantic_vad","eagerness":"low"}`), &td2))
	sem, err := td2.AsTurnDetection1()
	require.NoError(t, err)
	require.NotNil(t, sem.Eagerness)
	assert.Equal(t, spec31.Low, *sem.Eagerness)
}
