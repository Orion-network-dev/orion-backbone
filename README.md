# Orion network backbone

This is the networking daemons and networking backbone for `Orion`, all the 
applicative and bgp management are available on another repo: https://github.com/MatthieuCoder/Orion

This is a simple program called `oriond` that manages wireguard tunnels to other peers of the Orion network.

There is also a server program called `registry` that tracks all the Orion networks connections and all signaling between peers to archieve UDP-hole-punching and negociation between clients.
