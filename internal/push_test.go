package internal

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/pkg/content"
	"oras.land/oras-go/pkg/oras"
	"oras.land/oras-go/pkg/target"
)

func TestPush(t *testing.T) {
	copy := func(ctx context.Context, from target.Target, fromRef string, to target.Target, toRef string, opts ...oras.CopyOpt) (ocispec.Descriptor, error) {

		// Check that the from and to references are consistent

		if fromRef != "example.com/repository/bundle:v0.0.1" {
			t.Fatalf("unexpected from reference: %v", fromRef)
		}

		if toRef != "example.com/repository/bundle:v0.0.1" {
			t.Fatalf("unexpected to reference: %v", toRef)
		}

		// Dump everything from the "from" target into an in-memory registry
		// using the "from" and "to" references.

		store := content.NewMemory()

		if _, err := oras.Copy(context.Background(), from, fromRef, store, toRef); err != nil {
			t.Fatalf("copy: %v", err)
		}

		// Check that the manifest can be found and has the right media type.

		if _, desc, err := store.Resolve(context.Background(), toRef); err != nil {
			t.Fatalf("resolve manifest: %v", err)
		} else if desc.MediaType != "application/vnd.oci.image.manifest.v1+json" {
			t.Fatalf("invalid manifest media type: %v", desc.MediaType)
		}

		// Check that the configuration can be found and has the right media
		// type.

		if desc, _, found := store.GetByName("config.json"); !found {
			t.Fatalf("config not found")
		} else if desc.MediaType != "application/vnd.oci.image.config.v1+json" {
			t.Fatalf("invalid config media type: %v", desc.MediaType)
		}

		// Check that the manifest can be found, has the right media type, and
		// has the right content.

		if desc, data, found := store.GetByName("bundle.tar.gz"); !found {
			t.Fatalf("bundle not found")
		} else if desc.MediaType != "application/vnd.oci.image.layer.v1.tar+gzip" {
			t.Fatalf("invalid bundle media type: %v", desc.MediaType)
		} else if string(data) != "bundle content" {
			t.Fatalf("invalid bundle content: %v", string(data))
		}

		return ocispec.Descriptor{}, nil
	}

	bundle := filepath.Join(t.TempDir(), "bundle.tar.gz")

	if err := os.WriteFile(bundle, []byte("bundle content"), 0644); err != nil {
		t.Fatalf("write bundle content: %v", err)
	}

	if err := pushBundle(context.Background(), copy, "example.com/repository/bundle:v0.0.1", bundle); err != nil {
		t.Fatalf("error: %v", err)
	}
}
