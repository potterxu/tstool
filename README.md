# tstool

personal tool for ts analysis

## stat
```
statistis of ts stream

Usage:
tstool stat (interface ipAddress port) [flags]

Flags:
  -h, --help           help for stat
  -i, --interval int   Statistics update interval in seconds (default 3)
```

## cap
```
capture ts multicast

Usage:
tstool cap (interface ipAddress port) [flags]

Flags:
  -h, --help            help for cap
  -o, --output string   Output filename (default "./output.ts")
  -t, --time int        Captureing time
```

## parse

### psi
```
parse psi

Usage:
  tstool parse psi (filename) [flags]

Flags:
  -h, --help   help for psi
```
