use anyhow::Context;
use defguard_wireguard_rs::host::Peer;
use defguard_wireguard_rs::key::Key;
use defguard_wireguard_rs::{InterfaceConfiguration, WGApi, WireguardInterfaceApi};
use log::{debug, info};

use crate::reg_config;

fn gen_wireguard_key() -> Key {
    let private_key = wireguard_keys::Privkey::generate();
    return Key::new(private_key.map(|z| z));
}

pub(crate) fn make_wireguard_interface(
    config: &reg_config::RegistryConfiguration,
) -> anyhow::Result<WGApi> {
    // Creates the wg-api instance in order to control this wireguard interface.
    let wgapi = WGApi::new(config.wg_interface.clone(), false)
        .context("initialization of the wg-api instance")?;
    // Create the wireguard interface if it does not exist.
    wgapi.create_interface()?;

    // reading the current interface information in order to read the current private key.
    let information = wgapi
        .read_interface_data()
        .context("reading interface data")?;
    // use the current private key or create a new private key.
    let private_key = information.private_key.or_else(|| Some(gen_wireguard_key())).unwrap();
    debug!("using private key: {}", private_key);
    // infer the wireguard configuration
    let interface = InterfaceConfiguration {
        name: config.wg_interface.clone(),
        prvkey: format!("{private_key}"),
        address: "0.0.0.1".to_string(),
        port: config.wg_port as u32,
        peers: information.peers.into_values().collect::<Vec<Peer>>(),
    };

    wgapi.configure_interface(&interface)?;

    info!("Reconfigured wireguard interface.");

    Ok(wgapi)
}
