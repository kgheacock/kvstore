<h1>CSE Assignment</h1>
<h2>Assignment Overview:</h2>

<h3>kv store</h3>
<ul>
  <li>kvstore/dal contains the data structure and required methods CRUD methods</li>
<li>kvstore/kvs handles any interfacing with the DAL</li>
<li>kvstore/model contains definitions of structures used in kvstore/dal</li>
</ul>
<h3>router</h3>
<ul>
<li>router/middleware contains methods for proxying</li>
<li>router/router contains methods for assigning http routes to their handlers</li>
<li>router/model contains definitions of structures used in the router package</li>
  </ul>
<h3>hasher</h3>
<ul>
  <li>hasher/dal contains the entire consistent hashing library</li>
  <li>hasher/model contains definitions of structures used in consistent hashing library</li>
</ul>
<h3>vectorclock</h3>
<ul>
  <li>vectorclock/dal contains the vectorclock library</li>
  <li>vectorclock/model contains definitions of functions and structures used in vectorclock library</li>
</ul>
<h3>Dockerfile</h3>
Creation of the docker image has the following steps: 
  <ol>
  <li>specify the base image to be golang:alpine (an alpine linux golang deployment)</li>
  <li><optional> set working directory
    to app (note: this can be changed for ease of debugging)</li>
   <li> copy all files from ./ to the appropriate gopath on the new machine</li>
  <li> install git to aid with go commands</li>
  <li> install dependencies</li>
  <li> run a go install to install the program on the container</li>
  <li> set the path of the program to be run when the container is started</li>
  </ol>


