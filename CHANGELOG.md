# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.9.0] - 2026-09-14
### Added
- Security instructions in docs (#264)

### Changed
- Run debugger and monitor as root (needed for UDP) (#263)

## [0.8.37] - 2026-09-13
### Added
- Documentation website with VitePress, published to GitHub Pages (#245, #248, #249)

### Changed
- Unified release process into one tag-driven release (#252)
- Consolidated podNetworkingV2 into podNetworking (#247)
- Updated CRD memory limit customization (#246)

### Fixed
- Docs favicon and hero wording (#250, #251)

## [0.8.36] - 2026-09-11
### Fixed
- Build Vite stage on BUILDPLATFORM to avoid QEMU npm crash

## [0.8.35] - 2026-09-11
### Fixed
- Commit version bump reliably and create a GitHub release each deploy
