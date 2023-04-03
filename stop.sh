#!/bin/sh

if [ -e pid.txt ]
then
  kill -2 `cat pid.txt`
  rm pid.txt
else
  echo "Not running"
  exit 1
fi