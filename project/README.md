### WSL2でDockerのコマンドの実行に失敗する場合の対処方法

```bash
# docker daemonを起動する
sudo service docker start
```

### WSL2で`sudo apt update`に失敗する場合の対処方法

```bash
# nameserverを設定する
echo "nameserver 8.8.8.8" | sudo tee /etc/resolv.conf
```

### Dockerのコンテナの実行ユーザをホストのユーザにマッピングする

```bash
# daemon.jsonを生成する
echo '{"userns-remap" : "default"}' | sudo tee /etc/docker/daemon.json
# Dockerを再起動する
sudo systemctl restart docker
# /etc/subuidが生成されている
cat /etc/subuid
```

### WSL2でDockerをrootlessで実行する

#### rootlessモードのDockerの導入

```bash
# Docker実行用のユーザを作成する
sudo useradd -m -d /home/rootless -s /bin/bash rootless

# ユーザにパスワードを設定する
sudo passwd rootless rootless

# systemdの設定
sudo loginctl enable-linger rootless

# ユーザの切り替え
su rootless

# rootlessモードのDockerのinstall
export SKIP_IPTABLES=1
curl -fsSL https://get.docker.com/rootless | sh

# .bashrcに環境変数の設定
export PATH=/home/rootless/bin:$PATH
export DOCKER_HOST=unix:///run/user/{{rootlessのuid}}/docker.sock

# .bashrcの設定の反映
source ~/.bashrc

# ~/.config/systemd/user/docker.serviceのExecStartを書き換える
# ExecStart=/usr/bin/dockerd-rootless.sh --iptables=false
# --iptables=falseの記述を削除する
# ExecStart=/usr/bin/dockerd-rootless.sh

# 設定を反映させる
systemctl --user daemon-reload
systemctl --user restart docker

# cgroupsをv2に設定する
```