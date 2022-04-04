'#!/bin/bash

if [[ $# > 0 ]]; then

    echo "start copy single vm..."

    for f in `find /mnt/pvedata -name *$1*raw -type f`
    do
#	echo $f
#	dd if=/mnt/pvedata$f bs=2M of=/mnt/sync$f &
#       /usr/bin/rsync -e "ssh -c aes128-gcm@openssh.com -o Compression=no -x" -B 128k -a --inplace --no-whole-file $f node-agregt.sync:/$f &
        /opt/blcp/blcpc $f 4194304 node-agregt.sync $f &
    done

    exit 0
fi

echo "start copy all running vm..."
declare -a arr

for kvm in `/usr/sbin/qm list | grep running | awk '{print $1}'`
do
    for disk in `find /mnt/pvedata -name *$kvm*raw -type f`
    do
        #echo $disk
        arr+=(`echo $disk`)
    done
done

for lxc in `/usr/bin/lxc-ls --running`
do
    echo ""
done

for file in ${arr[@]}; do
    #echo $file
    #dd if=/mnt/pvedata$file bs=2M of=/mnt/sync$file &
    #/usr/bin/rsync -e "ssh -c aes128-gcm@openssh.com -o Compression=no -x" -B 128k -a --inplace --no-whole-file $file node-agregt.sync:/$file &
    /opt/blcp/blcpc $file 1572864 node-agregt.sync $file
done

#/usr/bin/rsync -e "ssh -c aes128-gcm@openssh.com -o Compression=no -x" -a --inplace --no-whole-file $images$kvm node1-engenr.sync://$images &
#/usr/bin/bbcp -r -f -k -w 180k -s 20 $images$kvm node1-engenr.sync:$images &

