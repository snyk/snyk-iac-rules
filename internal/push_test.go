package internal

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/pkg/content"
	"oras.land/oras-go/pkg/oras"
	"oras.land/oras-go/pkg/target"
)

func TestPush(t *testing.T) {
	copy := func(ctx context.Context, from target.Target, fromRef string, to target.Target, toRef string, opts ...oras.CopyOpt) (v1.Descriptor, error) {

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

		// Read the manifest corresponding to the target reference.

		_, manifestDesc, err := store.Resolve(context.Background(), toRef)
		if err != nil {
			t.Fatalf("resolve manifest: %v", err)
		} else if manifestDesc.MediaType != v1.MediaTypeImageManifest {
			t.Fatalf("invalid manifest media type: %v", manifestDesc.MediaType)
		}

		_, manifestData, found := store.Get(manifestDesc)
		if !found {
			t.Fatalf("manifest not found: %v", err)
		}

		var manifest v1.Manifest

		if err := json.Unmarshal(manifestData, &manifest); err != nil {
			t.Fatalf("manifest data is invalid: %v", err)
		}

		// Check that the configuration has the right media type.

		if manifest.Config.MediaType != v1.MediaTypeImageConfig {
			t.Fatalf("invalid config media type: %v", manifest.Config.MediaType)
		}

		// Check that the manifest has only one layer

		if len(manifest.Layers) != 1 {
			t.Fatalf("invalid layers: %v", manifest.Layers)
		}

		// Check that the only layer has the right media type and content.

		if desc, data, found := store.Get(manifest.Layers[0]); !found {
			t.Fatalf("bundle not found")
		} else if desc.MediaType != v1.MediaTypeImageLayerGzip {
			t.Fatalf("invalid bundle media type: %v", desc.MediaType)
		} else if string(data) != "bundle content" {
			t.Fatalf("invalid bundle content: %v", string(data))
		}

		return v1.Descriptor{}, nil
	}

	bundle := filepath.Join(t.TempDir(), "bundle.tar.gz")

	if err := os.WriteFile(bundle, []byte("bundle content"), 0644); err != nil {
		t.Fatalf("write bundle content: %v", err)
	}

	if err := pushBundle(context.Background(), copy, "example.com/repository/bundle:v0.0.1", bundle); err != nil {
		t.Fatalf("error: %v", err)
	}

}
