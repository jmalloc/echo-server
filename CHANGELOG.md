# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog], and this project adheres to
[Semantic Versioning].

<!-- references -->

[keep a changelog]: https://keepachangelog.com/en/1.0.0/
[semantic versioning]: https://semver.org/spec/v2.0.0.html

## [0.3.6] - 2023-10-31

- Add support for sending arbitrary headers in all responses

## [0.3.5] - 2023-03-07

- Update dependencies to patch some security vulnerabilites

## [0.3.4] - 2022-12-23

### Fixed

- Fix rendering of the HTTP request line to match RFC-2616

## [0.3.3] - 2022-07-11

### Fixed

- Prevent rendering empty header value along real header values
- Render all values for each header, not only the first

## [0.3.2] - 2022-07-08

### Changed

- Headers are now displayed in alphabetical order (thanks [@marcofranssen])

## [0.3.1] - 2021-09-12

### Added

- Add the `SEND_SERVER_HOSTNAME` environment variable, set to `false` to prevent the server from sending its hostname
- Add the `X-Send-Server-Hostname` request header

## [0.3.0] - 2021-08-16

### Added

- Add the `/.sse` route, which echoes the request using server-sent events

## [0.2.0] - 2021-06-03

### Added

- Add support for logging HTTP headers to stdout (thanks [@arulrajnet])

## [0.1.0] - 2020-04-20

- Initial release

<!-- references -->

[unreleased]: https://github.com/jmalloc/echo-server
[0.1.0]: https://github.com/jmalloc/echo-server/releases/v0.1.0
[0.2.0]: https://github.com/jmalloc/echo-server/releases/v0.2.0
[0.3.0]: https://github.com/jmalloc/echo-server/releases/v0.3.0
[0.3.1]: https://github.com/jmalloc/echo-server/releases/v0.3.1
[0.3.2]: https://github.com/jmalloc/echo-server/releases/v0.3.2
[0.3.4]: https://github.com/jmalloc/echo-server/releases/v0.3.4
[0.3.5]: https://github.com/jmalloc/echo-server/releases/v0.3.5

<!-- outside contributors -->

[@arulrajnet]: https://github.com/arulrajnet
[@marcofranssen]: https://github.com/marcofranssen

<!-- version template
## [0.0.1] - YYYY-MM-DD

### Added
### Changed
### Deprecated
### Removed
### Fixed
### Security
-->
