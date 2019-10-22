# CSE Assignment 2
Assignment Overview:

-----------------kv store---------------------------------------
-kvstore/dal contains the data structure and required methods CRUD methods
-kvstore/kvs handles any interfacing with the DAL
-kvstore/model contains definitions of structures used in kvstore/dal
----------------router------------------------------------------
-router/middleware contains methods for proxying
-router/router contains methods for assigning http routes to their handlers
-router/model contains definitions of structures used in the router package


---------------Dockerfile----------------------------------------
Creation of the docker image has the following steps: 1) specify the base image
to be golang:alpine (an alpine linux golang deployment) 2) <optional> set working directory 
to app (note: this can be changed for ease of debugging) 3) copy all files from ./ to 


