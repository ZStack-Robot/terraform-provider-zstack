// Copyright (c) ZStack.io, Inc.

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

// applyUuidOrNameFilter centralises the AI-vs-human lookup pattern for list
// data sources. When `uuid` is provided (non-null, non-unknown, non-empty),
// it short-circuits to a uuid= query and returns true so callers know the
// result is expected to contain at most one record. Otherwise it falls back
// to the existing name / name_pattern behaviour.
//
// Background: AI agents almost always have a stable UUID at hand (carried
// from prior steps / state) but no reliable name. Humans author with names.
// Letting both cohorts use the same data source — by adding `uuid` as a
// purely additive optional input — avoids forcing AI to guess at name
// patterns while keeping ergonomic name-based lookups intact.
//
// The function only mutates `q`; it does no validation of conflicts —
// schema-level ConflictsWith handles user-level errors before this runs.
func applyUuidOrNameFilter(q *param.QueryParam, uuid, name, namePattern types.String) (uuidApplied bool) {
	if !uuid.IsNull() && !uuid.IsUnknown() && uuid.ValueString() != "" {
		q.AddQ("uuid=" + uuid.ValueString())
		return true
	}
	if !name.IsNull() && !name.IsUnknown() && name.ValueString() != "" {
		q.AddQ("name=" + name.ValueString())
	}
	if !namePattern.IsNull() && !namePattern.IsUnknown() && namePattern.ValueString() != "" {
		q.AddQ("name~=" + namePattern.ValueString())
	}
	return false
}
