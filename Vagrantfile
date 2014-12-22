# -*- mode: ruby -*-
# vi: set ft=ruby :

# Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.vm.define "boot2docker"

  config.vm.box = "yungsang/boot2docker"

  config.vm.provision :shell do |sh|
    sh.inline = <<-EOT
      sudo echo 'EXTRA_ARGS="--label=vm=vitrualbox --label=box=yungsang/boot2docker"' > /var/lib/boot2docker/profile
      sudo /etc/init.d/docker restart
    EOT
  end
end
