# YASD

Yet Another Salesforce Dataloader

## Install

```bash
$ go get github.com/tzmfreedom/yasd
```

## Usage

Export records
```bash
$ yasd export -t {Salesforce Object Name}
```

Insert records
```bash
$ yasd insert -t {Salesforce Object Name} -f {path to source file} [--mapping {path to mapping file}] [--insert-nulls]
```

Update records
```bash
$ yasd update -t {Salesforce Object Name} -f {path to source file} [--mapping {path to mapping file}] [--insert-nulls]
```

Upsert records
```bash
$ yasd upsert -t {Salesforce Object Name} -f {path to source file} [--mapping {path to mapping file}] [--insert-nulls]
```

Delete records
```bash
$ yasd delete -t {Salesforce Object Name} -f {path to source file} [--mapping {path to mapping file}]
```

Undelete records
```bash
$ yasd undelete -t {Salesforce Object Name} -f {path to source file} [--mapping {path to mapping file}]
```

#### Common Option

* --username, -u

* --password, -p

* --endpoint, -e

* --api-version
  Specify Salesforce API Version (e.g. 43.0)

* --delimiter
  Specify CSV delimiter

* --encoding
  Specify CSV encoding, read/write

* --mapping
  Specify CSV header mapping file path

* --debug, -d
  If you set debug, cli output transmitting API SOAP XML to stdout.


#### DML SubCommand Option (Insert, Update, Upsert, Delete, Undelete)

* --file, -f

* --type, -t

* --success-file

* --error-file

* --query, -q

* --output, -o

* --format

* --batch-size


## Contribute

Just send pull request if needed or fill an issue!

## License

The MIT License See [LICENSE](https://github.com/tzmfreedom/yasd/blob/master/LICENSE) file.

