# Orion registry protocol v1

## Initialization handshake

* Sends a hello request
* Server does a broadcast of the join event to all the clients.
* Once the client is disconnected, a left event is emitted.

## Events handling

### Peer join events

Once a clients receives a join events, this means a new client
joined the network, if this is a new client, it creates a new private key,
computes the public key and generates a pre-shared key.

* Creates a private/public key pair.
* Creates a per-shared key.
* Does udp-hole-punching.
* Sends a peer-request to the new peer.

### Peer left events

* Deletes the corresponding peer event if it exists.

### Peer request

* Creates a private/public key pair.
* Sends the peer-response to the client.