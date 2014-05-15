a leveldb wrapper for levigo

simplify use leveldb in go

# Build leveldb

see [https://gist.github.com/siddontang/dfbc835e06e47d0f6297](https://gist.github.com/siddontang/dfbc835e06e47d0f6297) for build leveldb

# Install

you must first set CGO_CFLAGS, CGO_LDFLAGS to your leveldb and snappy directory.

dev.sh may help you:

    . ./dev.sh

# Performance

for better performance, I change some leveldb configurations to build:

+ db/dbformat.h

        // static const int kL0_SlowdownWritesTrigger = 8;
        static const int kL0_SlowdownWritesTrigger = 16;

        // static const int kL0_StopWritesTrigger = 12;
        static const int kL0_StopWritesTrigger = 64;

+ db/version_set.cc

        //static const int kTargetFileSize = 2 * 1048576;
        static const int kTargetFileSize = 32 * 1048576;

        //static const int64_t kMaxGrandParentOverlapBytes = 10 * kTargetFileSize;
        static const int64_t kMaxGrandParentOverlapBytes = 20 * kTargetFileSize;


