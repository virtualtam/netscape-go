# Change Log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/) and this
project adheres to [Semantic Versioning](https://semver.org/).

## Unreleased
### Changed

- Require Go 1.25
- Modernize Go code

## [v2.3.0](https://github.com/virtualtam/netscape-go/releases/tag/v2.3.0) - 2025-02-23
### Added

- Add Go fuzzing corpus
- Add Make target to seed the fuzzing corpus from test input data
- Setup Copywrite to manage license headers

### Changed

- Require Go 1.24
- Update CI workflow
- Fix nil pointer dereference issues detected by fuzzing
- Relocate benchmark input to `testdata/benchmark/`
- Relocate test input to `testdata/input/`


## [v2.2.0](https://github.com/virtualtam/netscape-go/releases/tag/v2.2.0) - 2023-10-26
### Changed

- Require Go 1.21
- Update CI workflows
- Marshal / encoder: HTML-escape bookmark description
- Unmarshal / decoder: HTML-unescape bookmark description


## [v2.1.0](https://github.com/virtualtam/netscape-go/releases/tag/v2.1.0) - 2022-06-03
### Changed

- Do not use time.Time pointers to represent nullable dates
- Introduce JSON marshaling methods with intermediary structs to handle nullable
  dates (via time.Time.IsZero())

### Fixed

- Initialize Decoder for proper UNIX timestamp format detection


## [v2.0.0](https://github.com/virtualtam/netscape-go/releases/tag/v2.0.0) - 2022-04-30
### Added

- Add a method to flatten a Document by returning a version with all bookmarks
  attached to the root folder

### Changed

- Bump Go module to v2
- Flatten and cleanup package structure to ease importing and using as a library


## [v1.1.0](https://github.com/virtualtam/netscape-go/releases/tag/v1.1.0) - 2022-03-25
### Added

- Add Make targets to run benchmarks and save CPU/Memory profiles

### Changed

- decoder: refactor date/time operations
- Use Markdown for the README

### Fixed

- encoder: sort folder and bookmark attributes by key for deterministic output


## [v1.0.0](https://github.com/virtualtam/netscape-go/releases/tag/v1.0.0) - 2022-03-06

Initial release.

### Added

- Unmarshal data using the Netscape Bookmark file format
- Marshal documents containing bookmarks and folders using the Netscape Bookmark
  file format
- Add support for nested folders
- Add support for folder metadata:

  - creation and update dates
  - arbitrary attributes
  - text description, with multi-line and inner markup support

- Add support for bookmark metadata:

  - creation and update dates
  - visibility
  - comma-separated tags
  - arbitrary attributes
  - text description, with multi-line and inner markup support

- Provide code and command-line examples to demonstrate usage
