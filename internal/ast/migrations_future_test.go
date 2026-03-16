// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: BUSL-1.1

package ast

import "testing"

// These tests document future AST operations needed for full migration coverage.
// Each is skipped with a description of the required capability.

func TestFuture_ExtractToResource(t *testing.T) {
	t.Skip("TODO: extract nested block/attribute from one resource into a new standalone resource block. " +
		"Example: aws_s3_bucket.versioning {} → new aws_s3_bucket_versioning resource wired via bucket = aws_s3_bucket.X.id")

	// When implemented, this test should:
	// 1. Parse an aws_s3_bucket with a versioning {} block
	// 2. Remove the versioning block from aws_s3_bucket
	// 3. Create a new aws_s3_bucket_versioning resource
	// 4. Wire it to the bucket via bucket attribute
	// 5. Verify the output has both resources with correct references
}

func TestFuture_MoveAttributeToBlock(t *testing.T) {
	t.Skip("TODO: move a top-level attribute into a nested block. " +
		"Example: aws_instance.cpu_core_count → aws_instance.cpu_options { core_count }")

	// When implemented, this test should:
	// 1. Parse an aws_instance with cpu_core_count = 2
	// 2. Move it into cpu_options { core_count = 2 }
	// 3. Verify the attribute is removed from top level and present in nested block
}

func TestFuture_FlattenBlock(t *testing.T) {
	t.Skip("TODO: inline nested block attributes into parent. " +
		"Example: aws_elasticache_replication_group.cluster_mode { num_node_groups } → top-level num_node_groups")

	// When implemented, this test should:
	// 1. Parse a block with cluster_mode { num_node_groups = 3, replicas_per_node_group = 2 }
	// 2. Flatten into top-level num_node_groups = 3, replicas_per_node_group = 2
	// 3. Verify cluster_mode block is removed and attributes are at top level
}

func TestFuture_RemoveResourceWithRefWarnings(t *testing.T) {
	t.Skip("TODO: remove a resource block and add FIXME comments at all reference sites across files. " +
		"Example: remove aws_opsworks_stack and add FIXME wherever aws_opsworks_stack.X is referenced")

	// When implemented, this test should:
	// 1. Parse multiple files with an aws_opsworks_stack and references to it
	// 2. Remove the resource block
	// 3. Add FIXME comments at each reference site in other files
	// 4. Verify the resource is gone and all reference sites have comments
}
