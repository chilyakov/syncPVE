#1/bin/bash

# vm migrate to node-engenr
# on start migration sync stopped on node-agregt and start on node-engenr
# if migrate is cancel, sync must be stopped on node-engenr and start on node-agregt again

tail -fn0 /var/log/pve/tasks/index | while read line
do

    echo $line | grep qmigrate | grep -v OK
    if [ $? = 0 ]
    then
        IFS=':'
        read -ra ADDR <<< "$line"
        echo ${ADDR[6]}":stop:" | nc node-engenr 7011 -w 1
        echo ${ADDR[6]}":start:" | nc localhost 7011 -w 1
        #/opt/syncvm/stop_sync.sh ${ADDR[6]}
    fi

done
