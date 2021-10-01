# Release process

## Tools

### GoReleaser

For the release process we use [GoReleaser](https://goreleaser.com/), which can be installed from their [installation page](https://goreleaser.com/install/).

## How it works

Click the [create a pull request][develop-release-pr] link to open a PR from `develop` to `main`

[develop-release-pr]: https://github.com/snyk/snyk-iac-custom-rules/compare/main...develop?expand=1&title=Release%20develop%20to%20production&body=Release%20stable%20to%20production

There are two release processes:
1. Release of Golang binaries via GitHub Releases
2. Release of NPM package

We will discuss both now, but both of them are dependant on GitHub tags for versioning.

### Versioning

We use [svu](https://github.com/caarlos0/svu) to generate a new tag based on the commit message. The tag can then be pushed as per the following example:
```bash
$ git tag -a v1.0.0 -m "Major release v1"
$ git push origin v1.0.0
```
### GitHub Release

Once the PR from `develop` to `master` is merged, the `Release SDK` GitHub action will generate a new GitHub tag based on the commit message and push an updated tag to GitHub. Then, using `goreleaser`, it will create a new release and publish it to GitHub, containing the SDK binaries.

### NPM package

Once the PR from `develop` to `master` is merged, the `Release SDK` GitHub action will publish a new version of the NPM package to the Snyk registry.

To test the NPM distribution process, run the following commands using the latest commit SHA:
```
$ ./scripts/release-npm.sh <commit> <tag>

$ npm i -g dist/

$ snyk-iac-custom-rules help
```


## CI/CD
As part of every PR from a feature branch to the `develop` branch, we run both CircleCI as well as the `E2E Tests` and `Contract Tests` GitHub Actions, which run our shellspec tests in Windows, Linux, and MacOS. The CircleCI pipeline runs the `golangci-lint` linter, `gofmt`, and `go mod tidy`, and then it runs `shellspec` end-to-end tests and the Golang unit tests on a Linux distribution.

Once the PR is merged into `develop`, the `E2E Tests` and `Contract Tests` GitHub Actions run again. These actions also run in PRs opened from `develop` to `main`.

Once the PR for the `main` branch has been merged, the `Release SDK` GitHub action runs, which increments the GitHub tag and creates a new GitHub release of the SDK and publishes the NPM package to the Snyk NPM Registry.

