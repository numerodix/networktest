===========
networktest
===========

.. image:: https://api.travis-ci.org/numerodix/networktest.png?branch=master
    :target: https://travis-ci.org/numerodix/networktest

.. image:: https://ci.appveyor.com/api/projects/status/ojcojhf1wlc837ug/branch/master?svg=true
    :target: https://ci.appveyor.com/project/numerodix/networktest

This tool answers the question: *am I connected to a (local) network and do I
have internet connectivity?*




How does it work?
=================

While the question *"can this computer reach the Internet?"* is quite easily
answered by trying to load ``http://www.google.com`` in a browser, common
network configurations have a number of parameters that may vary between one
case of *"it's not working"* and the next.

``networktest`` aims to be a tool that makes it easy to establish the key
network configuration characteristics and get a quick idea of what might be
wrong. It is a smoke test of your network setup.

The questions we ask are:

1. **Am I connected to any networks?** If so, what kind of network are they?
   Loopback interface? Link-local? Or is it a network that allows me to route
   packets to the Internet (the kind we want)?

2. **What ip addresses do I have on the networks I'm connected too?** These
   are addresses that should be reachable (ie. pingable) from myself.

3. **What gateways do I have?** A gateway (or router) is a host on my network
   that allows me to route packets to other networks.

4. **Can I reach Internet hosts by ip?** If so, my gateway is working.

5. **What DNS servers (nameservers) do I have?** DNS servers allow me to
   resolve Internet hostnames to ip addresses.

6. **Can I reach Internet hosts by hostname?** If so, I have "Internet
   connectivity" in the common sense.




Supported platforms
===================

``networktest`` is written in Go and targets Linux, FreeBSD, OS X, and Windows.




Usage
=====

To build::
    
    $ make

To detect IPv4 networking::

    $ bin/havenet
     + Scanning for networks...
        <docker0>   172.17.0.0      / 255.255.0.0    
        <eth0>      192.168.1.0     / 255.255.255.0  
        <wlan0>     192.168.1.0     / 255.255.255.0  
     + Detecting ips...
        <lo>        127.0.0.1       / 255.0.0.0         ping: 0.072 ms
        <docker0>   172.17.0.1      / 255.255.0.0       ping: 0.067 ms
        <eth0>      192.168.1.6     / 255.255.255.0     ping: 0.065 ms
        <wlan0>     192.168.1.10    / 255.255.255.0     ping: 0.044 ms
     + Detecting gateways...
        <eth0>      192.168.1.1       ping: 0.706 ms
         ip:        192.168.1.6    
         ip:        192.168.1.10   
     + Testing internet connection...
        b.root-servers.net.  192.228.79.201   ping: 181.0 ms
     + Detecting dns servers...
        127.0.1.1         ping: 0.052 ms
     + Testing internet dns...
        facebook.com      ping: 315.0 ms
        gmail.com         ping: 37.70 ms
        google.com        ping: 41.10 ms
        twitter.com       ping: 130.0 ms
        yahoo.com         ping: 154.0 ms

To detect IPv6 networking::

    $ bin/havenet -6
     + Scanning for networks...
        <lo>                                            ::1 / 128  [scope: host]
        <eth0>                                       fe80:: /  64  [scope: link]
        <wlan0>                                      fe80:: /  64  [scope: link]
     + Detecting ips...
        <lo>                                            ::1 / 128  ping: 0.047 ms
        <eth0>                    fe80::16da:fae1:c9ea:a4b9 /  64  ping: N/A
        <wlan0>                   fe80::762f:fe64:b7c7:7b7a /  64  ping: N/A
     + Detecting gateways...
        none found
