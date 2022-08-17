#!/bin/bash

# If wm migrate, stop sync (copy) raw disk files from this node

while true
do

    count=0

    for line in `pgrep task -a | grep qmigrate | sed 's/:/ /g'`
#  for line in `cat test.migrate | sed 's/:/ /g'`
    do

        if [[ $count -eq 8 ]]
        then
            #echo "count: "$count", line: "$line
            /opt/syncvm/stop_sync.sh $line
        fi

        ((count++))

        if [[ $count -eq 10 ]]
        then
           count=0
        fi

    done

    sleep 10

done
