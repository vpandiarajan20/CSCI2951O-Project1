#!/bin/bash

########################################
############# CSCI 2951-O ##############
########################################
E_BADARGS=65
if [ $# -ne 1 ]
then
    echo "Usage: $0 <input>"
    exit $E_BADARGS
fi

input="$1"

./project1.out "$input"
