# JSONDIFF

A simple little tool that produces readable diff of 2 JSON-able (read "convertible to `map[string]interface{}`") objects. Useful for diagnostics or debugging

## Installation

```
go get github.com/elgris/jsondiff
```

## Examples of the output

## Limitation

- Coloured output tested with `bash` only, not sure how it will behave with other terminals.
- The tool converts input data into `map[string]interface{}` with json encoding/decoding. Hence, types of input map will change during unmarshal step: integers become float64 and so on (check https://golang.org/pkg/encoding/json/ for details).

## License

MIT
