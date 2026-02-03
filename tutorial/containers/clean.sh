#!/bin/sh

for c in alpine-ssh ssh1 ssh2; do
    incus stop $c
    lxops destroy -name $c ssh.yaml
done
