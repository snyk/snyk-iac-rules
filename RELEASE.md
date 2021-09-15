# Release process

## Tools

### GoReleaser

For the release process we use [GoReleaser](https://goreleaser.com/), which can be installed from their [installation page](https://goreleaser.com/install/).

## How it works

If you're not creating a new major release of the SDK, then click the [create a pull request][develop-release-pr] link to open a PR from `develop` to `main`

[develop-release-pr]: https://github.com/snyk/snyk-iac-custom-rules/compare/main...develop?expand=1&title=Release%20develop%20to%20production&body=Release%20stable%20to%20production

Once this PR is merged, the `Release SDK` GitHub action will get the latest GitHub tag, increment its patch number, and push an updated tag to GitHub. Then, using `goreleaser`, it will create a new release and publish it to GitHub, containing the SDK binaries.

If you want to create a new major release of the SDK, then there is one extra step that needs to happen before merging the PR. That is, to create a new tag for the major version. Assuming we're at major version v0, then the following command must be run:
```
$ git tag -a v1.0.0 -m "Major release v1"
$ git push origin v1.0.0
```

If you want to create a new patch release of the SDK, then either follow the steps above or run the `./scripts/increment-release-tag.sh` script.

If you want to manually create a new release, download `goreleaser` and run the following command after creating a new tag:
```
$ goreleaser release --rm-dist
```

## CI/CD
As part of every PR from a feature branch to the `develop` branch, we run both CircleCI as well as the `E2E Tests` and `Contract Tests` GitHub Actions, which run our shellspec tests in Windows, Linux, and MacOS. The CircleCI pipeline runs the `golangci-lint` linter, `gofmt`, and `go mod tidy`, and then it runs `shellspec` end-to-end tests and the Golang unit tests on a Linux distribution.

Once the PR is merged into `develop`, the `E2E Tests` and `Contract Tests` GitHub Actions run again. These actions also run in PRs opened from `develop` to `main`.

Once the PR for the `main` branch has been merged, the `Release SDK` GitHub action runs, which increments the GitHub tag and creates a new GitHub release of the SDK.

