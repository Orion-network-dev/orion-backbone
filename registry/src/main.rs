mod reg_config;
mod wg;
mod ws;

use std::{
    fs,
    slice::Iter,
    sync::{mpsc::Sender, Arc, Mutex},
};

use crate::wg::make_wireguard_interface;
use ::config::Config;
use futures_util::{stream::StreamExt, FutureExt, SinkExt};
use protocol::*;
use tokio::sync::mpsc;
use warp::{filters::ws::Message, Filter};
use webpki::{TlsClientTrustAnchors, TrustAnchor};
use ws::do_hello;

const SERVER_TAGLINE: &str = "Welcome to Orion! Have fun!";

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    pretty_env_logger::init();

    // load the configuration
    let settings = Config::builder()
        .add_source(config::File::with_name("./registry.toml"))
        .add_source(config::Environment::with_prefix("ORIONREG"))
        .build()?
        .try_deserialize::<reg_config::RegistryConfiguration>()?;

    let wg = make_wireguard_interface(&settings)?;
    let server_information = Arc::new(OrionRegProtocolV1ServerHelloRequest {
        tagline: SERVER_TAGLINE.to_string(),
        public_address: OrionRegProtocolV1WireguardPeerInfo {
            address: settings.public_address,
            port: settings.wg_port,
        },
        nonce: "".to_string(),
    });
    let clients = Arc::new(Mutex::new(vec![]));

    let routes = warp::path("ws")
        .and(warp::ws())
        .map(move |ws: warp::ws::Ws| {
            let server_information = server_information.clone();
            let clients = clients.clone();
            ws.on_upgrade(|mut websocket| {
                tokio::spawn(async move {
                    let rootCAder =
                        fs::read("/home/matthieu/Documents/orionv3/dev/mtls/ca/rootCA.der")
                            .unwrap();
                    let tal = &[TrustAnchor::try_from_cert_der(&rootCAder).expect("invalid CA")];
                    let trust_anchor = Arc::new(TlsClientTrustAnchors(tal));
                    let information =
                        do_hello(&server_information, &mut websocket, trust_anchor.clone())
                            .await
                            .unwrap();
                    let (z, mut r) = mpsc::channel::<Message>(1);
                    clients.lock().unwrap().push((information.clone(), z));

                    {
                        let clients_l = clients.lock().unwrap();
                        let new_member = serde_json::to_string(&OrionRegProtocolV1JoinEvent {
                            name: information.name.clone(),
                            metadata: information.metadata.clone(),
                        })
                        .unwrap();
                        for (info, send) in clients_l.iter() {
                            let chan = send.clone();
                            let new_member = new_member.clone();
                            tokio::spawn(async move {
                                chan.send(Message::text(new_member.clone())).await.unwrap();
                            });
                        }
                    }

                    let (mut tx, mut rx) = websocket.split();

                    loop {
                        tokio::select! {
                            Some(message) = r.recv() => {
                                tx.send(message).await.unwrap();
                            },
                            Some(Ok(message)) = rx.next() => {

                            },
                            else => break,
                        }
                    }
                })
                .map(|z| z.unwrap())
            })
        });

    warp::serve(routes)
        .tls()
        .cert_path("dev/sslhost/cert.crt")
        .key_path("dev/sslhost/cert.key")
        .run(([0, 0, 0, 0], 3030))
        .await;
    Ok(())
}
