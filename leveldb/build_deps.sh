#!/bin/bash

#refer https://github.com/norton/lets/blob/master/c_src/build_deps.sh

#install leveldb and snappy

#you must set your own snappy and leveldb source directory
#snappy https://drive.google.com/file/d/0B0xs9kK-b5nMOWIxWGJhMXd6aGs/edit?usp=sharing
#leveldb https://leveldb.googlecode.com/files/leveldb-1.15.0.tar.gz
SNAPPY_SRC=./snappy
LEVELDB_SRC=./leveldb

SNAPPY_DIR=/usr/local/snappy
LEVELDB_DIR=/usr/local/leveldb

if [ ! -f $SNAPPY_DIR/lib/libsnappy.a ]; then
    (cd $SNAPPY_SRC && \
        ./configure --prefix=$SNAPPY_DIR && \
        make && \
        make install)
else
    echo "skip install snappy"
fi

if [ ! -f $LEVELDB_DIR/lib/libleveldb.a ]; then
    (cd $LEVELDB_SRC && \
        echo "echo \"PLATFORM_CFLAGS+=-I$SNAPPY_DIR/include\" >> build_config.mk" >> build_detect_platform &&
        echo "echo \"PLATFORM_CXXFLAGS+=-I$SNAPPY_DIR/include\" >> build_config.mk" >> build_detect_platform &&
        echo "echo \"PLATFORM_LDFLAGS+=-L $SNAPPY_DIR/lib -lsnappy\" >> build_config.mk" >> build_detect_platform &&
        make SNAPPY=1 && \
        make && \
        mkdir -p $LEVELDB_DIR/include/leveldb && \
        install include/leveldb/*.h $LEVELDB_DIR/include/leveldb && \
        mkdir -p $LEVELDB_DIR/lib && \
        cp -f libleveldb.* $LEVELDB_DIR/lib)
else
    echo "skip install leveldb"
fi

function add_path()
{
  # $1 path variable
  # $2 path to add
  if [ -d "$2" ] && [[ ":$1:" != *":$2:"* ]]; then
    echo "$1:$2"
  else
    echo "$1"
  fi
}

export CGO_CFLAGS="-I$LEVELDB_DIR/include -I$SNAPPY_DIR/include"
export CGO_LDFLAGS="-L$LEVELDB_DIR/lib -L$SNAPPY_DIR/lib -lsnappy"
export LD_LIBRARY_PATH=$(add_path $LD_LIBRARY_PATH $SNAPPY_DIR/lib)
export LD_LIBRARY_PATH=$(add_path $LD_LIBRARY_PATH $LEVELDB_DIR/lib)


go get github.com/jmhodges/levigo 
