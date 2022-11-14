# k8s cluster 구성
* kubernetes 패키지 설치 , master 노드와 worker 노드 간의 연결 설정    
* docker가 설치되어 있어야함
* 기존 kubernetes가 설치된 환경이라면, kubernetes 관련 프로세스 (kube*)를 모두 제거
* 기존 kubernetes가 설치된 환경이라면, 연결이 끊어진 mount-point (/var/lib/kubelet/...)를 모두 제거


### Hostname 설정
master/worker 노드임을 알 수 있게 hostname 변경, 재로그인시 변경
```
sudo hostnamectl set-hostname master

sudo hostnamectl set-hostname worker1
...
```


### SELinux 해제
/etc/selinux/config 파일에서
```
SELINUX=disabled
```
재시작후 적용


### Firewall 해제
```
sudo ufw disable
sudo systemctl stop firewalld
sudo systemctl disable firewalld
```


### Swap 해제
```
sudo swapoff -a && sudo sed -i '/swap/s/^/#/' /etc/fstab
```


### iptables가 브리지된 트래픽을 보도록 허용
아래 명령으로 br_netfilter모듈이 로드 되었는지 확인
```
lsmod | grep br_netfilter
```

명시적으로 로드하려면 sudo modprobe br_netfilter로 모듈을 로딩
```
sudo modprobe br_netfilter
```

쿠버네티스에서 br_netfilter를 사용하도록 설정
```
cat <<EOF | sudo tee /etc/modules-load.d/k8s.conf
br_netfilter
EOF

cat <<EOF | sudo tee /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
EOF

sudo sysctl --system
```


### kubeadm, kubelet 및 kubectl 설치
모든 컴퓨터에 다음 패키지를 설치

* kubeadm: 클러스터를 부트스트랩 함.
* kubelet: 클러스터의 모든 머신에서 실행되고 포드 및 컨테이너 시작과 같은 작업을 수행하는 구성 요소.
* kubectl: 클러스터와 통신하기 위한 명령줄 util

apt 패키지 업데이트 및 k8s apt 저장소를 사용하는데 필요한 패키지 설치
```
sudo apt-get update
sudo apt-get install -y apt-transport-https ca-certificates curl
```

Google Cloud 공개 서명 키를 다운로
```
sudo curl -fsSLo /usr/share/keyrings/kubernetes-archive-keyring.gpg https://packages.cloud.google.com/apt/doc/apt-key.gpg
```

Kubernetes apt리포지토리를 추가
```
echo "deb [signed-by=/usr/share/keyrings/kubernetes-archive-keyring.gpg] https://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee /etc/apt/sources.list.d/kubernetes.list
```

apt패키지 인덱스를 업데이트 하고 kubelet, kubeadm 및 kubectl을 설치하고 해당 버전을 고정
```
sudo apt-get update
sudo apt-get install -y kubelet kubeadm kubectl
sudo apt-mark hold kubelet kubeadm kubectl
```


### kubelet cgroup 드라이버 구성 (Docker container 기반)
```
sudo mkdir /etc/docker
cat <<EOF | sudo tee /etc/docker/daemon.json
{
  "exec-opts": ["native.cgroupdriver=systemd"],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "100m"
  },
  "storage-driver": "overlay2"
}
EOF
```

Docker를 다시 시작하고 부팅 시 활성화
```
sudo systemctl enable docker
sudo systemctl daemon-reload
sudo systemctl restart docker
sudo systemctl restart kubelet
```


### kubeadm을 사용하여 kubernetes cluster 구성하기 (Master 노드에서만)
kubeadm init 명령을 통해서 클러스터를 생성. Flannel 을 Network Addon 사용
```
sudo kubeadm init --pod-network-cidr=10.244.0.0/16 --apiserver-advertise-address=<Master 노드의 IP>
```

아래와 같이 화면에 출력된 클러스터 Join 명령줄을 복사하여 워커 노드들에서 실행
```
Your Kubernetes control-plane has initialized successfully!
To start using your cluster, you need to run the following as a regular user:

  mkdir -p $HOME/.kube
  sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
  sudo chown $(id -u):$(id -g) $HOME/.kube/config

Alternatively, if you are the root user, you can run:

  export KUBECONFIG=/etc/kubernetes/admin.conf

You should now deploy a pod network to the cluster.
Run "kubectl apply -f [podnetwork].yaml" with one of the options listed at:
  https://kubernetes.io/docs/concepts/cluster-administration/addons/

Then you can join any number of worker nodes by running the following on each as root:
kubeadm join <Master 노드의 IP>:6443 --token <Token> \
        --discovery-token-ca-cert-hash sha256:<Cert-Hash>
```

### kubeadm을 사용하여 kubernetes cluster 구성하기 (Worker 노드들에서)
kubeadmin init 명령의 출력으로 나온 member join 명령줄을 복사하여 실행

```
sudo kubeadm join <Master 노드의 IP>:6443 --token <Token> \
        --discovery-token-ca-cert-hash sha256:<Cert-Hash>
```


