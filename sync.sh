#!/usr/bin/bash


function sync() {
    rsync -avz --recursive -R --delete --prune-empty-dirs --exclude '.git' --info=progress2 -e 'ssh -p22 -o "StrictHostKeyChecking=no"' /home/$USER/src/github.com/arribada/smart-camera  root@raspberrypi.local:/
    rsync -avz --recursive -R --delete --prune-empty-dirs --exclude '.git' --info=progress2 -e 'ssh -p22 -o "StrictHostKeyChecking=no"' /home/$USER/src/github.com/groupgets/pylepton  root@raspberrypi.local:/
}

sync

while inotifywait -r -e modify,create,delete /home/$USER/src/github.com/arribada/smart-camera
do
    sync
done 