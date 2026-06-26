// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package provider

// Data is the provider custom machine config.
type Data struct {
	Datacenter   string `yaml:"datacenter"`
	ResourcePool string `yaml:"resource_pool"`
	Datastore    string `yaml:"datastore"`
	// StoragePolicy is the name of a vSphere Storage Policy (SPBM) to apply to the
	// cloned VM (home and disks). Optional; when empty the datastore default policy is used.
	StoragePolicy string `yaml:"storage_policy"`
	Network       string `yaml:"network"`
	Template      string `yaml:"template"`  // VM template name to clone from
	Folder        string `yaml:"folder"`    // VM folder path (optional)
	CACert        string `yaml:"ca_cert"`   // PEM-encoded CA certificate (optional)
	DiskSize      uint64 `yaml:"disk_size"` // GiB
	CPU           uint   `yaml:"cpu"`
	Memory        uint   `yaml:"memory"` // MiB
}
