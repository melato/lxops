
# launch a container called "ssh-template" and apply the cloud-config files
lxops launch -name ssh-template ssh.yaml

# if you script this, wait a few seconds to give some time
# to the container installation scripts to complete
sleep 5

# create a snapshot from this container
incus stop ssh-template
incus snapshot create ssh-template copy
incus publish ssh-template/copy --alias alpine-ssh

# list the new image
incus image list alpine-ssh

# delete the container.  We no longer need it.
lxops delete -name ssh-template ssh.yaml
