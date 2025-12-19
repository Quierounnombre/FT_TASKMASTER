#!/usr/bin/env zsh

i=$LIMIT

while (($i > 0));
do
    export CURRENT_TIME=$(date +"%T")
    sleep 1
	i=$((i - 1))
done
