#!/usr/bin/env bash

# put this script into /var/lib/vz/snippets
# add line to vm conf:
# hookscript: local:snippets/hooksync.sh

vmid=$1
phase=$2

if [ $phase == "pre-start" ]
then

    echo "pre-start $vmid"

fi

if [ $phase == "post-start" ]
then

    echo "post-start $vmid"

    # Send signal to syncvm daemon on port 7011 for run script (start_sync.sh $vmid)
    echo "$vmid:start:" | nc localhost 7011 -w 1

fi


if [ $phase == "pre-stop" ]
then

    echo "pre-stop $vmid"

fi


if [ $phase == "post-stop" ]
then

    echo "post-stop $vmid"

    # Stop sync
    echo "$vmid:stop:" | nc localhost 7011 -w 1

fi
