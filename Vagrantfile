Vagrant.configure(2) do |config|
  config.vm.define "rancheros-lite"

  config.vm.box = "ailispaw/rancheros-lite"

  config.vm.network :forwarded_port, guest: 2375, host: 2375, auto_correct: true, disabled: true

  config.vm.synced_folder ".", "/vagrant"

  if Vagrant.has_plugin?("vagrant-triggers") then
    config.trigger.after [:up, :resume] do
      info "Adjusting datetime after suspend and resume."
      run_remote "sudo ntpd -n -q -g -I eth0 > /dev/null; date"
    end
  end

  # Adjusting datetime before provisioning.
  config.vm.provision :shell, run: "always" do |sh|
    sh.inline = "ntpd -n -q -g -I eth0 > /dev/null; date"
  end

  config.vm.provision :docker do |d|
    d.build_image "/vagrant/godep-goxc/", args: "-t godep-goxc"
    d.run "godep-goxc",
      args: "--rm -v /vagrant:/gopath/src/github.com/ailispaw/talk2docker -w /gopath/src/github.com/ailispaw/talk2docker",
      cmd: "sh -c 'godep restore && make goxc'",
      auto_assign_name: false, daemonize: false
  end
end
