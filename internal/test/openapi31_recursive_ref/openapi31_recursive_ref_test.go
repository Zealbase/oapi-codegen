// Package openapi31_recursive_ref verifies that a `oneOf` branch expressed
// with `$recursiveRef: '#'` -- a self-reference back to the enclosing root
// schema -- generates a working recursive union type, exercising both the
// discriminator-mapping logic and general Go type generation. The test is
// structural -- it exercises the generated union accessors and JSON
// round-trips -- rather than string-matching the generated source.
package openapi31_recursive_ref

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	spec31 "github.com/oapi-codegen/oapi-codegen/v2/internal/test/openapi31_recursive_ref/spec_3_1"
)

// TestRecursiveRef_LeafBranch asserts that the `$ref`-based branch of the
// oneOf (LeafFilter) still round-trips normally alongside the
// `$recursiveRef` branch, and that its implicit discriminator-mapping key
// ("LeafFilter", derived from its $ref) is distinct from the
// `$recursiveRef` branch's own key.
func TestRecursiveRef_LeafBranch(t *testing.T) {
	var item spec31.CompoundFilter_Filters_Item
	key := "k"
	require.NoError(t, item.FromLeafFilter(spec31.LeafFilter{
		Type: spec31.Leaf,
		Key:  &key,
	}))

	encoded, err := json.Marshal(item)
	require.NoError(t, err)
	assert.JSONEq(t, `{"type":"LeafFilter","key":"k"}`, string(encoded))

	disc, err := item.Discriminator()
	require.NoError(t, err)
	assert.Equal(t, "LeafFilter", disc)

	back, err := item.AsLeafFilter()
	require.NoError(t, err)
	require.NotNil(t, back.Key)
	assert.Equal(t, "k", *back.Key)
}

// TestRecursiveRef_SelfReferentialBranch is the core regression check:
// before the fix, `$recursiveRef: '#'` inside CompoundFilter's own oneOf
// wasn't recognized as any kind of reference at all (kin-openapi doesn't
// parse the Draft 2019-09 keyword), so it fell into the "inline branch,
// no $ref, no constant discriminator value" path and generation failed
// with "discriminator: inline oneOf/anyOf branch has no $ref and no
// constant value ... to derive a mapping key from". It should instead
// behave exactly like `$ref: '#/components/schemas/CompoundFilter'`,
// producing a real recursive Go type (CompoundFilter containing
// []CompoundFilter_Filters_Item, one of whose variants is CompoundFilter
// itself), with its own implicit discriminator-mapping key ("CompoundFilter",
// derived from the synthesized ref path) distinct from LeafFilter's.
//
// Note: FromCompoundFilter/FromLeafFilter (like other ref-based union
// branches without an explicit x-discriminator-value) overwrite the
// branch's own `Type` field with the Go type name rather than a specific
// wire enum value, since CompoundFilter's `type` property isn't a `const`
// (it's `enum: [and, or]`, an operator selector orthogonal to which oneOf
// variant is in play). So this test drives values through JSON directly to
// exercise the recursive structure and dispatch by the synthesized
// discriminator key, rather than asserting on Type after a From* call.
func TestRecursiveRef_SelfReferentialBranch(t *testing.T) {
	raw := `{
		"type": "and",
		"filters": [
			{
				"type": "CompoundFilter",
				"filters": [
					{"type": "leaf", "key": "k"}
				]
			}
		]
	}`

	var outer spec31.CompoundFilter
	require.NoError(t, json.Unmarshal([]byte(raw), &outer))
	assert.Equal(t, spec31.And, outer.Type)
	require.Len(t, outer.Filters, 1)

	disc, err := outer.Filters[0].Discriminator()
	require.NoError(t, err)
	assert.Equal(t, "CompoundFilter", disc)

	middle, err := outer.Filters[0].AsCompoundFilter()
	require.NoError(t, err)
	require.Len(t, middle.Filters, 1)

	leafDisc, err := middle.Filters[0].Discriminator()
	require.NoError(t, err)
	assert.Equal(t, "leaf", leafDisc)

	leaf, err := middle.Filters[0].AsLeafFilter()
	require.NoError(t, err)
	assert.Equal(t, spec31.Leaf, leaf.Type)
	require.NotNil(t, leaf.Key)
	assert.Equal(t, "k", *leaf.Key)

	// Round-trip: build the same recursive structure programmatically via
	// From*, confirming FromCompoundFilter accepts a CompoundFilter value
	// containing itself without any special-casing needed by callers.
	var innerLeaf spec31.CompoundFilter_Filters_Item
	require.NoError(t, innerLeaf.FromLeafFilter(spec31.LeafFilter{Type: spec31.Leaf}))
	inner := spec31.CompoundFilter{Type: spec31.Or, Filters: []spec31.CompoundFilter_Filters_Item{innerLeaf}}

	var innerItem spec31.CompoundFilter_Filters_Item
	require.NoError(t, innerItem.FromCompoundFilter(inner))
	built := spec31.CompoundFilter{Type: spec31.And, Filters: []spec31.CompoundFilter_Filters_Item{innerItem}}

	encoded, err := json.Marshal(built)
	require.NoError(t, err)

	var roundTripped spec31.CompoundFilter
	require.NoError(t, json.Unmarshal(encoded, &roundTripped))
	require.Len(t, roundTripped.Filters, 1)
	gotMiddle, err := roundTripped.Filters[0].AsCompoundFilter()
	require.NoError(t, err)
	require.Len(t, gotMiddle.Filters, 1)
}
