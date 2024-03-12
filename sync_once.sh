#!/bin/bash

dir="/mnt/pvedata"
dst="node2-engenr.sync"

if [[ $# > 0 ]]; then

    for f in `find $dir -name *$1*raw -type f`
    do
       #echo $f
       /usr/bin/rsync -e "ssh -c aes128-gcm@openssh.com -o Compression=no -x" -B 128K -a --inplace --no-whole-file $f $dst:/$f
    done

fi


