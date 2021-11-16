#!/usr/bin/env bash

# Folder like this:
# ```
# + generator
#  - templates
#  - tto.conf
# - example.sh
# ```

GEN_DIR=$(dirname "$0")/generator

tto -conf "$GEN_DIR"/tto.conf -t "$GEN_DIR"/templates -o output -excludes QRT -clean

find output/java -name "*.java"  -print0 | xargs -0 sed -i 's/ int / Integer /g' && \
find output/java -name "*.java"  -print0 | xargs -0 sed -i 's/ long / Long /g' && \
find output/java -name "*.java" -print0 | xargs -0 sed -i 's/ float / Float /g'

#find platform-business/src -path "*/k2/entity" -print0 | \
#   xargs  -0 cp output/spring/src/main/java/com/github/system1/entity/*
exit 0;