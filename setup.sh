#!/usr/bin/env bash

# find the pluto root directory
ROOT=${PLUTO:-$HOME/pluto}

# get the repository
go get -u github.com/Zac-Garby/pluto

# copy the libraries to $ROOT/libraries if
# they're not already there
if [ -e $ROOT/libraries ]
then
    echo "libraries already exist. skipping"
else
    echo "libraries don't exist. creating them"

    # copy the libraries over
    cp -R $GOPATH/src/github.com/Zac-Garby/pluto/libraries $ROOT
fi

echo $ROOT