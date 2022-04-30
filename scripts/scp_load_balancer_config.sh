#!/bin/bash

# load balancer config file
filebeat_config='filebeat.yml'
prom_config='prometheus.yml'

# ssh key
key_file='ssh_key/terraform'

# ugly list concatenating of ips from terraform output
rows=$(terraform output -raw minitwit-swarm-leader-ip-address)
rows+=' '
rows+=$(terraform output -json minitwit-swarm-manager-ip-address | jq -r .[])
rows+=' '
rows+=$(terraform output -json minitwit-swarm-worker-ip-address | jq -r .[])

# scp the file
for ip in $rows; do
    scp -o "StrictHostKeyChecking no" -i $key_file $filebeat_config root@$ip:/root/filebeat.yml
    scp -o "StrictHostKeyChecking no" -i $key_file $prom_config root@$ip:/root/prometheus.yml
done
