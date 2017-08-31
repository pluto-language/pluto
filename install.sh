#!/usr/bin/env bash

# find the pluto root directory
ROOT=${PLUTO:-$HOME/pluto}

# get the repository
echo -n "installing pluto... "
go get -u github.com/pluto-language/pluto
echo "DONE"

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

# install plp
echo -n "installing plp... "
go get -u github.com/pluto-language/plp
echo "DONE"

echo "pluto and plp are installed!"
