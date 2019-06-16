# Changelog

All notable changes to rescope will be documented here.
Unreleased changes will go in the below heading.

## [Unreleased]

### Added
- Some unit test functions

### Fixed
- bugbounty.jp with missing scheme will no longer cause segfault.
- Scopes having avoided or conflicted targets on the last line should no longer cause out-of-bounds in removing them.
- Rare occurence where identifiers ending in `.*` extension, that also has multiple wildcards in domain did not parse correctly.

## [2.0] - 2019-06-26

### Added
- Support for [bugbounty.jp](https://bugbounty.jp)
- Support for [federacy.com](https://www.federacy.com/)
- New flag `--raw` that outputs naked (in-scope) definitions to file. Useful in working with other tools and programs.
- Support for resolving identifiers that conflict (overlap) with wilcarded excludes. Such conflict renders affected "in-scope" targets ineffective as excludes are prioritized in Burp/ZAP. This occurs when scopes are not properly defined, which if fairly common to see in BaaS programs.
- Support for avoiding certain third party resources, such as github.com, itunes.apple.com, play.google.com, etc, which is generally not something you want to scan/spider. Services are defined in [configs/avoid.txt](configs/avoid.txt). Met condition will prompt user as to whether affected targets should be ignored.


### Fixed
- Missing scopename prompt when parsing to ZAP without `--name` specified
- Targets like `www.*.example.com` and `*.*.example.com` should now parse correctly.
- Single IP's should now identify correctly.
- Bug that resulted in IP-ranges from being identified.

### Improved
- List handling to free up some unnecessary looping and improve extensibility.


## [1.1] - 2019-04-22

### Added
- Support for matching targets with s3 schema

### Fixed
- [#5](https://github.com/root4loot/rescope/issues/5) Targets separated by comma should now be grabbed correctly.  
- [#4](https://github.com/root4loot/rescope/issues/4) Intigriti programs should no longer parse with duplicate exclude definitions
- [#3](https://github.com/root4loot/rescope/issues/3) Bug that caused fatal exit upon providing full intigriti program URL
- [#2](https://github.com/root4loot/rescope/issues/2) Bug that caused duplicate scope definitions after parsing more than one program from one of the same affected services; hackerone, intigriti, yeswehack.
- Improper regex matching leading to strings having numbers and slashes to be matched as CIDR
- Wildcarded subdomains not parsing to Burp correctly

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
