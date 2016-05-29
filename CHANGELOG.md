# Change Log

All notable changes to this project will be documented in this file.
This project adheres to [Semantic Versioning](http://semver.org/).

## [Unreleased]

### Added
- Add a roadmap to the README.md

### Fixed
- The staticSpamDomains provider did not respect HTTP status codes.

## v0.4.0

### Added
- Add logging for the update command
- Add package documentation
- Add documentation for public functions
- Introduce a new list-spam-domains action

### Changed
- Combine static spam domains with dynamic ones
- Combine multiple referrer spam domain sources
- Display status as a percentage instead of a text-based status
- Allow to specify the number of days for the find-spam action
- Change the filter names prefix to "Referrer Spam Block Segment"

### Removed
- Remove the global status ... it didn't make much sense

### Fixed
- Fixed the update command. Updates did not work before.
- Fix template newline handling
- Fix the filesystem token store. Create the directory if it does not exist.

## v0.3.0

### Added

- Add `detect-spam` action which analyzes the last 30 days of analytics data for referrer spam

## v0.2.0

### Changed

- Add `accountID` argument to **update** and **remove** actions. So updating and removing will only be performed on specified accounts.
- Add optional `accountID` argument to **status** command to allow targeted spam-control status checks.
- Add optional `--quiet` or `-q` flag to **status** command to print the status overview in a format that can be parsed easier by tools like awk

## v0.1.0

### Added

- Add basic command line interface for displaying the spam-control status, updating the spam-control filters and removing the spam-control features. The authentication is done via oAuth2.
