# Change Log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/) and this
project adheres to [Semantic Versioning](https://semver.org/).

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
