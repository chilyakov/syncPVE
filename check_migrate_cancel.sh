#1/bin/bash

# vm migrate to node2-engenr

tail -fn0 /var/log/pve/tasks/index | while read line
do

    echo $line | grep qmigrate | grep -v OK
    if [ $? = 0 ]
    then
        IFS=':'
        read -ra ADDR <<< "$line"
        echo ${ADDR[6]}":stop:" | nc node2-engenr 7011 -w 1
        echo ${ADDR[6]}":start:" | nc localhost 7011 -w 1
        #/opt/syncvm/stop_sync.sh ${ADDR[6]}
    fi

done
