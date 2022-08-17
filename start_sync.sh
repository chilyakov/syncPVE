#!/bin/bash

if [[ $# > 0 ]]; then

    while true
    do
        for f in `find /mnt/pvedata -name *$1*raw -type f`
        do
           #echo $f
           /usr/bin/rsync -e "ssh -c aes128-gcm@openssh.com -o Compression=no -x" -B 128K -a --inplace --no-whole-file $f node-engenr.sync:/$f
        done

        sleep 60
    done

fi


