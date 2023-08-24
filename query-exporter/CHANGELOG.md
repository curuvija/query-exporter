# Changelog

This changelog contains only the last changes. Check https://github.com/curuvija/helm-charts/releases for each release.

## [2.0.0] - 2023-08-24

### Changed

- default docker image tag set to 2.9.0
- configuration loaded from secret instead of from configmap
- fixed yaml inconsistencies in deployment template

### Removed

- **Breaking:** removed configmap configuration in favor of more secure configuration from secret (check README for more details)
