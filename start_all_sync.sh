#!/bin/bash

systemctl start syncvm.service
sleep 5


declare -a synchost
cnt=0

cd /etc/pve/qemu-server

for vmid in `grep -r 'hooksync.sh' | sed 's/.conf/ /g' | awk '{print $1}'`
do

    synchost[$cnt]="$vmid"
    ((cnt++))

done

cd /root


for kvm in `/usr/sbin/qm list | grep running | awk '{print $1}'`
do

    if [[ "${synchost[*]}" =~ "${kvm}" ]]
    then
        #echo $kvm
        echo "$kvm:start:" | nc localhost 7011 -w 1
    fi
done
