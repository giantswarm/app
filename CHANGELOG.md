# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [8.0.0] - 2025-04-01

### Changed

- Changed `app.Validator` interface.
  - Replaced `projectName` string field with `isAdmissionController` boolean field.
  - Removed `enableManagedByLabel` field.
- When `app.Validator.isAdmissionController` is enabled, the existence checks for referenced config maps and secrets
  under App CRs (`.spec.config`, `.spec.userConfig`, `.spec.kubeConfig.secret`) will be skipped. This makes the
  existence of the `giantswarm.io/managed-by` label irrelevant as that was used to produce this behaviour that is
  the default behaviour now.

### Removed

- Removed `key.IsManagedBy` function as it is not relevant anymore.

## [7.1.0] - 2025-02-26

### Changed

- Make the `cordon-until` the only required annotation to cordon an app and make it expire automatically.

## [7.0.4] - 2025-01-30

### Added

- Add `ConfigMapExtraConfigs` and `SecretExtraConfigs` functions.

## [7.0.3] - 2025-01-22

## [7.0.2] - 2024-06-04

- Dependency updates

## [7.0.1] - 2024-04-15

- Dependency updates

## [7.0.0] - 2023-08-15

### Removed

- App: Remove `nginx-ingress-controller-app` exceptions. ([#294](https://github.com/giantswarm/app/pull/294))

## [6.15.6] - 2023-03-30

### Added

- Validation for unique App Operator version label for in-cluster apps.

## [6.15.5] - 2023-03-10

### Fixed

- Fix wrong field selector.

## [6.15.4] - 2023-03-09

### Fixed

- Getting the cluster ID for all the clusters types.

## [6.15.3] - 2023-03-08

### Changed

- Account for ClusterSingleton and NamespaceSingleton for CAPI clusters.

## [6.15.2] - 2023-02-01

## [6.15.1] - 2022-11-17

### Changed

- Adapt ClusterSingleton for CAPI.

## [6.15.0] - 2022-09-22

## Added

- Adding support for timeout fields in the App CR.

## [6.14.0] - 2022-09-13

### Changed

- Don't look for the `ingress-controller-values` ConfigMap cluster values when NGINX apps is installed for CAPI clusters. Instead look for the standard `<cluster_name>-cluster-values` ConfigMap.

## [6.13.0] - 2022-08-25

### Added

- Add a new validation rule for App Crs that checks for unique name of in-cluster Apps across the management cluster and rejects the request if it does not pass

## [6.12.0] - 2022-07-08

### Added

- Implement merging algorithm of multi layer App CR extra configuration layers RFC, see: https://github.com/giantswarm/rfc/tree/main/multi-layer-app-config#enhancing-app-cr

## [6.11.1] - 2022-06-15

### Removed

- Revert changes made in v6.11.0

## [6.11.0] - 2022-06-14

### Changed

- Use `application.giantswarm.io` CRDs from `apiextensions-application` instead of `apiextensions`.

### Added

- Add additional validation logic for the `chart-operator.giantswarm.io/app-namespace` annotation.

## [6.9.0] - 2022-03-11

### Added

- Add `key.AppCatalogEntryCompatibleProviders`.

## [6.8.1] - 2022-02-25

### Fixed

- Remove compatible providers validation for `AppCatalogEntry` as its overly strict.

## [6.8.0] - 2022-02-21

### Changed

- Downgrade k8s modules to `< 0.21.0` version and controller-runtime to `< 0.7.0` version.

## [6.7.0] - 2022-02-17

### Added

- Support setting extra labels and annotations when creating new App CR with app package.

## [6.6.2] - 2022-02-09

### Changed

- Make validation of `giantswarm.io/managedby` more general.

## [6.6.1] - 2022-01-26

### Fixed

- Fix key cluster methods to support org-namespaced App CRs.
- Fix validation order, so that labels are checked before App CR configs.

## [6.6.0] - 2022-01-24

### Added

- Add support for validating `giantswarm.io/cluster` label for org-namespaced App CRs.

## [6.5.1] - 2022-01-20

### Fixed

- Rename `ValidateResourcesExist` flag to `EnableManagedByLabel`.

## [6.5.0] - 2022-01-19

### Added

- Add `validateResourcesExist` flag to Validator for enabling/disabling support for `key.IsManagedByFlux`.

### Changed

- Narrow down the `key.IsManagedByFlux` impact on validation. It now skips only ConfigMap and Secret existence checks.

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

[Unreleased]: https://github.com/giantswarm/app/compare/v8.0.0...HEAD
[8.0.0]: https://github.com/giantswarm/app/compare/v7.1.0...v8.0.0
[7.1.0]: https://github.com/giantswarm/app/compare/v7.0.4...v7.1.0
[7.0.4]: https://github.com/giantswarm/app/compare/v7.0.3...v7.0.4
[7.0.3]: https://github.com/giantswarm/app/compare/v7.0.2...v7.0.3
[7.0.2]: https://github.com/giantswarm/app/compare/v7.0.1...v7.0.2
[7.0.1]: https://github.com/giantswarm/app/compare/v7.0.0...v7.0.1
[7.0.0]: https://github.com/giantswarm/app/compare/v6.15.6...v7.0.0
[6.15.6]: https://github.com/giantswarm/app/compare/v6.15.5...v6.15.6
[6.15.5]: https://github.com/giantswarm/app/compare/v6.15.4...v6.15.5
[6.15.4]: https://github.com/giantswarm/app/compare/v6.15.3...v6.15.4
[6.15.3]: https://github.com/giantswarm/app/compare/v6.15.2...v6.15.3
[6.15.2]: https://github.com/giantswarm/app/compare/v6.15.1...v6.15.2
[6.15.1]: https://github.com/giantswarm/app/compare/v6.15.0...v6.15.1
[6.15.0]: https://github.com/giantswarm/app/compare/v6.14.0...v6.15.0
[6.14.0]: https://github.com/giantswarm/app/compare/v6.13.0...v6.14.0
[6.13.0]: https://github.com/giantswarm/app/compare/v6.12.0...v6.13.0
[6.12.0]: https://github.com/giantswarm/app/compare/v6.11.1...v6.12.0
[6.11.1]: https://github.com/giantswarm/app/compare/v6.11.0...v6.11.1
[6.11.0]: https://github.com/giantswarm/app/compare/v6.10.0...v6.11.0
[6.10.0]: https://github.com/giantswarm/app/compare/v6.9.0...v6.10.0
[6.9.0]: https://github.com/giantswarm/app/compare/v6.8.1...v6.9.0
[6.8.1]: https://github.com/giantswarm/app/compare/v6.8.0...v6.8.1
[6.8.0]: https://github.com/giantswarm/app/compare/v6.7.0...v6.8.0
[6.7.0]: https://github.com/giantswarm/app/compare/v6.6.2...v6.7.0
[6.6.2]: https://github.com/giantswarm/app/compare/v6.6.1...v6.6.2
[6.6.1]: https://github.com/giantswarm/app/compare/v6.6.0...v6.6.1
[6.6.0]: https://github.com/giantswarm/app/compare/v6.5.1...v6.6.0
[6.5.1]: https://github.com/giantswarm/app/compare/v6.5.0...v6.5.1
[6.5.0]: https://github.com/giantswarm/app/compare/v6.4.0...v6.5.0
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