### Node의 상태 확인 (Master 노드에서만)
아래 명령으로 모든 노드가 Ready 상태인지 확인
```
kubectl get nodes
```

NotReady 라고 나온다면, CoreDNS Pod 들이 Pending 상태. 아래 명령으로 coredns 코드를 수정
```
kubectl edit cm coredns -n kube-system
```

24라인의 loop 명령어를 주석처리
```
cache 30
#loop
reload
loadbalance
```

Flannel network addon 설치
```
kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml
```

Flannel Pod 동작 확인
```
kubectl get pods -n kube-flannel -o wide
```

잠시 대기후 노드가 Running 상태로 바뀜
```
kubectl get nodes
```



## 기타 사항

### root 계정 외에 다른 계정에서 kubectl 커맨드 사용 가능 
```
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
```


### 클러스터 조인 토큰 찾기 (기본 24 시간후 만료)
```
kubeadm token list
```

만료되었다면
```
kubeadm token create
```


### 클러스터 디스커버리 토큰 해시 찾기
```
openssl x509 -pubkey -in /etc/kubernetes/pki/ca.crt | openssl rsa -pubin -outform der 2>/dev/null | openssl dgst -sha256 -hex | sed 's/^.* //'
```


### Pod network addon (
* [마스터 노드에서만 실행] Pod가 서로 통신 할 수 있도록 CNI(Container Network Interface) 기반 Pod 네트워크 추가 기능 구성한다.

* flannel (사용)
```
# kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml

podsecuritypolicy.policy/psp.flannel.unprivileged created
clusterrole.rbac.authorization.k8s.io/flannel created
clusterrolebinding.rbac.authorization.k8s.io/flannel created
serviceaccount/flannel created
configmap/kube-flannel-cfg created
daemonset.apps/kube-flannel-ds-amd64 created
daemonset.apps/kube-flannel-ds-arm64 created
daemonset.apps/kube-flannel-ds-arm created
daemonset.apps/kube-flannel-ds-ppc64le created
daemonset.apps/kube-flannel-ds-s390x created
```
* flannel 설치 확인
```
$ kubectl get pod --namespace=kube-system -o wide
NAME                                      READY   STATUS    RESTARTS   AGE     IP            NODE               NOMINATED NODE   READINESS GATES
coredns-5644d7b6d9-wdt7j                  1/1     Running   4          7d17h   10.244.0.11   bsc-kube-master    <none>           <none>
coredns-5644d7b6d9-x97vf                  1/1     Running   4          7d17h   10.244.0.10   bsc-kube-master    <none>           <none>
etcd-bsc-kube-master                      1/1     Running   4          7d17h   10.0.2.15     bsc-kube-master    <none>           <none>
kube-apiserver-bsc-kube-master            1/1     Running   4          7d17h   10.0.2.15     bsc-kube-master    <none>           <none>
kube-controller-manager-bsc-kube-master   1/1     Running   4          7d17h   10.0.2.15     bsc-kube-master    <none>           <none>
kube-flannel-ds-amd64-n7n7q               1/1     Running   4          7d17h   10.0.2.15     bsc-kube-master    <none>           <none>
kube-flannel-ds-amd64-qbcwp               1/1     Running   5          7d17h   10.0.2.15     bsc-kube-worker    <none>           <none>
kube-flannel-ds-amd64-qts5l               1/1     Running   1          3d22h   10.0.2.15     bsc-kube-worker2   <none>           <none>
kube-proxy-9sw4k                          1/1     Running   5          7d17h   10.0.2.15     bsc-kube-worker    <none>           <none>ㅋ
kube-proxy-b5hqc                          1/1     Running   4          7d17h   10.0.2.15     bsc-kube-master    <none>           <none>
kube-proxy-m5npl                          1/1     Running   1          3d22h   10.0.2.15     bsc-kube-worker2   <none>           <none>
kube-scheduler-bsc-kube-master            1/1     Running   4          7d17h   10.0.2.15     bsc-kube-master    <none>           <none>
```
* calico  (이전에 사용했었으나 불안정해서 flannel로 교체)
```
wget https://docs.projectcalico.org/manifests/calico.yaml
kubectl apply -f calico.yaml

# 정상 running 확인
root@master:~# kubectl get pods -n kube-system
NAMESPACE     NAME                                       READY   STATUS    RESTARTS   AGE
kube-system   calico-kube-controllers-647d84984b-rbp5b   1/1     Running   0          4m45s
kube-system   calico-node-gft4v                          1/1     Running   0          4m47s
kube-system   coredns-64897985d-2ndcp                    1/1     Running   0          9m1s
kube-system   coredns-64897985d-bn7r8                    1/1     Running   0          9m2s
kube-system   etcd-master                                1/1     Running   0          9m20s
kube-system   kube-apiserver-master                      1/1     Running   0          9m28s
kube-system   kube-controller-manager-master             1/1     Running   0          9m20s
kube-system   kube-proxy-bvcx9                           1/1     Running   0          9m2s
kube-system   kube-scheduler-master                      1/1     Running   0          9m20s
```
