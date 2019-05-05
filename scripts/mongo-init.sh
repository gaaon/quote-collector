#!/bin/bash

docker run -dit -p 27017:27017 -v $(PWD)/db:/data/db --name my-mongo -d mongo