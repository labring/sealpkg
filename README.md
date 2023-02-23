# runtime-ctl


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

5. goproxy

https://ghproxy.com/