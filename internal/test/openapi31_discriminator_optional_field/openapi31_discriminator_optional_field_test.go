// Package openapi31_discriminator_optional_field is a regression test for
// two compile-time bugs in the generated From*/Merge* union
// discriminator-setter helpers, discovered while generating types for the
// real OpenAI OpenAPI spec (a large real-world 3.1 spec):
//
//  1. When a union member's discriminator property is optional, the
//     generated Go field is a pointer to a named string-based type (e.g.
//     `*Envelope1Type`), not a plain string. The old template body did
//     `v.Type = "value"`, assigning a bare untyped string constant, which
//     does not compile against a pointer field.
//  2. When a union member is itself another discriminated union (a
//     "nested union"), it has no flat discriminator field of its own --
//     only a `union json.RawMessage` field -- so `v.Type = "value"` fails
//     to compile outright with "type X has no field or method Type".
//
// The fix stops trying to assign the discriminator value onto the
// member's typed Go field (`v.Type = ...`) altogether when the member
// doesn't declare that field as one of the union's own flat properties.
// Instead it patches the discriminator key straight into the marshaled
// JSON object, which works regardless of the member's field type, its
// pointer-ness, or its very existence.
package openapi31_discriminator_optional_field

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	spec31 "github.com/oapi-codegen/oapi-codegen/v2/internal/test/openapi31_discriminator_optional_field/spec_3_1"
)

// TestRequiredDiscriminator_PlainStringField covers the baseline case: a
// member whose discriminator property is required, so the generated field
// is a plain named-string type. FromEnvelope0 must still compile and work.
func TestRequiredDiscriminator_PlainStringField(t *testing.T) {
	value := "hello"
	var env spec31.Envelope
	require.NoError(t, env.FromEnvelope0(spec31.Envelope0{Value: &value}))

	encoded, err := json.Marshal(env)
	require.NoError(t, err)
	assert.JSONEq(t, `{"type":"Envelope0","value":"hello"}`, string(encoded))

	disc, err := env.Discriminator()
	require.NoError(t, err)
	assert.Equal(t, "Envelope0", disc)
}

// TestOptionalDiscriminator_PointerField covers sub-bug 1: Envelope1's
// `type` property is optional, so the generated field is
// `*Envelope1Type`, a pointer to a named type. Before the fix, the
// generated FromEnvelope1/MergeEnvelope1 didn't compile at all
// ("cannot use ... as *Envelope1Type value in assignment").
func TestOptionalDiscriminator_PointerField(t *testing.T) {
	value := "world"
	var env spec31.Envelope
	require.NoError(t, env.FromEnvelope1(spec31.Envelope1{Value: &value}))

	encoded, err := json.Marshal(env)
	require.NoError(t, err)
	assert.JSONEq(t, `{"type":"Envelope1","value":"world"}`, string(encoded))

	disc, err := env.Discriminator()
	require.NoError(t, err)
	assert.Equal(t, "Envelope1", disc)

	value2, err := env.AsEnvelope1()
	require.NoError(t, err)
	require.NotNil(t, value2.Type)
	assert.Equal(t, spec31.Envelope1Type("Envelope1"), *value2.Type)
}

// TestNestedUnionDiscriminator_NoFlatField covers sub-bug 2: NestedUnion is
// itself a discriminated union over LeafA/LeafB, so it has no flat `Type`
// field of its own -- only a `union json.RawMessage` field. Before the
// fix, the generated FromNestedUnion/MergeNestedUnion didn't compile at
// all ("type NestedUnion has no field or method Type"). This mirrors the
// real spec's InputItem union, one of whose members (Item) is itself a
// discriminated union.
func TestNestedUnionDiscriminator_NoFlatField(t *testing.T) {
	var leaf spec31.NestedUnion
	a := "leaf-value"
	require.NoError(t, leaf.FromLeafA(spec31.LeafA{A: &a}))

	var env spec31.Envelope
	require.NoError(t, env.FromNestedUnion(leaf))

	encoded, err := json.Marshal(env)
	require.NoError(t, err)
	// The outer Envelope's own discriminator ("type": "NestedUnion") wins
	// over whatever the inner union's discriminator ("type": "leafA") was,
	// since From* always overwrites the top-level union payload.
	assert.JSONEq(t, `{"type":"NestedUnion","a":"leaf-value"}`, string(encoded))

	disc, err := env.Discriminator()
	require.NoError(t, err)
	assert.Equal(t, "NestedUnion", disc)
}

// TestMergeOptionalDiscriminator_PointerField exercises MergeEnvelope1,
// the merge counterpart of TestOptionalDiscriminator_PointerField, which
// had the identical compile bug in the template.
func TestMergeOptionalDiscriminator_PointerField(t *testing.T) {
	value := "first"
	var env spec31.Envelope
	require.NoError(t, env.FromEnvelope1(spec31.Envelope1{Value: &value}))

	value2 := "second"
	require.NoError(t, env.MergeEnvelope1(spec31.Envelope1{Value: &value2}))

	encoded, err := json.Marshal(env)
	require.NoError(t, err)
	assert.JSONEq(t, `{"type":"Envelope1","value":"second"}`, string(encoded))
}
