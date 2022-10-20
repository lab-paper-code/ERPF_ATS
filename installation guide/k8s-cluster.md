# k8s cluster 구성
* kubernetes 패키지 설치 , master 노드와 worker 노드 간의 연결 설정    
* docker가 설치되어 있어야함
<br/>

### MAC 주소와 product_uuid가 모든 노드에 대해 고유한지 확인
* ifconfig -a 로 각 노드에 할당된 Mac address 및 IP를 조회하여 충돌이 나지 않는지 확인
```
ifconfig -a

$ produc_uuide 확인
cat /sys/class/dmi/id/product_uuid
```

### iptables가 브리지된 트래픽을 보도록 허용
* lsmod | grep br_netfilter br_netfilter모듈이 로드 되었는지 확인,
   명시적으로 로드하려면 sudo modprobe br_netfilter로 모듈을 로딩
   ```
   #sudo modprobe br_netfilter
   #lsmod | grep br

   br_netfilter           28672  0
   bridge                307200  1 br_netfilter
   stp                    16384  1 bridge
   llc                    16384  2 bridge,stp

   ```
* Linux 노드의 iptables가 브리지된 트래픽을 올바르게 보기 위한 요구 사항으로 구성 net.bridge.bridge-nf-call-iptables에서 가 1로 설정되어 있는지 확인
```
$cat <<EOF | sudo tee /etc/modules-load.d/k8s.conf
>br_netfilter
>EOF

$cat <<EOF | sudo tee /etc/sysctl.d/k8s.conf
>net.bridge.bridge-nf-call-ip6tables = 1
>net.bridge.bridge-nf-call-iptables = 1
>EOF
sudo sysctl --system

# 설정이 안되어 있다면
$sysctl net.bridge.bridge-nf-call-iptables=1
```

### Hostname 변경
* master 노드임을 알 수 있게 hostname 변경, 재로그인시 변경
```
$hostnamectl set-hostname master
```

### SELinux, firewall 해제
```
$ setenforce 0
setenforce: SELinux is disabled
$  ufw disable
```

### Swap 해제
```
$ swapon && cat /etc/fstab
$ swapoff -a && sed -i '/swap/s/^/#/' /etc/fstab
```

### kubeadm, kubelet 및 kubectl 설치
모든 컴퓨터에 다음 패키지를 설치

* kubeadm: 클러스터를 부트스트랩 함.
* kubelet: 클러스터의 모든 머신에서 실행되고 포드 및 컨테이너 시작과 같은 작업을 수행하는 구성 요소.
* kubectl: 클러스터와 통신하기 위한 명령줄 util

* apt 패키지 업데이트 및 k8s apt 저장소를 사용하는데 필요한 패키지 설치
```
sudo apt-get update
sudo apt-get install -y apt-transport-https ca-certificates curl
```

* Google Cloud 공개 서명 키를 다운로
```
sudo curl -fsSLo /usr/share/keyrings/kubernetes-archive-keyring.gpg https://packages.cloud.google.com/apt/doc/apt-key.gpg
```

* Kubernetes apt리포지토리를 추가
```
echo "deb [signed-by=/usr/share/keyrings/kubernetes-archive-keyring.gpg] https://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee /etc/apt/sources.list.d/kubernetes.list
```

* apt패키지 인덱스를 업데이트 하고 kubelet, kubeadm 및 kubectl을 설치하고 해당 버전을 고정
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

* Docker를 다시 시작하고 부팅 시 활성화
```
sudo systemctl enable docker
sudo systemctl daemon-reload
sudo systemctl restart docker
```

