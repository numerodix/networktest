=======
havenet
=======

.. image:: https://api.travis-ci.org/numerodix/gohavenet.png?branch=master
    :target: https://travis-ci.org/numerodix/gohavenet

This tool answers the question: *am I connected to a (local) network and do I
have internet connectivity?*



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
