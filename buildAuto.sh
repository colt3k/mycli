#!/bin/bash

go install .
PROG=mycli source ./bash_autocomplete
cp ./bash_autocomplete /usr/local/etc/bash_completion.d/mycli
source /usr/local/etc/bash_completion.d/mycli
