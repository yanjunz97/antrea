// Copyright 2022 Antrea Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package multicluster

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	mcsv1alpha1 "antrea.io/antrea/multicluster/apis/multicluster/v1alpha1"
	"antrea.io/antrea/pkg/antctl/raw/multicluster/common"
	mcscheme "antrea.io/antrea/pkg/antctl/raw/multicluster/scheme"
)

func TestJoin(t *testing.T) {
	existingClusterSet := &mcsv1alpha1.ClusterSet{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      "test-clusterset",
		},
		Status: mcsv1alpha1.ClusterSetStatus{
			ClusterStatuses: []mcsv1alpha1.ClusterStatus{
				{
					ClusterID: "leader-id",
					Conditions: []mcsv1alpha1.ClusterCondition{
						{
							Message: "Is the leader",
							Status:  v1.ConditionTrue,
							Type:    mcsv1alpha1.ClusterReady,
						},
					},
				},
			},
		},
	}

	secretContent := []byte(`#test file
---
apiVersion: v1
kind: Secret
metadata:
  name: token-secret
data:
  ca.crt: YWJjZAo=
  namespace: ZGVmYXVsdAo=
  token: YWJjZAo=
type: Opaque`)

	configContent := []byte(`apiVersion: multicluster.antrea.io/v1alpha1
kind: ClusterSetJoinConfig
clusterSetID: test-clusterset
clusterID: cluster-a
namespace: default
leaderClusterID: leader-id
leaderNamespace: leader-ns
leaderAPIServer: "http://localhost"
---
apiVersion: v1
kind: Secret
metadata:
  name: token-secret
data:
  ca.crt: YWJjZAo=
  namespace: ZGVmYXVsdAo=
  token: YWJjZAo=
type: Opaque`)

	tests := []struct {
		name           string
		expectedOutput string
		clusterID      string
		failureType    string
		secretFile     bool
		configFile     bool
	}{
		{
			name:           "join successfully",
			clusterID:      "cluster-a",
			expectedOutput: "Member cluster joined successfully",
		},
		{
			name:           "join successfully with Secret file",
			clusterID:      "cluster-a",
			expectedOutput: "Created the Secret from the config file",
			secretFile:     true,
		},
		{
			name:           "join successfully with config file",
			clusterID:      "cluster-a",
			expectedOutput: "Created the Secret from the config file",
			configFile:     true,
		},
		{
			name:           "fail to join due to empty ClusterID",
			clusterID:      "",
			expectedOutput: "the ClusterID of member cluster is required",
		},
		{
			name:           "fail to join and rollback",
			clusterID:      "cluster-a",
			failureType:    "create",
			expectedOutput: "failed to create object",
		},
	}
	cmd := NewJoinCommand()
	buf := new(bytes.Buffer)
	cmd.SetOutput(buf)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.Flag("clusterset").Value.Set("test-clusterset")

	joinOpts.ClusterSetID = "test-clusterset"
	joinOpts.LeaderClusterID = "leader-id"
	joinOpts.LeaderNamespace = "leader-ns"
	joinOpts.LeaderAPIServer = "http://localhost"
	joinOpts.TokenSecretName = "member-token"
	joinOpts.Namespace = "default"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			joinOpts.ClusterID = tt.clusterID
			joinOpts.k8sClient = fake.NewClientBuilder().WithScheme(mcscheme.Scheme).WithObjects(existingClusterSet).Build()
			if tt.failureType == "create" {
				joinOpts.k8sClient = common.FakeCtrlRuntimeClient{
					Client:      fake.NewClientBuilder().WithScheme(mcscheme.Scheme).WithObjects(existingClusterSet).Build(),
					ShouldError: true,
				}
			}
			if tt.secretFile {
				secret, err := os.CreateTemp("", "secret")
				if err != nil {
					log.Fatal(err)
				}
				defer os.Remove(secret.Name())
				secret.Write([]byte(secretContent))
				joinOpts.TokenSecretName = ""
				joinOpts.TokenSecretFile = secret.Name()
			}

			joinOpts.ConfigFile = ""
			if tt.configFile {
				config, err := os.CreateTemp("", "config")
				if err != nil {
					log.Fatal(err)
				}
				defer os.Remove(config.Name())
				config.Write([]byte(configContent))
				joinOpts.TokenSecretName = ""
				joinOpts.TokenSecretFile = ""
				joinOpts.ConfigFile = config.Name()
			}
			err := cmd.Execute()
			if err != nil {
				assert.Contains(t, err.Error(), tt.expectedOutput)
			} else {
				assert.Contains(t, buf.String(), tt.expectedOutput)
			}
		})
	}
}

