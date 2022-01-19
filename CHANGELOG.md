# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Add `validateResourcesExist` flag to Validator for enabling/disabling `key.IsManagedByFlux` validation.

### Changed

- Narrow down the `key.IsManagedByFlux` scope to validating ConfigMap and Secrets existence.

## [6.4.0] - 2022-01-18

### Added

- Add `key.ChartName` for removing workload cluster ID if its present.

## [6.3.0] - 2022-01-13

### Added

- Add `validateTargetNamespace` to ensure users are not allowed to create in-cluster Apps outside org- and WC-related namespaces.

## [6.2.0] - 2021-12-21

### Added

- Add `ValidateAppUpdate` to ensure `.spec.Namespace` is immutable in App CRs.

## [6.1.0] - 2021-12-09

### Added

- Support for App CRs with a `v` prefixed version. This enables Flux to automatically update the version based on its image tag.

## [6.0.0] - 2021-11-29

- Drop `apiextensions` dependency in favor of `apiextensions-application`.

## [5.6.0] - 2021-11-29

### Changed

- Allow setting `.spec.kubeConfig.inCluster` when generating app CRs.

## [5.5.0] - 2021-11-25

### Added

- Skip validation of configmap and secret names when `giantswarm.io/managedby`
label is set to `flux`.

## [5.4.0] - 2021-11-09

### Added

- Validate kubeconfig secret name is set when in cluster is false.

## [5.3.0] - 2021-09-15

### Added

- Make app CR namespace configurable.

## [5.2.3] - 2021-08-26

### Fixed

- Fix app-admission-controller webhook name in validation error matchers.

## [5.2.2] - 2021-08-17

### Fixed

- Don't restrict user values configmap name for NGINX Ingress Controller.

## [5.2.1] - 2021-08-17

### Fixed

- Check for nil or empty `.metadata.Name`  when validating `.spec.catalog`.

## [5.2.0] - 2021-08-16

### Changed

- Validate `Catalog` CRs instead of `AppCatalog` in App validation.

## [5.1.0] - 2021-08-10

- Reject app version label with legacy 1.0.0 value.

## [5.0.1] - 2021-06-15

### Fixed

- Fix typo in provider-specific CRD path in apiextensions.

## [5.0.0] - 2021-06-04

### Added

- Add `crd.LoadCRDs`, `crd.LoadCRD` functions.
- Add key functions for `Catalog` CRs.

### Changed

- Breaking change to replace `AppCatalog ` CRD with namespace scoped `Catalog`
CRD in `values` package.

## [4.13.0] - 2021-05-12

## [4.12.0] - 2021-05-06

### Changed

- Get metadata constants from k8smetadata library not apiextensions.

## [4.11.0] - 2021-04-27

### Added

- Add `InstallSkipCRDs` key function for app CRs.

## [4.10.0] - 2021-04-19

### Added

- Add validation for length of `metadata.name`.

## [4.9.0] - 2021-03-29

### Added

- Add validation for user configmap and secret names for apps in the default catalog.

## [4.8.0] - 2021-03-18

### Added

- Add `namespaceConfig` validation.

## [4.7.0] - 2021-03-05

### Added

- Add `key.AppTeam` function.

## [4.6.0] - 2021-03-02

### Added

- Add `compatibleProvider` metadata validation.

## [4.5.0] - 2021-02-25

### Added

- Add `namespace` metadata validation.
- Add `application.giantswarm.io/owners` annotation.

## [4.4.0] - 2021-02-19

### Added

- Add `key.AppCatalogEntryName` and `key.AppCatalogEntryTeam` functions.
- Add `application.giantswarm.io/team` annotation.

## [4.3.0] - 2021-02-03

### Added

- Add `key.ToChart` function.

## [4.2.0] - 2021-01-12

### Removed

- Remove unused errors from validation package.
- Do not set `config-controller.giantswarm.io/version` label to "0.0.0" on created App CRs.
- Remove `PauseReconciliation` option, responsible for setting `app-operator.giantswarm.io/paused` flag.

## [4.1.0] - 2021-01-05

### Added

- Add `PauseReconciliation` option, responsible for setting `app-operator.giantswarm.io/paused` flag.

### Removed

- Do not validate App CR configmap and secret names if managed by config-controller.

## [4.0.0] - 2020-12-03

### Changed

- Remove helmclient.MergeValue functions usage.
- Return interface map from merge functions.

## [3.7.0] - 2020-12-02

### Added

- Validate App CR configmap and secret names if managed by config-controller.

### Changed

- Change (unused yet) `config.giantswarm.io/major-version` annotation to `config.giantswarm.io/version`.

## [3.6.0] - 2020-12-01

### Added

- Support `ConfigMajorVersion` setting to set
  "config.giantswarm.io/major-version" annotation.
- Set "config-controller.giantswarm.io/version" label to "0.0.0" on created App
  CRs.

## [3.5.0] - 2020-11-27

### Added

- Return separate errors for cluster kubeconfig and configmap not existing
since there can be a delay creating them on cluster creation.

## [3.4.0] - 2020-11-26

### Added

- Allow configmap and secret configuration.

## [3.3.0] - 2020-11-23

### Added

- Add key functions for app labels.

## [3.2.0] - 2020-11-11

### Added

- Add key functions for cluster configmap and cluster kubeconfig names.

## [3.1.1] - 2020-11-10

### Fixed

