a leveldb wrapper for levigo

simplify use leveldb in go

# Build leveldb

see [https://gist.github.com/siddontang/dfbc835e06e47d0f6297](https://gist.github.com/siddontang/dfbc835e06e47d0f6297) for build leveldb

# Install

you must first set CGO_CFLAGS, CGO_LDFLAGS to your leveldb and snappy directory.

dev.sh may help you:

    . ./dev.sh

# Notice

I have changed this package to [https://github.com/siddontang/go-leveldb](https://github.com/siddontang/go-leveldb) and will not maintain here anymore.