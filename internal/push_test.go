package internal

import (
	"context"
	"testing"

	"github.com/containerd/containerd/content"
	"github.com/containerd/containerd/remotes"
	"github.com/open-policy-agent/opa/util/test"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/stretchr/testify/assert"
	"oras.land/oras-go/pkg/oras"
)

func mockPushParams() *PushCommandParams {
	return &PushCommandParams{
		BundleRegistry: "docker.io/test/test:latest",
	}
}

func TestPush(t *testing.T) {
	files := map[string]string{
		"bundle.tar.gz": `test`,
	}

	oldLogin := login
	defer func() {
		login = oldLogin
	}()
	login = func(context.Context) (remotes.Resolver, error) {
		return nil, nil
	}

	oldConfigurePushOpts := configurePushOpts
	defer func() {
		configurePushOpts = oldConfigurePushOpts
	}()
	configurePushOpts = func(configDescriptor ocispec.Descriptor) []oras.PushOpt {
		// verifies the config manifest
		assert.Equal(t, "application/vnd.oci.image.config.v1+json", configDescriptor.MediaType)
		assert.Equal(t, "sha256:0f1289b1eadb7bd9c8b634b9b677f5e881d1d91a1bebad042d0171daff0e7288", configDescriptor.Digest.String())

		return []oras.PushOpt{}
	}

	oldPush := push
	defer func() {
		push = oldPush
	}()
	push = func(ctx context.Context, resolver remotes.Resolver, ref string, provider content.Provider, descriptors []ocispec.Descriptor, opts ...oras.PushOpt) (ocispec.Descriptor, error) {
		// verifies the image name and tag and the bundle
		assert.Equal(t, "docker.io/test/test:latest", ref)
		assert.Equal(t, 1, len(descriptors))
		assert.Equal(t, "application/vnd.oci.image.layer.v1.tar+gzip", descriptors[0].MediaType)
		assert.Equal(t, "sha256:9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08", descriptors[0].Digest.String())

		return ocispec.Descriptor{}, nil
	}

	test.WithTempFS(files, func(root string) {
		pushParams := mockPushParams()

		err := RunPush([]string{root + "/" + "bundle.tar.gz"}, pushParams)
		assert.Nil(t, err)
	})
}
