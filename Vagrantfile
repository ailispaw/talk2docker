# A dummy plugin for DockerRoot to set hostname and network correctly at the very first `vagrant up`
module VagrantPlugins
  module GuestLinux
    class Plugin < Vagrant.plugin("2")
      guest_capability("linux", "change_host_name") { Cap::ChangeHostName }
      guest_capability("linux", "configure_networks") { Cap::ConfigureNetworks }
    end
  end
end

Vagrant.configure(2) do |config|
  config.vm.define "docker-root"

  config.vm.box = "ailispaw/docker-root"
  config.vm.box_version = ">= 1.3.9"

  config.vm.network :forwarded_port, guest: 2375, host: 2375, auto_correct: true, disabled: true

  config.vm.synced_folder ".", "/vagrant"

  config.vm.provision :docker do |d|
    d.build_image "/vagrant/godep-goxc/", args: "-t godep-goxc"
    d.run "godep-goxc",
      args: [
        "--rm",
        "-v /vagrant:/gopath/src/github.com/ailispaw/talk2docker",
        "-w /gopath/src/github.com/ailispaw/talk2docker"
      ].join(" "),
      cmd: "sh -c 'godep restore && make goxc'",
      auto_assign_name: false, daemonize: false, restart: false
  end
end
