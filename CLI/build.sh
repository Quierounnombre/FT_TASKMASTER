#!/bin/bash

rm /usr/local/bin/taskcli
rm /app/CLI/taskcli

go build -buildvcs=false -o taskcli /app/CLI


cp /app/CLI/taskcli /usr/local/bin/
chmod +x /usr/local/bin/taskcli

rm /app/CLI/taskcli
