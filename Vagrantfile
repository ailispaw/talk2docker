Vagrant.configure(2) do |config|
  config.vm.define "boot2docker"

  config.vm.box = "yungsang/boot2docker"

  config.vm.network "private_network", ip: "192.168.33.10"

  config.vm.synced_folder ".", "/vagrant", type: "nfs"

  if Vagrant.has_plugin?("vagrant-triggers") then
    config.trigger.after [:up, :resume] do
      info "Adjusting datetime after suspend and resume."
      run_remote "timeout -t 10 sudo /usr/local/bin/ntpclient -s -h pool.ntp.org; date"
    end
  end

  # Adjusting datetime before provisioning.
  config.vm.provision :shell, run: "always" do |sh|
    sh.inline = "timeout -t 10 sudo /usr/local/bin/ntpclient -s -h pool.ntp.org; date"
  end

  config.vm.provision :docker do |d|
    d.build_image "/vagrant/godep-goxc/", args: "-t godep-goxc"
    d.run "godep-goxc",
      args: "--rm -v /vagrant:/gopath/src/github.com/ailispaw/talk2docker -w /gopath/src/github.com/ailispaw/talk2docker",
      cmd: "sh -c 'godep restore && make goxc'",
      auto_assign_name: false, daemonize: false
  end
end
