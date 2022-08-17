#!/bin/bash

if [[ $# > 0 ]]; then

     for proc in `ps fx | grep "start_sync.sh" | grep $1 | awk '{print $1}'`
     do
        pkill -P $proc && kill $proc
        pkill -P $proc && kill $proc
     done

fi
