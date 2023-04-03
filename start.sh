#!/bin/sh



if [ -e pid.txt ]
then
  echo "Already running"
  exit 1
else
  nohup ./bikep >> log.txt &
  echo $! > pid.txt
fi