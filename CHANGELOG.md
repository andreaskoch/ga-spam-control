# Change Log

All notable changes to this project will be documented in this file.
This project adheres to [Semantic Versioning](http://semver.org/).

## [Unreleased]

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