- Move validation package to pkg.

## [3.1.0] - 2020-11-05

### Added

- Add validation package extracted from the validation resource in app-operator.

## [3.0.0] - 2020-11-04

- Add values service extracted from app-operator.

### Added

- Add annotation and key packages extracted from app-operator.

### Changed

- Updated apiextensions to v3.4.0.
- Prepare module v3.

## [2.0.0] - 2020-08-11

### Changed

- Updated Kubernetes dependencies to v1.18.5.

## [0.2.3] - 2020-06-23

### Changed

- Update apiextensions to avoid displaying empty strings in app CRs.

## [0.2.2] - 2020-06-01

### Changed

- Set version label value to 0.0.0 so control plane app CRs are reconciled by
  app-operator-unique.

## [0.2.1] - 2020-04-24

- Fix module path (was accidentaly declared as gitlab.com/...).

## [0.2.0] - 2020-04-24

### Changed

- migrate from dep to go modules (build-only changes)

## [0.1.0] - 2020-04-24

### Added

- First release

[Unreleased]: https://github.com/giantswarm/app/compare/v6.4.0...HEAD
[6.4.0]: https://github.com/giantswarm/app/compare/v6.3.0...v6.4.0
[6.3.0]: https://github.com/giantswarm/app/compare/v6.2.0...v6.3.0
[6.2.0]: https://github.com/giantswarm/app/compare/v6.1.0...v6.2.0
[6.1.0]: https://github.com/giantswarm/app/compare/v6.0.0...v6.1.0
[6.0.0]: https://github.com/giantswarm/app/compare/v5.6.0...v6.0.0
[5.6.0]: https://github.com/giantswarm/app/compare/v5.5.0...v5.6.0
[5.5.0]: https://github.com/giantswarm/app/compare/v5.4.0...v5.5.0
[5.4.0]: https://github.com/giantswarm/app/compare/v5.3.0...v5.4.0
[5.3.0]: https://github.com/giantswarm/app/compare/v5.2.3...v5.3.0
[5.2.3]: https://github.com/giantswarm/app/compare/v5.2.2...v5.2.3
[5.2.2]: https://github.com/giantswarm/app/compare/v5.2.1...v5.2.2
[5.2.1]: https://github.com/giantswarm/app/compare/v5.2.0...v5.2.1
[5.2.0]: https://github.com/giantswarm/app/compare/v5.1.0...v5.2.0
[5.1.0]: https://github.com/giantswarm/app/compare/v5.0.1...v5.1.0
[5.0.1]: https://github.com/giantswarm/app/compare/v5.0.0...v5.0.1
[5.0.0]: https://github.com/giantswarm/app/compare/v4.13.0...v5.0.0
[4.13.0]: https://github.com/giantswarm/app/compare/v4.12.0...v4.13.0
[4.12.0]: https://github.com/giantswarm/app/compare/v4.11.0...v4.12.0
[4.11.0]: https://github.com/giantswarm/app/compare/v4.10.0...v4.11.0
[4.10.0]: https://github.com/giantswarm/app/compare/v4.9.0...v4.10.0
[4.9.0]: https://github.com/giantswarm/app/compare/v4.8.0...v4.9.0
[4.8.0]: https://github.com/giantswarm/app/compare/v4.7.0...v4.8.0
[4.7.0]: https://github.com/giantswarm/app/compare/v4.6.0...v4.7.0
[4.6.0]: https://github.com/giantswarm/app/compare/v4.5.0...v4.6.0
[4.5.0]: https://github.com/giantswarm/app/compare/v4.4.0...v4.5.0
[4.4.0]: https://github.com/giantswarm/app/compare/v4.3.0...v4.4.0
[4.3.0]: https://github.com/giantswarm/app/compare/v4.2.0...v4.3.0
[4.2.0]: https://github.com/giantswarm/app/compare/v4.1.0...v4.2.0
[4.1.0]: https://github.com/giantswarm/app/compare/v4.0.0...v4.1.0
[4.0.0]: https://github.com/giantswarm/app/compare/v3.7.0...v4.0.0
[3.7.0]: https://github.com/giantswarm/app/compare/v3.5.0...v3.7.0
[3.6.0]: https://github.com/giantswarm/app/compare/v3.5.0...v3.6.0
[3.5.0]: https://github.com/giantswarm/app/compare/v3.4.0...v3.5.0
[3.4.0]: https://github.com/giantswarm/app/compare/v3.3.0...v3.4.0
[3.3.0]: https://github.com/giantswarm/app/compare/v3.2.0...v3.3.0
[3.2.0]: https://github.com/giantswarm/app/compare/v3.1.1...v3.2.0
[3.1.1]: https://github.com/giantswarm/app/compare/v3.1.0...v3.1.1
[3.1.0]: https://github.com/giantswarm/app/compare/v3.0.0...v3.1.0
[3.0.0]: https://github.com/giantswarm/app/compare/v2.0.0...v3.0.0
[2.0.0]: https://github.com/giantswarm/app/compare/v0.2.3...v2.0.0
[0.2.3]: https://github.com/giantswarm/app/compare/v0.2.2...v0.2.3
[0.2.2]: https://github.com/giantswarm/app/compare/v0.2.1...v0.2.2
[0.2.1]: https://github.com/giantswarm/app/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/giantswarm/app/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/giantswarm/app/releases/tag/v0.1.0
