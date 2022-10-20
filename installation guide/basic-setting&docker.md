## basic setting and docker install

### SSH 및 방화벽 해제
```
apt-get update
(sudo) apt-get update

ssh openssh-server 설치
(sudo) apt-get install openssh-server

ssh 클라이언트와 서버를 동시에 설치
(sudo) apt-get install ssh

```

```
(sudo) ufw enable
(sudo) ufw allow 22
(sudo) ufw reload

ssh 서비스 시작

(sudo) service ssh start

제대로 구동 되었는지 확인


(sudo) service ssh status
(sudo) ps -ef | grep sshd
(sudo) netstat -ntlp | grep sshd

```
### Docker

```
https를 사용해서 레포지토리를 사용할 수 있도록 필요한 패키지를 설치한다.
$ sudo apt-get install -y  apt-transport-https ca-certificates curl software-properties-common

Docker 공식 리포지토리에서 패키지를 다운로드 받았을 때 위변조 확인을 위한 GPG 키를 추가
$ curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add

Docker.com 의 GPG 키가 등록됐는지 확인한다
$ apt-key fingerprint

Docker 공식 저장소를 리포지토리로 등록한다
$ add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable”

저장소 등록정보에 기록됐는지 확인한다
$ grep docker /etc/apt/sources.list
deb [arch=amd64] https://download.docker.com/linux/ubuntu bionic stable

리포지토리 정보를 갱신
$ sudo-apt update

docker container engine 을 설치한다
$ apt-get install -y docker-ce

도커 서비스 상태, 버전 확인
$ ps -ef | grep docker

root     14503     1  0 00:58 ?        00:00:00 /usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock
root     17393  5202  0 01:03 pts/1    00:00:00 grep --color=auto docker
$ docker --version
Docker version 20.10.12, build e91ed57
```

[참고](https://kindloveit.tistory.com/18)