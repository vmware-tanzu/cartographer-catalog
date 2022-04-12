#!/usr/bin/env bash
# Copyright 2022 VMware
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset

readonly ROOT=$(cd $(dirname $0)/.. && pwd)

readonly CARTOGRAPHER_VERSION=0.3.0
readonly CERT_MANAGER_VERSION=1.7.2
readonly KAPP_CONTROLLER_VERSION=0.34.0
readonly KNATIVE_SERVING_VERSION=1.3.0
readonly KPACK_VERSION=0.5.2
readonly SECRETGEN_CONTROLLER_VERSION=0.8.0
readonly SOURCE_CONTROLLER_VERSION=0.22.4
readonly CARTOGRAPHER_CONVENTIONS_VERSION=0.1.0-build.3

main() {
        cd $ROOT/hack

        test $# -eq 0 && abort_with_help

        for cmd in $@; do
                case $cmd in
                start)
                        start_local_registry
                        start_kind_cluster
                        ;;

                apply-dependencies)
                        install_cert_manager
                        install_kapp_controller
                        install_secretgen_controller

                        install_cartographer
                        install_cartographer_conventions
                        install_knative_serving
                        install_kpack
                        install_source_controller
                        ;;

                apply-ootb)
                        install_ootb
                        ;;

                *)
                        abort_with_help
                        ;;
                esac
        done
}

abort_with_help() {
        echo "usage: $0 [cmd ...]"
        echo "cmd: (start|apply-dependencies|apply-ootb)"
        exit 1
}

start_local_registry() {
        local container_name=registry

        docker container inspect $container_name &>/dev/null && {
                echo "registry already exists"
                return
        }

        docker run \
                --detach \
                --name $container_name \
                --publish 5000:5000 \
                registry:2
}

start_kind_cluster() {
        local container_name="kind-control-plane"
        local image="kindest/node:v1.23.4@sha256:0e34f0d0fd448aa2f2819cfd74e99fe5793a6e4938b328f657c8e3f81ee0dfb9"
        local local_registry

        local_registry=$(local_ip_addr):5000

        docker container inspect $container_name &>/dev/null && {
                echo "cluster already exists"
                return
        }

        cat <<EOF | kind create cluster --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: kind
containerdConfigPatches:
- |-
  [plugins."io.containerd.grpc.v1.cri".registry]
    [plugins."io.containerd.grpc.v1.cri".registry.mirrors]
      [plugins."io.containerd.grpc.v1.cri".registry.mirrors."${local_registry}"]
        endpoint = ["http://${local_registry}"]
    [plugins."io.containerd.grpc.v1.cri".registry.configs]
      [plugins."io.containerd.grpc.v1.cri".registry.configs."${local_registry}".tls]
        insecure_skip_verify = true
nodes:
  - role: control-plane
    image: $image
EOF

        kubectl config set-context --current --namespace default
        kubectl config get-contexts
        kubectl cluster-info
}

install_ootb() {
        local local_registry
        local_registry=$(local_ip_addr):5000

        ytt --ignore-unknown-comments \
                -f $ROOT/src \
                --data-value registry.server=$local_registry \
                --data-value registry.repository=test |
                kapp deploy --yes -a ootb -f-
}

install_cartographer() {
        ytt --ignore-unknown-comments \
                -f "./overlays/strip-resources.yaml" \
                -f https://github.com/vmware-tanzu/cartographer/releases/download/v$CARTOGRAPHER_VERSION/cartographer.yaml |
                kapp deploy --yes -a cartographer -f-
}

install_cartographer_conventions() {
        ytt --ignore-unknown-comments \
                -f https://github.com/vmware-tanzu/cartographer-conventions/releases/download/v$CARTOGRAPHER_CONVENTIONS_VERSION/cartographer-conventions-v$CARTOGRAPHER_CONVENTIONS_VERSION.yaml \
                -f "./overlays/strip-resources.yaml" |
                kapp deploy --yes -a cartographer-conventions -f-
}

install_cert_manager() {
        ytt --ignore-unknown-comments \
                -f "./overlays/strip-resources.yaml" \
                -f https://github.com/jetstack/cert-manager/releases/download/v$CERT_MANAGER_VERSION/cert-manager.yaml |
                kapp deploy --yes -a cert-manager -f-
}

install_kapp_controller() {
        ytt --ignore-unknown-comments \
                -f "./overlays/strip-resources.yaml" \
                -f https://github.com/vmware-tanzu/carvel-kapp-controller/releases/download/v$KAPP_CONTROLLER_VERSION/release.yml |
                kapp deploy --yes -a kapp-controller -f-
}

install_secretgen_controller() {
        ytt --ignore-unknown-comments \
                -f "./overlays/strip-resources.yaml" \
                -f https://github.com/vmware-tanzu/carvel-secretgen-controller/releases/download/v$SECRETGEN_CONTROLLER_VERSION/release.yml |
                kapp deploy --yes -a secretgen-controller -f-
}

install_source_controller() {
        kubectl create namespace gitops-toolkit || true

        kubectl create clusterrolebinding gitops-toolkit-admin \
                --clusterrole=cluster-admin \
                --serviceaccount=gitops-toolkit:default || true

        ytt --ignore-unknown-comments \
                -f "./overlays/strip-resources.yaml" \
                -f https://github.com/fluxcd/source-controller/releases/download/v$SOURCE_CONTROLLER_VERSION/source-controller.crds.yaml \
                -f https://github.com/fluxcd/source-controller/releases/download/v$SOURCE_CONTROLLER_VERSION/source-controller.deployment.yaml |
                kapp deploy --yes -a gitops-toolkit --into-ns gitops-toolkit -f-
}

install_kpack() {
        local local_registry
        local_registry=$(local_ip_addr):5000

        kapp deploy --yes -a kpack \
                -f <(
                        ytt \
                                --ignore-unknown-comments \
                                -f https://github.com/pivotal/kpack/releases/download/v$KPACK_VERSION/release-$KPACK_VERSION.yaml \
                                -f ./overlays/strip-resources.yaml \
                                -f ./kpack --data-value builder_image=$local_registry/builder
                )

        echo "waiting clusterbuilder to be ready..."
        kubectl wait --for=condition=ready clusterbuilder default --timeout=2m
}

install_knative_serving() {
        ytt --ignore-unknown-comments \
                -f https://github.com/knative/serving/releases/download/knative-v$KNATIVE_SERVING_VERSION/serving-core.yaml \
                -f https://github.com/knative/serving/releases/download/knative-v$KNATIVE_SERVING_VERSION/serving-crds.yaml \
                -f "./overlays/strip-resources.yaml" |
                kapp deploy --yes -a knative-serving -f-
}

local_ip_addr() {
        python - <<-EOF
	import socket

	s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
	s.connect(("8.8.8.8", 80))
	print(s.getsockname()[0])
	s.close()
	EOF
}

main "$@"
