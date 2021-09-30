# Release process

## Tools

### GoReleaser

For the release process we use [GoReleaser](https://goreleaser.com/), which can be installed from their [installation page](https://goreleaser.com/install/).

## How it works

If you're not creating a new major release of the SDK, then click the [create a pull request][develop-release-pr] link to open a PR from `develop` to `main`

[develop-release-pr]: https://github.com/snyk/snyk-iac-rules/compare/main...develop?expand=1&title=Release%20develop%20to%20production&body=Release%20stable%20to%20production

Once this PR is merged, the `Release SDK` GitHub action will generate a new GitHub Tag from the commit message using [svu](https://github.com/caarlos0/svu) and push it to GitHub. Then, using `goreleaser`, it will create a new release and publish it to GitHub, containing the SDK binaries.

## CI/CD
As part of every PR from a feature branch to the `develop` branch, we run both CircleCI as well as the `E2E Tests` and `Contract Tests` GitHub Actions, which run our shellspec tests in Windows, Linux, and MacOS. The CircleCI pipeline runs the `golangci-lint` linter, `gofmt`, and `go mod tidy`, and then it runs `shellspec` end-to-end tests and the Golang unit tests on a Linux distribution.

Once the PR is merged into `develop`, the `E2E Tests` and `Contract Tests` GitHub Actions run again. These actions also run in PRs opened from `develop` to `main`. [Open Production Release PR](https://github.com/snyk/snyk-iac-rules/compare/main...develop?expand=1)

Once the PR for the `main` branch has been merged, the `Release SDK` GitHub action runs, which increments the GitHub tag and creates a new GitHub release of the SDK.

