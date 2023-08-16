# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2023-16-Aug

### Added
- Labels for service and pod.
- Prometheus scrape config example configuration added to the readme.
- Resources requests and limits in values by default.
- readOnlyRootFilesystem set to true on securityContext.
- runAsNonRoot set to true on securityContext.
- runAsUser set to true on securityContext.
### Changed
- Changed service monitor port name to http.
- Changed service monitor selector.
- Liveness and readiness probes configuration moved from template to values.
- query-exporter docker image version changed from 2.8.0 to 2.8.1.=

### Removed
- Removed latest as a default docker image.

## [0.1.0] - 2022-03-30

### Added
- Initial Helm Chart with basic functionallity


