## Decimal v1.4.5

- Rebase on upstream
- Require go 1.24 and update go deps

Previous releases contained merged pull requests.

## Decimal v1.4.1

- Fix BSON encoding/decoding not working with omitempty or pointer fields

## Decimal v1.4.0

Forked as https://github.com/greatcloak/decimal.

- Update readme with project status and fork info
- Refactor files/tests for readability and fix go linter errors
- Add BSON marshaler/unmarshaller
- Allow empty string when unmarshalling JSON

## Decimal v1.3.1

#### ENHANCEMENTS

- Reduce memory allocation in case of initialization from big.Int [#252](https://github.com/shopspring/decimal/pull/252)

#### BUGFIXES

- Fix binary marshalling of decimal zero value [#253](https://github.com/shopspring/decimal/pull/253)

## Decimal v1.3.0

#### FEATURES

- Add NewFromFormattedString initializer [#184](https://github.com/shopspring/decimal/pull/184)
- Add NewNullDecimal initializer [#234](https://github.com/shopspring/decimal/pull/234)
- Add implementation of natural exponent function (Taylor, Hull-Abraham) [#229](https://github.com/shopspring/decimal/pull/229)
- Add RoundUp, RoundDown, RoundCeil, RoundFloor methods [#196](https://github.com/shopspring/decimal/pull/196) [#202](https://github.com/shopspring/decimal/pull/202) [#220](https://github.com/shopspring/decimal/pull/220)
- Add XML support for NullDecimal [#192](https://github.com/shopspring/decimal/pull/192)
- Add IsInteger method [#179](https://github.com/shopspring/decimal/pull/179)
- Add Copy helper method [#123](https://github.com/shopspring/decimal/pull/123)
- Add InexactFloat64 helper method [#205](https://github.com/shopspring/decimal/pull/205)
- Add CoefficientInt64 helper method [#244](https://github.com/shopspring/decimal/pull/244)

#### ENHANCEMENTS

- Performance optimization of NewFromString init method [#198](https://github.com/shopspring/decimal/pull/198)
- Performance optimization of Abs and Round methods [#240](https://github.com/shopspring/decimal/pull/240)
- Additional tests (CI) for ppc64le architecture [#188](https://github.com/shopspring/decimal/pull/188)

#### BUGFIXES

- Fix rounding in FormatFloat fallback path (roundShortest method, fix taken from Go main repository) [#161](https://github.com/shopspring/decimal/pull/161)
- Add slice range checks to UnmarshalBinary method [#232](https://github.com/shopspring/decimal/pull/232)

## Decimal v1.2.0

#### BREAKING

- Drop support for Go version older than 1.7 [#172](https://github.com/shopspring/decimal/pull/172)

#### FEATURES

- Add NewFromInt and NewFromInt32 initializers [#72](https://github.com/shopspring/decimal/pull/72)
- Add support for Go modules [#157](https://github.com/shopspring/decimal/pull/157)
- Add BigInt, BigFloat helper methods [#171](https://github.com/shopspring/decimal/pull/171)

#### ENHANCEMENTS

- Memory usage optimization [#160](https://github.com/shopspring/decimal/pull/160)
- Updated travis CI golang versions [#156](https://github.com/shopspring/decimal/pull/156)
- Update documentation [#173](https://github.com/shopspring/decimal/pull/173)
- Improve code quality [#174](https://github.com/shopspring/decimal/pull/174)

#### BUGFIXES

- Revert remove insignificant digits [#159](https://github.com/shopspring/decimal/pull/159)
- Remove 15 interval for RoundCash [#166](https://github.com/shopspring/decimal/pull/166)
