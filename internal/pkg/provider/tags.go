// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/vmware/govmomi/vapi/rest"
	"github.com/vmware/govmomi/vapi/tags"
	"github.com/vmware/govmomi/vim25/mo"
	"go.uber.org/zap"
)

// splitTagSpec splits a machine class tag entry into its optional category and
// tag name parts. The "category/name" form pins the tag to a category; a plain
// name matches across all categories. The first slash separates the parts, so
// a tag name may not contain a slash unless a category is given.
func splitTagSpec(spec string) (category, name string) {
	if before, after, found := strings.Cut(spec, "/"); found {
		return before, after
	}

	return "", spec
}

// matchTagsByName returns all tags with the given name. A tag name is only
// unique within its category, so multiple matches are possible.
func matchTagsByName(allTags []tags.Tag, name string) []tags.Tag {
	var matches []tags.Tag

	for _, tag := range allTags {
		if tag.Name == name {
			matches = append(matches, tag)
		}
	}

	return matches
}

// resolveTagIDs resolves machine class tag specs ("name" or "category/name")
// to vCenter tag IDs. Plain names must be unambiguous across categories.
func resolveTagIDs(ctx context.Context, manager *tags.Manager, specs []string) ([]string, error) {
	ids := make([]string, 0, len(specs))

	var (
		allTags     []tags.Tag
		tagsFetched bool
	)

	for _, spec := range specs {
		category, name := splitTagSpec(spec)

		if name == "" {
			return nil, fmt.Errorf("invalid tag %q: empty tag name", spec)
		}

		if category != "" {
			cat, err := manager.GetCategory(ctx, category)
			if err != nil {
				return nil, fmt.Errorf("failed to find tag category %q: %w", category, err)
			}

			tag, err := manager.GetTagForCategory(ctx, name, cat.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to find tag %q in category %q: %w", name, category, err)
			}

			ids = append(ids, tag.ID)

			continue
		}

		if !tagsFetched {
			var err error

			allTags, err = manager.GetTags(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to list tags: %w", err)
			}

			tagsFetched = true
		}

		matches := matchTagsByName(allTags, name)

		switch len(matches) {
		case 0:
			return nil, fmt.Errorf("tag %q not found", name)
		case 1:
			ids = append(ids, matches[0].ID)
		default:
			return nil, fmt.Errorf("tag name %q exists in multiple categories, use \"category/name\" to disambiguate", name)
		}
	}

	return ids, nil
}

// attachTags attaches the given machine class tags to the VM. Tags are
// resolved and attached via the vSphere Automation (REST) API, which uses its
// own session next to the SOAP one; a fresh session is created per call.
// All tags are resolved before the first attach, and attaching an already
// attached tag is a no-op, so the operation is safe to retry.
func (p *Provisioner) attachTags(ctx context.Context, vmRef mo.Reference, specs []string) error {
	restClient := rest.NewClient(p.vsphereClient.Client)

	if err := restClient.Login(ctx, p.userInfo); err != nil {
		return fmt.Errorf("failed to log in to the vSphere automation API: %w", err)
	}

	defer func() {
		if err := restClient.Logout(ctx); err != nil {
			p.logger.Debug("failed to log out of the vSphere automation API", zap.Error(err))
		}
	}()

	manager := tags.NewManager(restClient)

	ids, err := resolveTagIDs(ctx, manager, specs)
	if err != nil {
		return err
	}

	for i, id := range ids {
		if err := manager.AttachTag(ctx, id, vmRef); err != nil {
			return fmt.Errorf("failed to attach tag %q: %w", specs[i], err)
		}
	}

	return nil
}
