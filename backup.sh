#!/bin/bash

tmpdir="/mnt/vz"
dstdir="/mnt/fsbackup/`hostname -a`"

declare -a exclud=()

echo `date +%d.%m.%Y" "%T` "backup is starting..." > backup.log

for kvm in `/usr/sbin/qm list | grep running | awk '{print $1}'`
do

if [[ ! "${exclud[*]}" =~ "${kvm}" ]]
then

    for dir in `find /mnt/pvedata -name $kvm -type d`
    do
        for file in `ls $dir | grep disk`
        do
            #echo $dir/$file

           check=$(grep -c $file backup.exclude) #skip VM disks from file
           if [ $check -ne 0 ]; then
              echo "skip $dir/$file" >> backup.log
              continue
           fi

            dname=`echo $dir | sed 's/\// /g' | awk '{print $3}'`
            #echo $dname-$file
            dd if=$dir/$file bs=512K | zstd -1 -T0 -f -o $tmpdir/$file.zst 2>/dev/null

            cp -f $tmpdir/$file.zst $dstdir/$dname
            if [[ $(date +%u) -eq 6 ]]; then cp -f $tmpdir/$file.zst $dstdir/week/$dname; fi
            rm -f $tmpdir/$file.zst
            echo "backup success $dir/$file" >> backup.log
        done
    done

fi
done

echo `date +%d.%m.%Y" "%T` "backup complete" >> backup.log



