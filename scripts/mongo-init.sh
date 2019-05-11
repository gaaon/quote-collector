#!/bin/bash

docker run -dit -p 27017:27017 -v $(pwd)/db:/data/db --name my-mongo -d mongo