#!/bin/bash

echo -e "\n--> Checking that all the necessary files exist\n"
# Check that all files exist
[ ! -f "secrets" ] && echo "secrets does not exist. Please generate a ssh_key" && exit

echo -e "\n--> Loading environment variables from secrets file\n"
source secrets

echo -e "\n--> Checking that environment variables are set\n"
# check that all variables are set
[ -z "$TF_VAR_do_token" ] && echo "TF_VAR_do_token is not set" && exit

terraform destroy -auto-approve
rm -rf .terraform