#!/bin/sh

incus stop ssh-template
lxops delete -name ssh-template ssh.yaml
incus image delete alpine-ssh
