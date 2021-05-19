#!/bin/bash
set -euo pipefail

mkdir -p /tmp/vstore-test
for i in {1..2}; do
  ffmpeg -f lavfi -i testsrc=duration=10:size=1280x720:rate=30 /tmp/vstore-test/video$i.mpg
done
