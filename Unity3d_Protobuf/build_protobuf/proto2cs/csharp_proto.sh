#!/bin/sh


echo "build proto"
SRC_DIR=/Users/tinyult/go/src/github.com/name5566/protobuf
DST_DIR=/Users/tinyult/go/src/github.com/name5566/leafserver/src/server/msg

protoc -I=$SRC_DIR --csharp_out=$DST_DIR $SRC_DIR/*.proto

echo "done"