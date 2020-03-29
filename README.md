# csvtojson

A simple command line utility to parse csv files and convert them to JSON/YML.

## Installation

```
$ go get -u go.krishnaiyereaswaran.com/csvtojson
```

## Usage

```
csvtojson [flags]

Flags:
      --allow-malformed          allow parsing malformed CSV
  -c, --csv-file string          input csv file name
  -d, --debug                    print detailed logs for errors
      --fill-empty-with string   value to fill empty cells with. --allow-malformed must be set for this to be effective
  -h, --help                     help for csvtojson
  -o, --out-file string          output file name. Use yaml or json based on the required format.
      --yaml                     marshal to yaml instead of json
```


## License

The contents of this repository is licensed under the Apache 2.0 license. Please check the [LICENSE](./LICENSE) file for more information.