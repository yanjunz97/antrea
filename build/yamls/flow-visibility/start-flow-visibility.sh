#!/usr/bin/env bash

# Copyright 2022 Antrea Authors
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

set -eo pipefail

THIS_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

start_flow_visibility() {
  echo "=== Starting Flow Visibility ==="
  # install Clickhouse operator
  kubectl apply -f https://raw.githubusercontent.com/Altinity/clickhouse-operator/master/deploy/operator/clickhouse-operator-install-bundle.yaml
  kubectl create namespace flow-visibility
  kubectl create configmap grafana-dashboard-config -n flow-visibility --from-file=${THIS_DIR}/grafana/provisioning/dashboards/
  kubectl apply -f ${THIS_DIR}/flow-visibility.yaml -n flow-visibility

  echo "=== Waiting for Clickhouse and Grafana to be ready ==="
  sleep 10
  kubectl wait --for=condition=ready pod -l app=clickhouse-operator -n kube-system --timeout=60s
  kubectl wait --for=condition=ready pod -l app=clickhouse -n flow-visibility --timeout=60s
  kubectl wait --for=condition=ready pod -l app=grafana -n flow-visibility --timeout=60s

  # get NodeName of Grafana Pod
  NODE_NAME=$(kubectl get pod -l app=grafana -n flow-visibility -o jsonpath='{.items[0].spec.nodeName}')
  # get NodeIP
  NODE_IP=$(kubectl get nodes ${NODE_NAME} -o jsonpath='{.status.addresses[0].address}')
  # get NodePort of Grafana Service
  GRAFANA_NODEPORT=$(kubectl get svc grafana -n flow-visibility -o jsonpath='{.spec.ports[*].nodePort}')
  echo "=== Grafana Service is listening on ${NODE_IP}:${GRAFANA_NODEPORT} ==="
}

start_flow_visibility