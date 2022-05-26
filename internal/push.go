package internal

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/containerd/containerd/remotes"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	auth "oras.land/oras-go/pkg/auth/docker"
	orascontext "oras.land/oras-go/pkg/context"

	"oras.land/oras-go/pkg/content"
	"oras.land/oras-go/pkg/oras"

	util "github.com/snyk/snyk-iac-rules/util"
)

var push = oras.Push

type PushCommandParams struct {
	BundleRegistry string
}

const configContents = `{
	"mediaType": "application/vnd.oci.image.config.v1+json"
}`

func RunPush(args []string, params *PushCommandParams) error {
	ctx := orascontext.Background()
	memoryStore := content.NewMemoryStore()

	return pushBundle(ctx, memoryStore, params.BundleRegistry, args[0])
}

func pushBundle(ctx context.Context, memoryStore *content.Memorystore, repository string, path string) error {
	resolver, err := login(ctx)
	if err != nil {
		return err
	}

	// load the generated bundle and in-memory config.json into OCI descriptors
	bundle, err := loadBundle(ctx, memoryStore, path)
	if err != nil {
		return err
	}
	config := loadConfig(memoryStore)

	desriptors := []ocispec.Descriptor{
		bundle,
	}
	pushOpts := configurePushOpts(config)
	_, err = push(ctx, resolver, repository, memoryStore, desriptors, pushOpts...)
	if err != nil {
		return fmt.Errorf("Failed to push bundle to container registry: " + err.Error())
	}

	return nil
}

var login = func(ctx context.Context) (remotes.Resolver, error) {
	cli, err := auth.NewClient()
	if err != nil {
		return nil, fmt.Errorf("Failed to create authentication client: " + err.Error())
	}

	resolver, err := cli.Resolver(ctx, http.DefaultClient, false)
	if err != nil {
		return nil, fmt.Errorf("Failed to login: " + err.Error())
	}
	return resolver, nil
}

var loadBundle = func(ctx context.Context, memoryStore *content.Memorystore, path string) (ocispec.Descriptor, error) {
	bundle, err := os.Open(path)
	if err != nil {
		return ocispec.Descriptor{}, fmt.Errorf("Failed to find bundle: " + err.Error())
	}
	defer bundle.Close()

	bundleContents, err := ioutil.ReadAll(bundle)
	if err != nil {
		return ocispec.Descriptor{}, fmt.Errorf("Failed to read bundle contents: " + err.Error())
	}

	return memoryStore.Add(path, util.CustomTarballLayerMediaType, bundleContents), nil
}

var loadConfig = func(memoryStore *content.Memorystore) ocispec.Descriptor {
	descriptor := memoryStore.Add("config.json", util.CustomConfigMediaType, []byte(configContents))
	descriptor.Annotations = nil
	return descriptor
}

var configurePushOpts = func(config ocispec.Descriptor) []oras.PushOpt {
	return []oras.PushOpt{
		oras.WithPushStatusTrack(os.Stdout),
		oras.WithConfig(config),
	}
}
