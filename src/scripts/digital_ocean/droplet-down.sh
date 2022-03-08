#!/bin/bash
name=$1

doctl compute droplet delete $name --force