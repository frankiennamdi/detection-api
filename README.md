# detection-api

## Introduction

The detection api provides a service that identifies suspicious travel using the speed of travel to and from 
consecutive geo location derived from the IP information of login events. When the service receives an event it saves
it as part of the event history for that user, and evaluates the event prior and following the current event, if they
exist. It then calculates the speed it will takes to travel between the two location and flag suspicious 
activity if it beyond a configured threshold. The default speed is 500, but this is overridden in the makefile
to 100 to get more violations during testing. You can reset it to 500 when you start the server. See make file for
options.

```
 make run SUSPICIOUS_SPEED=500
```

OR

```
 make run-image SUSPICIOUS_SPEED=500
```

## Requirements
1. Go 1.13 
2. Make, my version is GNU Make 4.2.1
3. Docker
4. Tested on Mac OS

## Build and Executions

The application contain a set of make build targets. 

1. test, build and run locally 

```
 make run
```

2. test, build image and run image in docker 

```
 make run-image
```

3. Generate sample events

``` 
make run-generator
```
or 

```
make run-generator NUM_OF_EVENTS=100
```

To generate a specific number of randomized events for the samples users. Please be aware that 
I am only using the three IPs provided in the sample and hence the randomness is limited. 

Please see Makefile for more information. 

## Docker volume mapping with caveat
The location for the database is in the resource/event-db folder. And the name is configurable. When running in docker 
you can map the volume to the local storage e.g. `-v $(PWD)/resources/event-db:/app/resources/event-db` in the 
`make run-image` target. The caveat is that this make the application slower from my observation. 

## External Dependencies

1. Used open source DB drivers as indicated in the mod files, also used for DB migration.
2. Used mux REST framework.
3. Notable mention, for configuration handling I used a library I started a few months ago. I wanted
a library that inject environment variables into yaml files when available and binds them to go struct. I
did not see any good one out there, so I wrote one. It is still quite raw, but it got the job done here 
for me. https://github.com/frankiennamdi/go-configuration

## Constraints/ Design Decisions

1. **event UUID** is the primary key.
2. **username** and **timestamp** are unique keys.
3. duplicates based on the aforementioned constraints are ignored and only the original request is used in evaluation.
4. Each events results in an open action on the database, this is in order to not keep the db open for too long
and allow the database to manage the access. The drawback of this means that it is possible to reach the max 
number of file descriptors allowed on the system. To make the application graceful, we introduced the **DB_MAX_CONN**
and **IP_GEO_DB_MAX_CONN** to allow you tune this based on the limit of the host machine. Currently they are
both set to **200**
5. Used db file as I felt this was more usable offline. 

## Possible Future Improvements

* Possible use of a more tradition database or use of the memory version of SQLite. 
* Possible integration tests that involves a running server.
* Possible use of a more traditional database to handle the transaction load.
* Clean up the use of pointers where needed. I tried to balance the need for nil values and check, immutability and copying. 
In cases where I had the difficult choices I tried to hide the struct properties and not allow modification 
after creation, except through the constructors in some cases. This means that the most properties are unexported
and can only be read by functions that return their value after construction. 


