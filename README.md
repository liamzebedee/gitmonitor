gitmonitor
==========

Scans $HOME and prints when Git repo's are changed. Example:

```
2018/01/30 00:07:49 Done
2018/01/30 00:07:55 event: "/Users/liamz/Documents/electron-quick-start/.git/index": CHMOD
2018/01/30 00:07:55 event: "/Users/liamz/Documents/electron-quick-start/.git/index.lock": CREATE
2018/01/30 00:07:55 event: "/Users/liamz/Documents/electron-quick-start/.git/index.lock": REMOVE
```

MIT license.

## Install
 - `go build && ./gitmon` or
 - `go get github.com/liamzebedee/gitmonitor`

## Precompiled binaries
See `dist/`