### kubeadm을 사용하여 kubernetes cluster 구성하기
* [마스터 노드에서만 실행] kubeadm init 명령을 통해서 클러스터를 생성
```
$ kubeadm init

[preflight] Running pre-flight checks
[preflight] Pulling images required for setting up a Kubernetes cluster
[preflight] This might take a minute or two, depending on the speed of your internet connection
[preflight] You can also perform this action in beforehand using 'kubeadm config images pull'
[certs] Using certificateDir folder "/etc/kubernetes/pki"
[certs] Generating "ca" certificate and key
[certs] Generating "apiserver" certificate and key
[certs] apiserver serving cert is signed for DNS names [kubernetes kubernetes.default kubernetes.default.svc kubernetes.default.svc.cluster.local master] and IPs [10.96.0.1 172.30.1.16]
[certs] Generating "apiserver-kubelet-client" certificate and key
[certs] Generating "front-proxy-ca" certificate and key
[certs] Generating "front-proxy-client" certificate and key
[certs] Generating "etcd/ca" certificate and key
[certs] Generating "etcd/server" certificate and key
[certs] etcd/server serving cert is signed for DNS names [localhost master] and IPs [172.30.1.16 127.0.0.1 ::1]
[certs] Generating "etcd/peer" certificate and key
[certs] etcd/peer serving cert is signed for DNS names [localhost master] and IPs [172.30.1.16 127.0.0.1 ::1]
[certs] Generating "etcd/healthcheck-client" certificate and key
[certs] Generating "apiserver-etcd-client" certificate and key
[certs] Generating "sa" key and public key
[kubeconfig] Using kubeconfig folder "/etc/kubernetes"
[kubeconfig] Writing "admin.conf" kubeconfig file
[kubeconfig] Writing "kubelet.conf" kubeconfig file
[kubeconfig] Writing "controller-manager.conf" kubeconfig file
[kubeconfig] Writing "scheduler.conf" kubeconfig file
[kubelet-start] Writing kubelet environment file with flags to file "/var/lib/kubelet/kubeadm-flags.env"
[kubelet-start] Writing kubelet configuration to file "/var/lib/kubelet/config.yaml"
[kubelet-start] Starting the kubelet
[control-plane] Using manifest folder "/etc/kubernetes/manifests"
[control-plane] Creating static Pod manifest for "kube-apiserver"
[control-plane] Creating static Pod manifest for "kube-controller-manager"
[control-plane] Creating static Pod manifest for "kube-scheduler"
[etcd] Creating static Pod manifest for local etcd in "/etc/kubernetes/manifests"
[wait-control-plane] Waiting for the kubelet to boot up the control plane as static Pods from directory "/etc/kubernetes/manifests". This can take up to 4m0s

[kubelet-check] Initial timeout of 40s passed.
[apiclient] All control plane components are healthy after 68.017142 seconds
[upload-config] Storing the configuration used in ConfigMap "kubeadm-config" in the "kube-system" Namespace
[kubelet] Creating a ConfigMap "kubelet-config-1.23" in namespace kube-system with the configuration for the kubelets in the cluster
NOTE: The "kubelet-config-1.23" naming of the kubelet ConfigMap is deprecated. Once the UnversionedKubeletConfigMap feature gate graduates to Beta the default name will become just "kubelet-config". Kubeadm upgrade will handle this transition transparently.
[upload-certs] Skipping phase. Please see --upload-certs
[mark-control-plane] Marking the node master as control-plane by adding the labels: [node-role.kubernetes.io/master(deprecated) node-role.kubernetes.io/control-plane node.kubernetes.io/exclude-from-external-load-balancers]
[mark-control-plane] Marking the node master as control-plane by adding the taints [node-role.kubernetes.io/master:NoSchedule]
[bootstrap-token] Using token: 5v5tqi.9fbmk1t5bq4rr89l
[bootstrap-token] Configuring bootstrap tokens, cluster-info ConfigMap, RBAC Roles
[bootstrap-token] configured RBAC rules to allow Node Bootstrap tokens to get nodes
[bootstrap-token] configured RBAC rules to allow Node Bootstrap tokens to post CSRs in order for nodes to get long term certificate credentials
[bootstrap-token] configured RBAC rules to allow the csrapprover controller automatically approve CSRs from a Node Bootstrap Token
[bootstrap-token] configured RBAC rules to allow certificate rotation for all node client certificates in the cluster
[bootstrap-token] Creating the "cluster-info" ConfigMap in the "kube-public" namespace
[kubelet-finalize] Updating "/etc/kubernetes/kubelet.conf" to point to a rotatable kubelet client certificate and key
[addons] Applied essential addon: CoreDNS
[addons] Applied essential addon: kube-proxy

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
kubeadm join 172.30.1.16:6443 --token 5v5tqi.9fbmk1t5bq4rr89l \
        --discovery-token-ca-cert-hash sha256:e5d020b5cd0ba5f0a45b9482f8b41fdc1704e9a8fbc88c961b43b986cfc57c38
```





