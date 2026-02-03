
## lxops extract copies files from the image to a
## "alpine-ssh" template container
## It copies the files that are needed to create
## instance-specific non-root disk devices
lxops extract -name alpine-ssh ssh.yaml

lxops launch -name ssh1 ssh.yaml
lxops launch -name ssh2 ssh.yaml
