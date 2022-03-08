#!/bin/bash
smallestdroplet=s-1vcpu-1gb
london=lon1
name=$1
dockerimg=docker-20-04

# Turning lines into comma separted list: https://stackoverflow.com/a/17759537
sshkeys=$(doctl compute ssh-key list --format ID --no-header true | xargs | sed -e 's/ /,/g')
ip=$(doctl compute droplet create $name --size $smallestdroplet --region $london --image $dockerimg --ssh-keys $sshkeys --wait --format PublicIPv4 --no-header)
echo Droplet can now be accessed at $ip 