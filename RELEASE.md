# Release process

## Tools

### GoReleaser

For the release process we use [GoReleaser](https://goreleaser.com/), which can be installed from their [installation page](https://goreleaser.com/install/).

## How it works

There are two release processes:
1. Via GoReleaser, which includes the release of Golang binaries via GitHub Releases, Docker images, dep and rpm files, and brew and scoop repositories.
2. Via `npm publish`, which releases a new NPM package

Both of them are dependant on GitHub tags for versioning.

### Versioning

We use [svu](https://github.com/caarlos0/svu) to generate a new tag based on the commit message. 

### GitHub Release

Once a PR is merged, the CircleCI workflow will generate a new GitHub tag based on the commit message and push an updated tag to GitHub. Then, using `goreleaser`, it will create a new release and publish it to GitHub, containing the SDK binaries.

### NPM package

Once a PR is merged, the CircleCI workflow will publish a new version of the NPM package to the Snyk registry.

To test the NPM distribution process, run the following command using the latest tag in GitHub:
```
$ ./scripts/release-npm.sh --tag=v0.2.3

$ npm i -g dist/

$ snyk-iac-rules help
```


## CI/CD
As part of every PR from a feature branch to the release branch, we run both CircleCI as well as the `E2E Tests` and `Contract Tests` GitHub Actions, which run our shellspec tests in Windows, Linux, and MacOS. The CircleCI pipeline runs the `golangci-lint` linter, `gofmt`, and `go mod tidy`, and then it runs `shellspec` end-to-end tests and the Golang unit tests on a Linux distribution.

Once the PR is merged into `develop`, the CircleCI release process described in this document runs.

