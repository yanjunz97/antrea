// Copyright 2022 Antrea Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Code generated by MockGen. DO NOT EDIT.
// Source: antrea.io/antrea/pkg/flowaggregator/querier (interfaces: FlowAggregatorQuerier)

// Package testing is a generated GoMock package.
package testing

import (
	querier "antrea.io/antrea/pkg/flowaggregator/querier"
	gomock "github.com/golang/mock/gomock"
	intermediate "github.com/vmware/go-ipfix/pkg/intermediate"
	reflect "reflect"
)

// MockFlowAggregatorQuerier is a mock of FlowAggregatorQuerier interface
type MockFlowAggregatorQuerier struct {
	ctrl     *gomock.Controller
	recorder *MockFlowAggregatorQuerierMockRecorder
}

// MockFlowAggregatorQuerierMockRecorder is the mock recorder for MockFlowAggregatorQuerier
type MockFlowAggregatorQuerierMockRecorder struct {
	mock *MockFlowAggregatorQuerier
}

// NewMockFlowAggregatorQuerier creates a new mock instance
func NewMockFlowAggregatorQuerier(ctrl *gomock.Controller) *MockFlowAggregatorQuerier {
	mock := &MockFlowAggregatorQuerier{ctrl: ctrl}
	mock.recorder = &MockFlowAggregatorQuerierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFlowAggregatorQuerier) EXPECT() *MockFlowAggregatorQuerierMockRecorder {
	return m.recorder
}

// GetFlowRecords mocks base method
func (m *MockFlowAggregatorQuerier) GetFlowRecords(arg0 *intermediate.FlowKey) []map[string]interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFlowRecords", arg0)
	ret0, _ := ret[0].([]map[string]interface{})
	return ret0
}

// GetFlowRecords indicates an expected call of GetFlowRecords
func (mr *MockFlowAggregatorQuerierMockRecorder) GetFlowRecords(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFlowRecords", reflect.TypeOf((*MockFlowAggregatorQuerier)(nil).GetFlowRecords), arg0)
}

// GetRecordMetrics mocks base method
func (m *MockFlowAggregatorQuerier) GetRecordMetrics() querier.Metrics {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRecordMetrics")
	ret0, _ := ret[0].(querier.Metrics)
	return ret0
}

// GetRecordMetrics indicates an expected call of GetRecordMetrics
func (mr *MockFlowAggregatorQuerierMockRecorder) GetRecordMetrics() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRecordMetrics", reflect.TypeOf((*MockFlowAggregatorQuerier)(nil).GetRecordMetrics))
}
