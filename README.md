[![Build Status](https://travis-ci.org/tzmfreedom/yasd.svg?branch=master)](https://travis-ci.org/tzmfreedom/yasd)

# YASD

Yet Another Salesforce Dataloader

## Install

For Windows user with cmd.exe
```
@"%SystemRoot%\System32\WindowsPowerShell\v1.0\powershell.exe" -NoProfile -InputFormat None -ExecutionPolicy Bypass -Command "iex ((New-Object System.Net.WebClient).DownloadString('http://install.freedom-man.com/yasd.ps1'))"
```
For Windows user with cmd.exe
```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force; iex ((New-Object System.Net.WebClient).DownloadString('http://install.freedom-man.com/yasd.ps1'))
```

For Linux, MacOS user
```bash
$ curl -sL http://install.freedom-man.com/yasd | bash
```

For golang user
```bash
$ go get github.com/tzmfreedom/yasd
```

## Usage

Export records
```bash
$ yasd export -q {SOQL}
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

