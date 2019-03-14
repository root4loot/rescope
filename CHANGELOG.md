# Changelog

All notable changes to rescope will be documented here.
Unreleased changes will go in the below heading.

## [Unreleased]

## [0.2] - 2019-03-14
### Changed
- How Burp/ZAP is specified from the cli. Now using flags instead.
- Printing of identified targets. rescope will now include a leading +/- for each target to better indicate which are includes and excludes. Perhaps that'll make things clearer for those who are color blind.
- rescope will now handle !EXCLUDE tag having leading/preceding text on the same line.

### Fixed
- A filepath issue that resulted in panic upon parsing to ZAP from executable that resided outside of package directory.
- An issue that resulted in IP ranges/CIDR from being parsed correctly.
- Burp parsing inaccuracy for certain targets having both http(s) and ports. Rescope will now include protocol ports (80|443) and host port when parsing to Burp.
- Minor issue that prevented --version from being displayed.

### Added
- CHANGELOG.md

## [0.1] - 2019-03-06
### Added
- Repository to Github (https://github.com/root4loot/rescope)
- Blog Post: https://root4loot.com/post/introducing_rescope/
