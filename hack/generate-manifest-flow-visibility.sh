#!/usr/bin/env bash

# Copyright 2022 Antrea Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -eo pipefail

function echoerr {
    >&2 echo "$@"
}

_usage="Usage: $0 [--mode (dev|release)] [--keep] [--help|-h]
Generate a YAML manifest for the Clickhouse-Grafana Flow-visibility Solution, using Kustomize, and
print it to stdout.
        --mode (dev|release)        Choose the configuration variant that you need (default is 'dev').
        --keep                      Debug flag which will preserve the generated kustomization.yml.
        --volume (ram|pv)           Choose the volume provider that you need (default is 'ram').
        --storageclass -sc <name>   Provide the StorageClass used to dynamically provision the 
                                    Persistent Volume for ClickHouse storage.
        --local <path>              Create the Persistent Volume for ClickHouse with a provided
                                    local path.
        --nfs <hostname:path>       Create the Persistent Volume for ClickHouse with a provided
                                    NFS server hostname or IP address and the path exported in the
                                    form of hostname:path.
        --no-ch-monitor             Generate a manifest without the ClickHouse monitor.

This tool uses kustomize (https://github.com/kubernetes-sigs/kustomize) to generate manifests for
Clickhouse-Grafana Flow-visibility Solution. You can set the KUSTOMIZE environment variable to the
path of the kustomize binary you want us to use. Otherwise we will look for kustomize in your PATH
and your GOPATH. If we cannot find kustomize there, we will try to install it."

function print_usage {
    echoerr "$_usage"
}

function print_help {
    echoerr "Try '$0 --help' for more information."
}

MODE="dev"
KEEP=false
VOLUME="ram"
STORAGECLASS=""
LOCALPATH=""
NFSPATH=""
CHMONITOR=true

while [[ $# -gt 0 ]]
do
key="$1"
case $key in
    --mode)
    MODE="$2"
    shift 2
    ;;
    --keep)
    KEEP=true
    shift
    ;;
    --volume)
    VOLUME="$2"
    shift 2
    ;;
    -sc|--storageclass)
    STORAGECLASS="$2"
    shift 2
    ;;
    --local)
    LOCALPATH="$2"
    shift 2
    ;;
    --nfs)
    NFSPATH="$2"
    shift 2
    ;;
    --no-ch-monitor)
    CHMONITOR=false
    shift 1
    ;;
    -h|--help)
    print_usage
    exit 0
    ;;
    *)    # unknown option
    echoerr "Unknown option $1"
    exit 1
    ;;
esac
done

if [ "$MODE" != "dev" ] && [ "$MODE" != "release" ]; then
    echoerr "--mode must be one of 'dev' or 'release'"
    print_help
    exit 1
fi

if [ "$MODE" == "release" ] && [ -z "$IMG_NAME" ]; then
    echoerr "In 'release' mode, environment variable IMG_NAME must be set"
    print_help
    exit 1
fi

if [ "$MODE" == "release" ] && [ -z "$IMG_TAG" ]; then
    echoerr "In 'release' mode, environment variable IMG_TAG must be set"
    print_help
    exit 1
fi

if [ "$VOLUME" != "ram" ] && [ "$VOLUME" != "pv" ]; then
    echoerr "--volume must be one of 'ram' or 'pv'"
    print_help
    exit 1
fi

if [ "$VOLUME" == "pv" ] && [ "$LOCALPATH" == "" ] && [ "$NFSPATH" == "" ] && [ "$STORAGECLASS" == "" ]; then
    echoerr "When deploy with 'pv', one of '--local', '--nfs', '--storageclass' should be set"
    print_help
    exit 1
fi

if [ "$LOCALPATH" != "" ] && [ "$NFSPATH" != "" ]; then
    echoerr "Cannot set '--local' and '--nfs' at the same time"
    print_help
    exit 1
fi

