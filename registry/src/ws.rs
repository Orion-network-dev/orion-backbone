use std::sync::Arc;

use anyhow::{anyhow, bail, Result};
use base64::prelude::*;
use futures_util::{SinkExt, StreamExt};
use log::info;
use protocol::{OrionRegProtocolV1ServerHelloRequest, OrionRegProtocolV1ServerHelloResponse};
use rand::RngCore;
use warp::filters::ws::{Message, WebSocket};
use webpki::{DnsNameRef, Time, TlsClientTrustAnchors};

static ALL_SIGALGS: &[&webpki::SignatureAlgorithm] = &[
    &webpki::RSA_PKCS1_2048_8192_SHA256,
    &webpki::RSA_PKCS1_2048_8192_SHA384,
    &webpki::RSA_PKCS1_2048_8192_SHA512,
    &webpki::RSA_PKCS1_3072_8192_SHA384,
];

pub(crate) async fn do_hello<'a>(
    server_information: &OrionRegProtocolV1ServerHelloRequest,
    ws: &mut WebSocket,
    trust_anchors: Arc<TlsClientTrustAnchors<'a>>,
) -> Result<OrionRegProtocolV1ServerHelloResponse> {
    let mut nonce = [0u8; 4096];
    rand::thread_rng().fill_bytes(&mut nonce);
    let mut server_information_nonce =
        OrionRegProtocolV1ServerHelloRequest::clone(&server_information);
    server_information_nonce.nonce = BASE64_STANDARD.encode(&nonce).to_string();

    let request = serde_json::to_string(&server_information_nonce)?;
    ws.send(Message::text(request)).await?;

    let response = ws.next().await;

    match response {
        Some(Ok(message)) => {
            if message.is_text() {
                let data: OrionRegProtocolV1ServerHelloResponse = serde_json::from_str(
                    &message
                        .to_str()
                        .map_err(|_| anyhow!("unable to decode the string from the handshake"))?,
                )?;

                let certificate_der = BASE64_STANDARD.decode(data.certificate.clone())?;
                let end_entity_cert =
                    webpki::EndEntityCert::try_from(certificate_der.as_slice()).unwrap();

                let is_dns_valid = end_entity_cert
                    .verify_is_valid_for_dns_name(
                        DnsNameRef::try_from_ascii_str(&data.name).unwrap(),
                    )
                    .is_ok();
                let is_trusted = end_entity_cert
                    .verify_is_valid_tls_client_cert_ext(
                        &ALL_SIGALGS,
                        &trust_anchors,
                        &[],
                        Time::try_from(std::time::SystemTime::now())?,
                    )
                    .is_ok();

                if is_trusted && is_dns_valid {
                    let signature = BASE64_STANDARD.decode(data.signature.clone())?;
                    end_entity_cert.verify_signature(
                        &webpki::RSA_PKCS1_2048_8192_SHA512,
                        &nonce,
                        &signature,
                    )?;
                    return Ok(data);
                } else {
                    info!("{is_trusted} {is_dns_valid}");

                    //

                    bail!("certificate is not valid");
                }
            } else {
                bail!("handshake is no in cleartext")
            }
        }
        _ => {
            bail!("unable to complete the hello handshake")
        }
    }
}
