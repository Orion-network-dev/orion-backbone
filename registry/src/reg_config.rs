use std::net::IpAddr;

use serde::{Deserialize, Serialize};

#[derive(Deserialize, Serialize)]
pub struct RegistryConfiguration {
    pub public_address: IpAddr,
    pub wg_port: u16,
    pub wg_interface: String,
}