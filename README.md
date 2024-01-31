# Usage

## As a command

```sh
$ go install go.katupy.io/ugen/cmd/ugen@latest
$ ~/go/bin/ugen
D5XlMsS2e5Xb
$ ~/go/bin/ugen -h
Usage:
  ugen [flags]

Flags:
  -a, --any                Use any character from the random generator.
      --base64             Output base64 encoded strings. Sets --any.
  -c, --count int          The number of unique strings to generate. (default 1)
  -d, --digit              Use only digits. If length > 1, digit will never start with 0.
  -h, --help               help for ugen
      --hex                Output hex encoded strings. Sets --any.
  -i, --interval string    Generate a random number within the provided open ended, i.e., [_,_), interval. E.g., -1000,1000. If a single number is provided, the interval begins with zero. Ignores --length.
  -l, --length int         Length of the generated string. (default 12)
      --lower              Output characters in lower case.
      --prefix string      Write prefix before each generated string.
  -s, --separator string   Separator for generated strings. (default "\n")
      --suffix string      Write suffix after each generated string.
      --ulid               Generate a ULID.
      --ulid-as-uuid       Generate a ULID displayed as a UUID. Sets --ulid.
      --upper              Output characters in upper case.
      --uuid               Generate a random (v4) UUID.
      --uuid7              Generate a v7 UUID.
```

## As a module

```sh
go get -u go.katupy.io/ugen
```
