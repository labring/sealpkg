# runtime-ctl

1. 1.26
```shell
# Check support for kube-v1.26+
if [[ "${KUBE_XY//./}" -ge 126 ]] && [[ "${SEALOS_XYZ//./}" -le 413 ]] && [[ -z "$sealosPatch" ]]; then
  echo "INFO::skip $KUBE(kube>=1.26) when $SEALOS(sealos<=4.1.3)"
  echo https://kubernetes.io/blog/2022/11/18/upcoming-changes-in-kubernetes-1-26/#cri-api-removal
  exit
fi
```
2. docker
```shell
IMAGE_KUBE=kubernetes-$CRI_TYPE
  if [[ "${KUBE_XY//./}" -ge 126 ]]; then
    sudo cp -au "$MOUNT_CRI"/cri/cri-dockerd.tgz cri/
  else
    sudo cp -au "$MOUNT_CRI"/cri/cri-dockerd.tgzv125 cri/cri-dockerd.tgz
  fi
  DOCKER_XY=$(until curl -sL "https://github.com/kubernetes/kubernetes/raw/release-$KUBE_XY/build/dependencies.yaml" | yq '.dependencies[]|select(.name == "docker")|.version'; do sleep 3; done)
  case $DOCKER_XY in
  18.09 | 19.03 | 20.10)
    sudo cp -au "$MOUNT_CRI/cri/docker-$DOCKER_XY.tgz" cri/docker.tgz
    ;;
  *)
    sudo cp -au "$MOUNT_CRI/cri/docker.tgz" cri/
    ;;
  esac
```
3. manifest
```shell
if [[ "${kube_major//./}" -ge 126 ]]; then
  if ! [[ "${sealos_major//./}" -le 413 ]] || [[ -n "$sealosPatch" ]]; then
    echo "Verifying the availability of unstable"
  else
    echo "INFO::skip kube(>=1.26) building when sealos <= 4.1.3"
    exit
  fi
  FROM_CRI=$(sudo buildah from "$IMAGE_CACHE_NAME:cri-amd64")
  MOUNT_CRI=$(sudo buildah mount "$FROM_CRI")
  case $CRI_TYPE in
  containerd)
    if ! [[ "$(sudo cat "$MOUNT_CRI"/cri/.versions | grep CONTAINERD | awk -F= '{print $NF}')" =~ v1\.([6-9]|[0-9][0-9])\.[0-9]+ ]]; then
      echo https://kubernetes.io/blog/2022/11/18/upcoming-changes-in-kubernetes-1-26/#cri-api-removal
      exit
    fi
    ;;
  docker)
    if ! [[ "$(sudo cat "$MOUNT_CRI"/cri/.versions | grep CRIDOCKER | awk -F= '{print $NF}')" =~ v0\.[3-9]\.[0-9]+ ]]; then
      echo https://github.com/Mirantis/cri-dockerd/issues/125
      exit
    fi
    ;;
  esac
fi
```

4. docker
```shell
# kube 1.16(docker-18.09)
  until curl -sLo docker-18.09.tgz "https://download.docker.com/linux/static/stable/$ALIAS_ARCH/docker-18.09.9.tgz"; do sleep 3; done
  # kube 1.17-20(docker-19.03)
  until curl -sLo docker-19.03.tgz "https://download.docker.com/linux/static/stable/$ALIAS_ARCH/docker-19.03.15.tgz"; do sleep 3; done
  # kube 1.21-23(docker-20.10)
  until curl -sLo docker-20.10.tgz "https://download.docker.com/linux/static/stable/$ALIAS_ARCH/docker-20.10.23.tgz"; do sleep 3; done
  # kube 1.1x-25(cri-dockerd v0.2.x)
  # kube 1.26-2x(cri-dockerd v0.3.x)
  until curl -sLo "docker.tgz" "https://download.docker.com/linux/static/stable/$ALIAS_ARCH/docker-$DOCKER.tgz"; do sleep 3; done
```