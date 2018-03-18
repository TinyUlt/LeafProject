#!/bin/sh


echo "build proto"
SRC_DIR=/Users/tinyult/go/src/github.com/name5566/protobuf
DST_DIR=/Users/tinyult/go/src/github.com/name5566/leafserver/src/server/msg
DST_DIR2=/Users/tinyult/go/src/github.com/name5566/leafclient/src/client/msg
protoc -I=$SRC_DIR --go_out=$DST_DIR $SRC_DIR/*.proto
protoc -I=$SRC_DIR --go_out=$DST_DIR2 $SRC_DIR/*.proto
echo "done"