### etc
* root 계정 외에 다른 계정에서 kubectl 커맨드 사용 가능 
```
$ mkdir -p $HOME/.kube
$ sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
$ sudo chown $(id -u):$(id -g) $HOME/.kube/config
```

* [마스터 노드에서만 실행] 쿠버네티스 클러스터에 조인하기 위한 명령어 구문을 저장 
```
$ cat > token.sh
kubeadm join 172.30.1.16:6443 --token 5v5tqi.9fbmk1t5bq4rr89l \
        --discovery-token-ca-cert-hash sha256:e5d020b5cd0ba5f0a45b9482f8b41fdc1704e9a8fbc88c961b43b986cfc57c38
        
$ chmod +x token.sh (실행권한 부여)
```

### Pod network addon
* [마스터 노드에서만 실행] Pod가 서로 통신 할 수 있도록 CNI(Container Network Interface) 기반 Pod 네트워크 추가 기능 구성한다.

* flannel (사용)
```
# kubectl apply -f <https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml>

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

### 워커 노드 세팅

* kubeadm init 후에 저장한 token.sh을 각 worker 노드로 전송한다.
```
root@master:~# scp token.sh ubuntu@155.230.35.174:~/palisade3@155.230.35.174's password:
token.sh

root@master:~# scp token.sh ubuntu@155.230.35.175:~/
palisade4@155.230.35.175's password:
token.sh

root@master:~# scp token.sh ubuntu@155.230.35.179:~/
palisade5@155.230.35.175's password:
token.sh
```

* token.sh을 전달받은 각 노드는 실행해서 cluster에 연결
```
root@worker2:/home/palisade3# ./token.sh
[preflight] Running pre-flight checks
[preflight] Reading configuration from the cluster...
[preflight] FYI: You can look at this config file with 'kubectl -n kube-system get cm kubeadm-config -o yaml'
W0106 13:23:42.637670   15397 utils.go:69] The recommended value for "resolvConf" in "KubeletConfiguration" is: /run/systemd/resolve/resolv.conf; the provided value is: /run/systemd/resolve/resolv.conf
[kubelet-start] Writing kubelet configuration to file "/var/lib/kubelet/config.yaml"
[kubelet-start] Writing kubelet environment file with flags to file "/var/lib/kubelet/kubeadm-flags.env"
[kubelet-start] Starting the kubelet
[kubelet-start] Waiting for the kubelet to perform the TLS Bootstrap...

This node has joined the cluster:
* Certificate signing request was sent to apiserver and a response was received.
* The Kubelet was informed of the new secure connection details.

Run 'kubectl get nodes' on the control-plane to see this node join the cluster.
```

* Master 노드에서 최종 cluster에 붙은 모든 노드들을 확인
```
root@master:/home/palisade2# kubectl get nodes -o wide
NAME      STATUS   ROLES           AGE   VERSION   INTERNAL-IP      EXTERNAL-IP   OS-IMAGE             KERNEL-VERSION      CONTAINER-RUNTIME
master    Ready    control-plane   18h   v1.25.3   155.230.36.27    <none>        Ubuntu 20.04.4 LTS   5.15.0-50-generic   containerd://1.6.6
worker1   Ready    <none>          18h   v1.25.3   155.230.35.174   <none>        Ubuntu 20.04.4 LTS   5.15.0-50-generic   containerd://1.6.6
worker2   Ready    <none>          18h   v1.25.3   155.230.35.175   <none>        Ubuntu 20.04.4 LTS   5.15.0-50-generic   containerd://1.6.6
worker3   Ready    <none>          18h   v1.25.3   155.230.35.179   <none>        Ubuntu 20.04.4 LTS   5.15.0-50-generic   containerd://1.6.6

```




[참고(calico)](https://kindloveit.tistory.com/23) 
[참고(flannnel)](https://medium.com/finda-tech/overview-8d169b2a54ff)