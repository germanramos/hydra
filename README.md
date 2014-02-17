Hydra
=====

"If a head is cut off, two more will take its place"

http://innotech.github.io/hydra/

# What is Hydra?
Hydra is multi-cloud application discovery, management and balancing service.
It attempts to ease the routing and balancing burden from servers and delegate it on the client (browser, mobile app, etc).

Hydra is composed of 3 server-side applications:
* <a href="https://github.com/innotech/hydra_server">Hydra Server</a>
* <a href="https://github.com/innotech/hydra_app_manager">AppManager</a>
* <a href="https://github.com/innotech/hydra_basic_probe">AppManager Info Server (Probe)</a>

Also there is several client-side library:
* <a href="https://github.com/innotech/hydra-javascript-client">Hydra Javascript Client</a>
* <a href="https://github.com/innotech/hydra_node_client">Hydra NodeJS Client</a>
* <a href="https://github.com/innotech/hydra-java-client">Hydra Java Client</a>

Finally, there is web monitor that connects to an Hydra Server an print all the information:
* <a href="https://github.com/innotech/hydra/tree/master/src/app_manager_sysmon">Hydra Sys Monitor</a>

To use Hydra, you need to deploy all the 3 server-side applications. 
* An Hydra server per cloud to comunicate your applications to the world.
* An AppManager per app to notify information to Hydra.
* An AppManager Info Server on each single server your application is deployed in. AppManager will gather the status information connecting to it.

For information on how to deploy each individual piece, check out their respective repositories.

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
