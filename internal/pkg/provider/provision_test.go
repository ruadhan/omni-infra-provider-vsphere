// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package provider

import "testing"

func TestClusterFolderName(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		requestSetID string
		expected     string
	}{
		{requestSetID: "prod-control-planes", expected: "prod"},
		{requestSetID: "prod-workers", expected: "prod"},
		{requestSetID: "my-cluster-control-planes", expected: "my-cluster"},
		{requestSetID: "my-cluster-workers", expected: "my-cluster"},
		// Custom-named machine sets keep the full ID.
		{requestSetID: "prod-storage-nodes", expected: "prod-storage-nodes"},
		// Degenerate IDs are kept as-is rather than reduced to an empty name.
		{requestSetID: "-workers", expected: "-workers"},
		{requestSetID: "workers", expected: "workers"},
	} {
		t.Run(test.requestSetID, func(t *testing.T) {
			t.Parallel()

			if actual := clusterFolderName(test.requestSetID); actual != test.expected {
				t.Errorf("clusterFolderName(%q) = %q, expected %q", test.requestSetID, actual, test.expected)
			}
		})
	}
}
