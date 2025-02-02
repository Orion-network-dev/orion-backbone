# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  config.vm.box = "debian/bookworm64"

  config.vm.provision "shell", inline: <<-SHELL
    apt-get update && apt-get install -y curl
    curl -s https://deb.frrouting.org/frr/keys.gpg | tee /usr/share/keyrings/frrouting.gpg > /dev/null
    echo deb '[signed-by=/usr/share/keyrings/frrouting.gpg]' https://deb.frrouting.org/frr $(lsb_release -s -c) frr-stable | tee -a /etc/apt/sources.list.d/frr.list
    apt-get update && apt-get install --allow-downgrades -y /vagrant/dist/orion-backbone_*_amd64.deb
  SHELL

  config.vm.define "registry" do |registry|
    registry.vm.network "private_network", ip: "192.168.50.200"
    registry.vm.hostname = "registry"

    registry.vm.provision "shell", inline: <<-SHELL
      cp /vagrant/secret/registry/registry.pem /etc/oriond/identity.pem
      systemctl enable --now orion-registry
    SHELL
    registry.vm.network "forwarded_port", guest: 64431, host: 64431
  end

  # Provision two orion node VMs 
  (0..1).each do |n| 
    config.vm.define "node#{n}" do |node|
      node.vm.network "private_network", ip: "192.168.50.1#{n}"
      node.vm.hostname = "node#{n}"
  
      node.vm.provision "shell", inline: %{
        sed -i "s/bgpd=no/bgpd=yes/g" /etc/frr/daemons
        sed -i "s/pimd=no/pimd=yes/g" /etc/frr/daemons
        sed -i "s/pim6d=no/pim6d=yes/g" /etc/frr/daemons
        systemctl restart frr
        cp /vagrant/secret/node#{n}/identity.pem /etc/oriond/identity.pem
        echo "192.168.50.200 reg.orionet.re" \> /etc/hosts
        systemctl enable --now oriond
      }
    end
  end
end
