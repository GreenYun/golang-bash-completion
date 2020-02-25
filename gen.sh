#!/bin/bash

source=$(cat /etc/bash_completion)
INSTALL_DIR=${source#\.}
INSTALL_DIR=${INSTALL_DIR%/*}

[[ -z $INSTALL_DIR ]] && INSTALL_DIR=/usr/share/bash-completion

cat Makefile.in > Makefile

files=( "go" "gofmt" )
for i in ${files[@]}; do
    echo -e "\tinstall ./${i} ${INSTALL_DIR}/completions/${i}" >> Makefile
done