func TestJoinOptValidate(t *testing.T) {
	tests := []struct {
		name           string
		expectedOutput string
		opts           *ClusterSetJoinConfig
		secretFile     bool
	}{
		{
			name:           "empty ClusterID",
			expectedOutput: "the ClusterID of leader cluster is required",
			opts: &ClusterSetJoinConfig{
				TokenSecretName: "token-a",
			},
		},
		{
			name:           "empty API Server",
			expectedOutput: "the API server of the leader cluster is required",
			opts: &ClusterSetJoinConfig{
				TokenSecretName: "token-a",
				ClusterID:       "cluster-a",
				LeaderClusterID: "leader-id",
			},
		},
		{
			name:           "empty Secret",
			expectedOutput: "a member token Secret must be provided through the Secret name, or Secret file, or Secret manifest in the config file",
			opts: &ClusterSetJoinConfig{
				ClusterID:       "cluster-a",
				LeaderClusterID: "leader-id",
				LeaderAPIServer: "http://localhost",
			},
		},
		{
			name:           "empty leader Namespace",
			expectedOutput: "the leader cluster Namespace is required",
			opts: &ClusterSetJoinConfig{
				ClusterID:       "cluster-a",
				LeaderClusterID: "leader-id",
				LeaderAPIServer: "http://localhost",
				TokenSecretName: "token-a",
			},
		},
		{
			name:           "empty ClusterSet ID",
			expectedOutput: "the ClusterSet ID is required",
			opts: &ClusterSetJoinConfig{
				ClusterID:       "cluster-a",
				LeaderClusterID: "leader-id",
				LeaderAPIServer: "http://localhost",
				TokenSecretName: "token-a",
				LeaderNamespace: "default",
			},
		},
		{
			name:           "empty member ClusterID",
			expectedOutput: "the ClusterID of member cluster is required",
			opts: &ClusterSetJoinConfig{
				LeaderClusterID: "leader-id",
				LeaderAPIServer: "http://localhost",
				TokenSecretName: "token-a",
				LeaderNamespace: "default",
				ClusterSetID:    "test-clusterset",
			},
		},
		{
			name:           "empty kubeconfig",
			expectedOutput: "flag accessed but not defined: kubeconfig",
			opts: &ClusterSetJoinConfig{
				LeaderClusterID: "leader-id",
				LeaderAPIServer: "http://localhost",
				TokenSecretName: "token-a",
				LeaderNamespace: "default",
				ClusterSetID:    "test-clusterset",
				ClusterID:       "cluster-a",
			},
		},
		{
			name:           "failed to unmarshal Secret file",
			expectedOutput: "failed to unmarshall Secret from token Secret file",
			opts: &ClusterSetJoinConfig{
				LeaderClusterID: "leader-id",
				LeaderAPIServer: "http://localhost",
				LeaderNamespace: "default",
				ClusterSetID:    "test-clusterset",
				ClusterID:       "cluster-a",
			},
			secretFile: true,
		},
	}

	secretContent := []byte(`apiVersion: v1
	kind: Secret
	metadata:
	  name: token-secret
	data:
	  ca.crt: a
	  namespace: ZGVmYXVsdAo=
	  token: YWJjZAo=
	type: Opaque`)
	cmd := &cobra.Command{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.secretFile {
				secret, err := os.CreateTemp("", "secret")
				if err != nil {
					log.Fatal(err)
				}
				defer os.Remove(secret.Name())
				secret.Write([]byte(secretContent))
				tt.opts.TokenSecretName = ""
				tt.opts.TokenSecretFile = secret.Name()
			}
			err := tt.opts.validateAndComplete(cmd)
			if err != nil {
				assert.Contains(t, err.Error(), tt.expectedOutput)
			} else {
				t.Error("Expected to get error but got nil")
			}
		})
	}
}
