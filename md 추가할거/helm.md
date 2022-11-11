# helm , helm chart 설치 

* kubernetes가 설치되어 있어야함



### helm 설치 및 버전 확인
linux ubuntu 20.04 환경에서 설치 
아래의 커맨드로 파일 다운로드 후 스크립트 실행 
```
$ curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
$ chmod 700 get_helm.sh
$ ./get_helm.sh
```

아래 명령어를 통해 설치됨과 버전을 확인할 수 있음.

```
$ root@master:~$ helm version
version.BuildInfo{Version:"v3.9.2", GitCommit:"1addefbfe665c350f4daf868a9adc5600cc064fd", GitTreeState:"clean", GoVersion:"go1.17.12"}
```

### helm 사용하기

helm에서 chart를 사용하기위해 미리 배포된 repository를 사용하기 위해선 repo를 추가해야 한다. repo add 명령어로 repository를 다운받을 수 있다. 
아래 list 명령어로 실행중인 chart의 리스트를 확인할 수 있다.이 때 namespace를 확인하고 기입해야 한다. 
```
$ helm repo add [NAME] [URL] [flags]
```

list 명령어로 실행중인 chart의 리스트를 확인할 수 있다. 이때 namespace를 확인하고 기입해야 한다. 
```
$ root@master:~$ helm list -A
NAME            NAMESPACE       REVISION        UPDATED                                 STATUS          CHART                            APP VERSION
jsonexporter    ksv             6               2022-11-01 23:06:29.259408297 +0900 KST deployed        prometheus-json-exporter-0.4.0   v0.5.0
prometheus      ksv             4               2022-10-31 20:19:34.776034436 +0900 KST deployed        kube-prometheus-stack-39.11.0    0.58.0
pvc-autoresizer ksv             1               2022-10-24 16:11:48.440464113 +0900 KST deployed        pvc-autoresizer-0.5.0            0.5.0     
```



