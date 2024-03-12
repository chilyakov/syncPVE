#!/bin/bash

tmpdir="/mnt/pvedata/sas4"
dstdir="/mnt/fsbackup/`hostname -a`"

declare -a exclud=(142,144)

for kvm in `/usr/sbin/qm list | grep running | awk '{print $1}'`
do

if [[ ! "${exclud[*]}" =~ "${kvm}" ]]
then

    for dir in `find /mnt/pvedata -name $kvm -type d`
    do
	for file in `ls $dir | grep disk`
	do
            if [[ "$file" == 'vm-111-disk-1.raw' ]]; then
                continue
            fi

            if [[ "$dir/$file" == '/mnt/pvedata/sas4/images/143/vm-143-disk-1.raw' ]]; then
                continue
            fi

            if [[ "$dir/$file" == '/mnt/pvedata/sas4/images/141/vm-141-disk-0.raw' ]]; then
                continue
            fi

            if [[ "$dir/$file" == '/mnt/pvedata/sas4/images/135/vm-135-disk-0.qcow2' ]]; then
                continue
            fi

            if [[ "$dir/$file" == '/mnt/pvedata/sas4/images/115/vm-115-disk-0.raw' ]]; then
                continue
            fi



	    #echo $dir/$file
            dname=`echo $dir | sed 's/\// /g' | awk '{print $3}'`
            #echo $dname-$file

# without compress for kvm 111
            if [[ "$kvm" == '111' ]]; then

                cp -f $dir/$file $dstdir/$dname
                if [[ $(date +%u) -eq 6 ]]; then cp -f $dir/$file $dstdir/week/$dname; fi
                continue
            fi

# with compress
	    dd if=$dir/$file bs=512K | zstd -1 -T0 -f -o $tmpdir/$file.zst 2>/dev/null

            cp -f $tmpdir/$file.zst $dstdir/$dname
            if [[ $(date +%u) -eq 6 ]]; then cp -f $tmpdir/$file.zst $dstdir/week/$dname; fi
            rm -f $tmpdir/$file.zst
	done
    done

fi
done

#exit

#for f in `find /mnt/pvedata/sas2/ -name *.zst`
#do
#     cp -f $f /mnt/fsbackup/172.18.18.30
#    cp -f $f /mnt/synology/172.18.18.13
#done

#rm -f /mnt/pvedata/sas2/*.zst


