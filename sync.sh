#!/usr/bin/bash


function sync() {
    rsync -avz --recursive -R --delete --prune-empty-dirs --exclude '.git' --info=progress2 -e 'ssh -p22 -o "StrictHostKeyChecking=no"' /home/krasi/src/github.com/arribada/smart-camera  root@192.168.0.116:/
}

sync

while inotifywait -r -e modify,create,delete /home/krasi/src/github.com/arribada/smart-camera
do
    sync
done 