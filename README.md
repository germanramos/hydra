Hydra
=====

"If a head is cut off, two more will take its place"

http://innotech.github.io/hydra/

# What is Hydra?
Hydra is multi-cloud application discovery, management and balancing service.
It attempts to ease the routing and balancing burden from servers and delegate it on the client (browser, mobile app, etc).

Hydra is composed of 2 server-side applications:
* <a href="https://github.com/innotech/hydra">Hydra Server (This repository)</a>
* <a href="https://github.com/innotech/hydra-basic-probe">Hydra Basic Probe</a>

Also there is several client-side library:
* <a href="https://github.com/innotech/hydra-javascript-client">Hydra Javascript Client</a>
* <a href="https://github.com/innotech/hydra_node_client">Hydra NodeJS Client</a>
* <a href="https://github.com/innotech/hydra-java-client">Hydra Java Client</a>

Finally, there is web monitor that connects to an Hydra Server an print all the information:
* <a href="https://github.com/innotech/hydra-web-monitor">Hydra Sys Monitor</a>

To use Hydra, you need to deploy: 
* At least one Hydra server (per cloud prefered) to comunicate your applications to the world.
* An Hydra Basic Probe on each single server your application is deployed in. 

For information on how to deploy each individual piece, check out their respective repositories.

# Hydra Server

The Hydra server is also composed by one core application and one or more add-ons. These add-ons are "workers" and Hydra use them in order to delegate the balance calculation. You can deploy as many workers as you want in the same server than Hydra or in different servers. The communication between Hydra Core and the workers is made with a TCP connection using ZeroMQ.

## Hydra Core
* <a href="https://github.com/innotech/hydra/blob/master/Documentation/configuration.md">Instalation and Configuration</a>
* <a href="https://github.com/innotech/hydra/blob/master/Documentation/development_enviroment.md">Development Environment</a>
* <a href="https://github.com/innotech/hydra/blob/master/Documentation/roadmap.md">Roadmap</a>

## Hydra Workers
* <a href="https://github.com/innotech/hydra-worker-round-robin">Round Robin Worker</a>
* <a href="https://github.com/innotech/hydra-worker-map-sort">Map and Sort Worker</a>
* <a href="https://github.com/innotech/hydra-worker-sort-by-number">Sort by Number Workder</a>
* <a href="https://github.com/innotech/hydra-worker-map-by-limit">Map by Limir Worker</a>

# License

(The MIT License)

Authors:  
Germán Ramos &lt;german.ramos@gmail.com&gt;  
Pascual de Juan &lt;pascual.dejuan@gmail.com&gt;  
Jonas da Cruz &lt;unlogic@gmail.com&gt;  
Luis Mesas &lt;luismesas@gmail.com&gt;  
Alejandro Penedo &lt;icedfiend@gmail.com&gt;  
Jose María San José &lt;josem.sanjose@gmail.com&gt;  

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
'Software'), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
