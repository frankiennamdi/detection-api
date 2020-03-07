# detection-api

## Introduction

The detection api provide a service that identifies suspicious travel using the speed of travel to and from 
consecutive geo location from login events from IP information. When the service receives an events it saves
it and the event history for that user, and evaluates the event prior and following the current event, if they
exist. It then calculates the speed it will require to travel between the two location and flag suspicious 
activity if it beyond a configured threshold. The default speed is 500, but this is overridden in the makefile
to 100 to get more violations. You can reset it to 500 when you start the server

```
 make run SUSPICIOUS_SPEED=500
```

OR

```
 make run-image SUSPICIOUS_SPEED=500
```

## Requirements
1. go 1.13 
2. make, my version is GNU Make 4.2.1
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

To generate a specific number of randomized events for the the samples users. Please be aware that 
I am only using the three IPs provided in the sample and hence the randomness is limited. 

Please see Makefile for more information. 

## External Dependencies

1. Used open source DB drivers as indicated in the mod files, also used for DB migration
2. Used mux REST server
3. Notable mention, for configuration handling I used a library I started a few month ago. I wanted
a library that inject environment variable into yaml files when available and binds them to go struct. I
did not see any good one out there, so I wrote one. It is still quite raw, but it go the job done here 
for me. https://github.com/frankiennamdi/go-configuration

## Constraints

1. event UUID is the primary key.
2. username and timestamp are unique keys.
3. duplicates based on the constraints are ignored and only the original request is used in evaluation.

## Possible Future Improvements

* Possible integration test that involves a running server.
* Possible use of a more traditional database. 


