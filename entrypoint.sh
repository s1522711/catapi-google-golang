#!/bin/bash
cd /home/container

# Output Current Go Version
echo "Current Go Version: $(go version)"

# Replace Startup Variables
MODIFIED_STARTUP=`eval echo $(echo ${STARTUP} | sed -e 's/{{/${/g' -e 's/}}/}/g')`
echo ":/home/container$ ${MODIFIED_STARTUP}"

# Run the Startup Command
${MODIFIED_STARTUP}