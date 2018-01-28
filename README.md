### Word Processor

This repository contains the simple implementation of two servers
of which are:

* The input server (located in the `input` directory) listens for
 POST http requests at `/` on the port specified in the command line flags
 (defaults to `5555`) which reads the raw body and appends it to 
 the file location also supplied in flags (defaults to `/tmp/data.txt`)
* The second server (located in the `stats` directory) reads the same file location (again this is supplied in flags)
 and processes the words in the file and determines some statistics which are
 returned at the endpoint `/stats`. The file pointer and statistics are refreshed
 every 10 seconds.
 
 ##### Running
 I have included a Makefile which has two tasks
 - `make start` will build the go binaries and then use docker-compose to start
   the two services.
 - `make stop` will stop the running binaries (also deleting the attached volume) and remove
  the built binaries from the system.

You can also run by compiling the binaries yourself and running them on your local machine.