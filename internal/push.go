package internal

import (
	"context"
	"fmt"
	"os"

	orascontext "oras.land/oras-go/pkg/context"
	"oras.land/oras-go/pkg/target"

	"oras.land/oras-go/pkg/content"
	"oras.land/oras-go/pkg/oras"

	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/snyk/snyk-iac-rules/util"
)

type PushCommandParams struct {
	BundleRegistry string
}

func RunPush(args []string, params *PushCommandParams) error {
	return pushBundle(orascontext.Background(), oras.Copy, params.BundleRegistry, args[0])
}

const configContents = `{
	"mediaType": "application/vnd.oci.image.config.v1+json"
}`

type copy func(ctx context.Context, from target.Target, fromRef string, to target.Target, toRef string, opts ...oras.CopyOpt) (v1.Descriptor, error)

func pushBundle(ctx context.Context, copy copy, repository string, bundlePath string) error {
	bundleData, err := os.ReadFile(bundlePath)
	if err != nil {
		return fmt.Errorf("read bundle: %v", err)
	}

	registry, err := content.NewRegistry(content.RegistryOptions{})
	if err != nil {
		return fmt.Errorf("create registry: %v", err)
	}

	store := content.NewMemory()

	bundleDesc, err := store.Add("bundle.tar.gz", util.CustomTarballLayerMediaType, bundleData)
	if err != nil {
		return fmt.Errorf("add bundle: %v", err)
	}

	configDesc, err := store.Add("config.json", util.CustomConfigMediaType, []byte(configContents))
	if err != nil {
		return fmt.Errorf("add config: %v", err)
	}

	manifestData, manifestDesc, err := content.GenerateManifest(&configDesc, nil, bundleDesc)
	if err != nil {
		return fmt.Errorf("generate manifest: %v", err)
	}

	if err := store.StoreManifest(repository, manifestDesc, manifestData); err != nil {
		return fmt.Errorf("store manifest: %v", err)
	}

	fmt.Printf("Uploading %s %s\n", bundleDesc.Digest.Encoded()[:12], bundlePath)

	if _, err := copy(ctx, store, repository, registry, repository); err != nil {
		return fmt.Errorf("Failed to push bundle to container registry: %v", err)
	}

	return nil
}
