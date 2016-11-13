<img src="http://golang.org/doc/gopher/frontpage.png" alt="Golang logo" align="right"/>

[![GoDoc](https://godoc.org/github.com/SparrowDb/sparrowdb?status.svg)](https://godoc.org/github.com/SparrowDb/sparrowdb)
[![Build Status](https://travis-ci.org/SparrowDb/sparrowdb.svg?branch=master)](https://travis-ci.org/SparrowDb/sparrowdb)
[![Go Report Card](https://goreportcard.com/badge/github.com/SparrowDb/sparrowdb)](https://goreportcard.com/report/github.com/SparrowDb/sparrowdb)

Whats is SparrowDB?
====================
SparrowDB is an image database that works like an append-only object store. Sparrow has tools that allow image processing and HTTP server to access images.


Sparrow Object Store
====================
Sparrow consists of three files – the actual Sparrow store file containing the images data, plus an index file and a bloom filter file.

There is a corresponding data definition record followed by the image bytes for each image in the storage file. The index file provides the offset of the data definition in the storage file.


Features
====================
1. Built-in HTTP API so you don't have to write any server side code to get up and running.
2. Optimizations for image storing.
3. Websocket server to provide real time information about the server.


Getting started
====================
This short guide will walk you through getting a basic server up and running, and demonstrate some simple reads and writes.



Using Sparrow
====================
Creating a database:
	
	curl -X PUT http://127.0.0.1:8081/api/database_name

Show databases:

    curl -X GET http://127.0.0.1:8081/api/_all


Sending an image to database:

	curl -i -X PUT -H "Content-Type: multipart/form-data"  \
        -F "uploadfile=@image.jpg" \
        http://127.0.0.1:8081/api/database_name/image_key


Querying an image:

	curl -X GET http://127.0.0.1:8081/api/database_name/image_key


Accessing image from browser:
	
	http://localhost:8081/g/database_name/image_key


Token
====================

If is set in database configuration file, generate_token = true, SparrowDB will generate a token for each image uploaded. The token’s value is randomly assigned by and stored in database. The token effectively eliminates attacks aimed at guessing valid URLs for photos.

Accessing image from browser with token:
	
	http://localhost:8081/g/database_name/image_key/token_value


Image Processing
====================

SparrowDB uses [bild](https://github.com/anthonynsimon/bild) to allow image processing using [LUA](https://github.com/yuin/gopher-lua) script.

All SparrowDB scripts must be in 'scripts' folder.

Example of script that converts image to grayscale:

```lua
-- loads sparrowdb module
local s = require("sparrowdb")

-- get image data
image = s.getInputImage()

-- convert image to grayscale
out = s.grayscale(image)

-- set new output for image
s.setOutputImage(out)
```


Replication
====================

SparrowDB integrates with [NATS](https://github.com/nats-io/gnatsd) to replicate all data using pub/sub.

To use Sparrow in cluster mode, change the following options in config/sparrow.xml:

	<enable_cluster>false</enable_cluster>


And set NATS servers:

	<publisher_servers>nats://localhost:4222</publisher_servers>


Or multiple servers:

	<publisher_servers>nats://localhost:4222, nats://localhost:4223, nats://localhost:4224</publisher_servers>



License
====================
This software is under MIT license.
