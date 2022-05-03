#!/bin/bash

echo -e "\n--> Checking that all the necessary files exist\n"
filebeat_config="filebeat.yml"
prom_config="prometheus.yml"
env_config=".env"
docker_config="docker-compose.yaml"

docker_dir="docker"
src_dir="src"

[ ! -f "$docker_config" ] && echo "$docker_config file does not exist" && exit
[ ! -f "$filebeat_config" ] && echo "$filebeat_config does not exist" && exit
[ ! -f "$prom_config" ] && echo "$prom_config does not exist" && exit
[ ! -f "$env_config" ] && echo "$env_config does not exist" && exit
[ ! -d "$docker_dir" ] && echo "$docker_dir directory does not exist" && exit
[ ! -d "$src_dir" ] && echo "$src_dir directory does not exist" && exit

echo -e "--> Provisioning a Droplet, this may take a couple of minutes...\n"
fourgdroplet=s-2vcpu-4gb
london=lon1
name=$1
dockerimg=docker-20-04

# Turning lines into comma separted list: https://stackoverflow.com/a/17759537
sshkeys=$(doctl compute ssh-key list --format ID --no-header true | xargs | sed -e 's/ /,/g')
ip=$(doctl compute droplet create $name --size $fourgdroplet --region $london --image $dockerimg --ssh-keys $sshkeys --wait --format PublicIPv4 --no-header)

echo -e "--> Copying project files to the Droplet"
sleep 10
scp -o 'StrictHostKeyChecking no' $filebeat_config $prom_config $env_config $docker_config root@$ip:/root/

echo -e "--> Copying Project Files to the Droplet"
sleep 10
scp -o 'StrictHostKeyChecking no' -r $docker_dir $src_dir root@$ip:/root/

echo -e "--> Starting the service"
sleep 10
ssh -o 'StrictHostKeyChecking no' root@$ip "docker-compose up --build --detach"

echo -e "The site can now be accesed @ http://$ip:8080"
echo -e "Grafana can now be accessed @ http://$ip:3000"
echo -e "Kibana can now be accessed @ http://$ip:5601"
echo -e "You can ssh into the Droplet using ssh root@$ip"
echo -e "You can destroy the droplet using ./src/scripts/digital_ocean/droplet-down.sh $name"