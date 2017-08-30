#!/usr/bin/env bash

# find the pluto root directory
ROOT=${PLUTO:-$HOME/pluto}

# get the repository
go get -u github.com/Zac-Garby/pluto

# copy the libraries to $ROOT/libraries if
# they're not already there
if [ -e $ROOT/libraries ]
then
    echo "libraries already exist"
else
    echo "libraries don't exist. creating them"
    
    # make the libraries directory
    mkdir -p $ROOT/libraries
fi

# copy the libraries over
cp -R $GOPATH/src/github.com/Zac-Garby/pluto/libraries $ROOT

echo "Pluto is installed! Type 'pluto' to run the REPL."
