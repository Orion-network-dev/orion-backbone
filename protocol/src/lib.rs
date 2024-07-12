use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::net::IpAddr;

#[derive(Serialize, Deserialize, Clone)]
pub struct OrionRegProtocolV1WireguardPeerInfo {
    pub address: IpAddr,
    pub port: u16,
}

#[derive(Serialize, Deserialize, Clone)]
pub struct OrionRegProtocolV1ServerHelloRequest {
    pub tagline: String,
    // this is used to discover the various endpoints.
    pub public_address: OrionRegProtocolV1WireguardPeerInfo,
    // send a nonce
    pub nonce: String,
}

#[derive(Serialize, Deserialize, Clone)]
pub struct OrionRegProtocolV1ServerHelloResponse {
    pub name: String,
    pub metadata: HashMap<String, String>,

    pub certificate: String,
    pub signature: String,
}

#[derive(Serialize, Deserialize, Clone)]
pub struct OrionRegProtocolV1JoinEvent {
    pub name: String,
    pub metadata: HashMap<String, String>,
}

#[derive(Serialize, Deserialize)]
pub struct OrionRegProtocolV1LeftEvent {}

#[derive(Serialize, Deserialize)]
pub struct OrionRegProtocolV1PeerRequest {}

#[derive(Serialize, Deserialize)]
pub struct OrionRegProtocolV1PeerResponse {}

#[derive(Serialize, Deserialize)]
pub struct OrionRegProtocolV1PeerAcknowledge {}

#[derive(Serialize, Deserialize)]
pub struct OrionRegProtocolV1StunRequest {}

#[derive(Serialize, Deserialize)]
pub struct OrionRegProtocolV1StunResponse {}

///
#[derive(Serialize, Deserialize)]
#[serde(tag = "t")]
pub enum OrionRegProtocolV1Message {
    /// First message when a client joins which show the server's information.
    HelloRequest(OrionRegProtocolV1ServerHelloRequest),
    /// Response of a `HelloRequest` which contains the peer's information and metadata.
    HelloResponse(OrionRegProtocolV1ServerHelloResponse),
    /// Event when a client joins which contains the details and metadata of a peer.
    JoinEvent(OrionRegProtocolV1JoinEvent),
    /// Event when a client lefts.
    LeftEvent(OrionRegProtocolV1LeftEvent),
    /// Event when a client wants to peer with someone.
    PeerRequest(OrionRegProtocolV1PeerRequest),
    /// Response of a `PeerRequest` which contains the new tunnel credentials.
    PeerResponse(OrionRegProtocolV1PeerResponse),
    /// Final message of a Peering system.
    /// Precedes a `PeerResponse` and means the `PeerResponse` credentials
    /// is ready to be used.
    PeerAcknowledge(OrionRegProtocolV1PeerAcknowledge),

    /// Asks the server for a stun request.
    StunPeerRequest(OrionRegProtocolV1StunRequest),
    /// Response of a `StunPeerRequest` once the peer is connected.
    StunPeerResponse(OrionRegProtocolV1StunResponse),
}
