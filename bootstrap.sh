#!/bin/bash

echo -e "\n--> Bootstrapping Minitwit\n"

echo -e "\n--> Loading environment variables from secrets file\n"
source secrets

echo -e "\n--> Checking that environment variables are set\n"
# check that all variables are set
[ -z "$TF_VAR_do_token" ] && echo "TF_VAR_do_token is not set" && exit
[ -z "$SPACE_NAME" ] && echo "SPACE_NAME is not set" && exit
[ -z "$STATE_FILE" ] && echo "STATE_FILE is not set" && exit
[ -z "$AWS_ACCESS_KEY_ID" ] && echo "AWS_ACCESS_KEY_ID is not set" && exit
[ -z "$AWS_SECRET_ACCESS_KEY" ] && echo "AWS_SECRET_ACCESS_KEY is not set" && exit

echo -e "\n--> Checking that all the necessary files exist\n"
# Check that all files exist
[ ! -f "ssh_key/terraform" ] && echo "ssh_key/terraform does not exist. Please generate a ssh_key" && exit
[ ! -f "docker-stack.yml" ] && echo "docker-stack.yml file does not exist" && exit
[ ! -f "filebeat.yml" ] && echo "filebeat.yml does not exist" && exit
[ ! -f "prometheus.yml" ] && echo "prometheus.yml" && exit

echo -e "\n--> Initializing terraform\n"
# initialize terraform
terraform init \
    -backend-config "bucket=$SPACE_NAME" \
    -backend-config "key=$STATE_FILE" \
    -backend-config "access_key=$AWS_ACCESS_KEY_ID" \
    -backend-config "secret_key=$AWS_SECRET_ACCESS_KEY"

# check that everything looks good
echo -e "\n--> Validating terraform configuration\n"
terraform validate

# create infrastructure
echo -e "\n--> Creating Infrastructure\n"
terraform apply -auto-approve -parallelism=1

# ensure all nodes have the necessary config files
# sleep to reduce the number of failed connections
sleep 5

echo -e "\n--> Copying the config files to the necessary nodes\n"

# sleep to reduce the number of failed connections
sleep 5

echo -e "\n--> Copying the config files to all the nodes"
bash src/scripts/terraform/scp_config_to_all_nodes.sh

# sleep to reduce the number of failed connections
sleep 5

# deploy the stack to the cluster
echo -e "\n--> Deploying the Minitwit stack to the cluster\n"
ssh \
    -o 'StrictHostKeyChecking no' \
    root@$(terraform output -raw minitwit-swarm-leader-ip-address) \
    -i ssh_key/terraform \
    'docker stack deploy minitwit -c docker-stack.yml'

echo -e "\n--> Done bootstrapping Minitwit"
echo -e "--> The system needs to initialize, this can take up to a couple of minutes..."
echo -e "--> Site will be avilable @ http://$(terraform output -raw public_ip):8080"
echo -e "--> You can check the status of swarm cluster @ http://$(terraform output -raw minitwit-swarm-leader-ip-address):8888"
echo -e "--> You can view logs at @ http://$(terraform output -raw minitwit-swarm-leader-ip-address):5601"
echo -e "--> You can monitor the system @ http://$(terraform output -raw minitwit-swarm-leader-ip-address):3000"
echo -e "--> ssh to swarm leader with 'ssh root@\$(terraform output -raw minitwit-swarm-leader-ip-address) -i ssh_key/terraform'"
echo -e "--> To remove the infrastructure run: make destroy-prod"
