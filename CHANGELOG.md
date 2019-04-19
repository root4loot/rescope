# Changelog

All notable changes to rescope will be documented here.
Unreleased changes will go in the below heading.

## [Unreleased]
### Fixed
- #3: Bug that caused fatal exit upon providing full intigriti program URL
- #2: Bug that caused duplicate scope definitions after parsing more than one program from one of the same affected services; hackerone, intigriti, yeswehack.

## [1.0] - 2019-04-16

### Added
- New feature that makes it possible to parse scopes directly from public bugbounty programs.
- New flag (-u | --url) related to the above.
- Option to specify includes (aside from excludes) with the optional !INCLUDE tag.
- New flag (--itag) related to the above.
- Gopkg files for dep management.

### Fixed
- Minor bug that caused certain number formats in string to match as valid targets.
- Bug that prevented targets with ports from being set in Zap context.
- Bug that prevented targets with leading wildcard and no immediate dot from being fully matched.

### Changed
- Replaced the previous vendor package 'github.com/fatih/color' with 'github.com/gookit/color' for adding colors support as the former project was no longer maintained.
- General commenting and code impovements.
- Changed --extag to -etag and removed the short version.


## [0.3] - 2019-03-30
### Improvements
- Cleaner project structure. Packages now live in `internal/` rather than project root.

### Added
- File `configs/services` which lists a bunch of service names and ports. With this, rescope is able to identify ports for targets that has schemes but no port specified. For now this is used when parsing to Burp. Example:
     - `ftps://example.com` => `port: 990`
     - `https://example.com:21` => `port: 443,21`

- Port 80,443 to Burp scope when identifier has no scheme and no port. This'll prevent people from touching ports that're otherwise (not clearly defined) out of bounds. Example:
     - `example.com` => `port: 80,443`

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
