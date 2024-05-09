# 🏭 UDS <name-upper> Package

[![Latest Release](https://img.shields.io/github/v/release/<organization>/uds-package-<name>)](https://github.com/<organization>/uds-package-<name>/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/<organization>/uds-package-<name>/tag-and-release.yaml)](https://github.com/<organization>/uds-package-<name>/actions/workflows/tag-and-release.yaml)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/<organization>/uds-package-<name>/badge)](https://api.securityscorecards.dev/projects/github.com/<organization>/uds-package-<name>)

This package is designed for use as part of a bundle deployed on [UDS Core](https://github.com/defenseunicorns/uds-core).

## Prerequisites

<name-upper> requires a postgres database. Wiring coder to your dependencies is done primarily via helm values, which will require the use of a bundle created with uds-cli.

## Flavors

| Flavor | Description | Example Creation |
| ------ | ----------- | ---------------- |
| upstream | Uses upstream images within the package. | `zarf package create . -f upstream` |

## Releases

The released packages can be found in [ghcr](https://github.com/<organization>/uds-package-<name>/pkgs/container/packages%2Fuds%2F<name>).

## UDS Tasks (for local dev and CI)

*For local dev, this requires you install [uds-cli](https://github.com/defenseunicorns/uds-cli?tab=readme-ov-file#install)

> :white_check_mark: **Tip:** To get a list of tasks to run you can use `uds run --list`!

## Contributing

Please see the [CONTRIBUTING.md](./CONTRIBUTING.md)