if [ "$NFSPATH" != "" ]; then
    pathPair=(${NFSPATH//:/ })
    if [ ${#pathPair[@]} != 2 ]; then
        echoerr "--nfs must be in the form of hostname:path"
        print_help
        exit 1
    fi
fi

THIS_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

source $THIS_DIR/verify-kustomize.sh

if [ -z "$KUSTOMIZE" ]; then
    KUSTOMIZE="$(verify_kustomize)"
elif ! $KUSTOMIZE version > /dev/null 2>&1; then
    echoerr "$KUSTOMIZE does not appear to be a valid kustomize binary"
    print_help
    exit 1
fi

KUSTOMIZATION_DIR=$THIS_DIR/../build/yamls/flow-visibility

TMP_DIR=$(mktemp -d $KUSTOMIZATION_DIR/overlays.XXXXXXXX)

pushd $TMP_DIR > /dev/null

BASE=../../base


if $CHMONITOR; then
    mkdir chmonitor && cd chmonitor
    cp $KUSTOMIZATION_DIR/patches/chmonitor/*.yml .
    touch kustomization.yml
    $KUSTOMIZE edit add base $BASE
    $KUSTOMIZE edit add patch --path chMonitor.yml --group clickhouse.altinity.com --version v1 --kind ClickHouseInstallation --name clickhouse
    BASE=../chmonitor
    cd ..
    mkdir $MODE && cd $MODE
    touch kustomization.yml
    $KUSTOMIZE edit add base $BASE
    # ../../patches/$MODE may be empty so we use find and not simply cp
    find ../../patches/$MODE -name \*.yml -exec cp {} . \;

    if [ "$MODE" == "dev" ]; then
        $KUSTOMIZE edit set image flow-visibility-clickhouse-monitor=projects.registry.vmware.com/antrea/flow-visibility-clickhouse-monitor:latest
        $KUSTOMIZE edit add patch --path imagePullPolicy.yml --group clickhouse.altinity.com --version v1 --kind ClickHouseInstallation --name clickhouse
    fi

    if [ "$MODE" == "release" ]; then
        $KUSTOMIZE edit set image flow-visibility-clickhouse-monitor=$IMG_NAME:$IMG_TAG
    fi
    BASE=../$MODE
    cd ..
fi

if [ "$VOLUME" == "ram" ]; then
    mkdir ram && cd ram
    cp $KUSTOMIZATION_DIR/patches/ram/*.yml .
    touch kustomization.yml
    $KUSTOMIZE edit add base $BASE
    $KUSTOMIZE edit add patch --path mountRam.yml --group clickhouse.altinity.com --version v1 --kind ClickHouseInstallation --name clickhouse
fi

if [ "$VOLUME" == "pv" ]; then
    mkdir pv && cd pv
    cp $KUSTOMIZATION_DIR/patches/pv/*.yml .
    touch kustomization.yml
    $KUSTOMIZE edit add base $BASE

    if [[ $STORAGECLASS != "" ]]; then
        sed -i.bak -E "s/STORAGECLASS_NAME/$STORAGECLASS/" mountPv.yml
    else
        sed -i.bak -E "s/STORAGECLASS_NAME/clickhouse-storage/" mountPv.yml
    fi
    if [[ $LOCALPATH != "" ]]; then
        sed -i.bak -E "s~LOCAL_PATH~$LOCALPATH~" createLocalPv.yml
        $KUSTOMIZE edit add base createLocalPv.yml
    fi
    if [[ $NFSPATH != "" ]]; then
        sed -i.bak -E "s~NFS_SERVER_PATH~${pathPair[0]}~" createNfsPv.yml
        sed -i.bak -E "s~NFS_SERVER_ADDRESS~${pathPair[1]}~" createNfsPv.yml
        $KUSTOMIZE edit add base createNfsPv.yml
    fi
    $KUSTOMIZE edit add patch --path mountPv.yml --group clickhouse.altinity.com --version v1 --kind ClickHouseInstallation --name clickhouse
fi

$KUSTOMIZE build

popd > /dev/null


if $KEEP; then
    echoerr "Kustomization file is at $TMP_DIR/$MODE/kustomization.yml"
else
    rm -rf $TMP_DIR
fi
