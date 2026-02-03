#!/bin/sh

types() {
    for t in $*; do
        lxops-template exec -t doc/type.tpl -D type=$t -o md/$t.md
    done
}

types Config Filesystem Device Pattern HostPath
lxops-template build -c doc/build.yaml -i doc/pages/readme -o .
lxops-template build -c doc/build.yaml -i doc/pages -o md


