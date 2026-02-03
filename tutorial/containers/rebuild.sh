#!/bin/sh

incus stop ssh1
lxops delete -name ssh1 ssh.yaml
lxops launch -name ssh1 ssh.yaml
