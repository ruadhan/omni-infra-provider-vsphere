// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package provider //nolint:testpackage // white-box test of unexported splitTagSpec and matchTagsByName

import (
	"testing"

	"github.com/vmware/govmomi/vapi/tags"
)

func TestSplitTagSpec(t *testing.T) {
	for _, test := range []struct {
		spec     string
		category string
		name     string
	}{
		{spec: "backup-daily", category: "", name: "backup-daily"},
		{spec: "backup-policy/daily", category: "backup-policy", name: "daily"},
		{spec: "env/prod/extra", category: "env", name: "prod/extra"},
		{spec: "/name-only", category: "", name: "name-only"},
		{spec: "category-only/", category: "category-only", name: ""},
		{spec: "", category: "", name: ""},
	} {
		t.Run(test.spec, func(t *testing.T) {
			category, name := splitTagSpec(test.spec)

			if category != test.category || name != test.name {
				t.Errorf("splitTagSpec(%q) = (%q, %q), expected (%q, %q)", test.spec, category, name, test.category, test.name)
			}
		})
	}
}

func TestMatchTagsByName(t *testing.T) {
	allTags := []tags.Tag{
		{ID: "urn:1", Name: "daily", CategoryID: "backup-policy"},
		{ID: "urn:2", Name: "prod", CategoryID: "env"},
		{ID: "urn:3", Name: "daily", CategoryID: "snapshots"},
	}

	for _, test := range []struct {
		name        string
		expectedIDs []string
	}{
		{name: "prod", expectedIDs: []string{"urn:2"}},
		{name: "daily", expectedIDs: []string{"urn:1", "urn:3"}},
		{name: "missing", expectedIDs: nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			matches := matchTagsByName(allTags, test.name)

			if len(matches) != len(test.expectedIDs) {
				t.Fatalf("matchTagsByName(%q) returned %d matches, expected %d", test.name, len(matches), len(test.expectedIDs))
			}

			for i, match := range matches {
				if match.ID != test.expectedIDs[i] {
					t.Errorf("matchTagsByName(%q)[%d] = %q, expected %q", test.name, i, match.ID, test.expectedIDs[i])
				}
			}
		})
	}
}
