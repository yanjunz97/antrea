// Copyright 2020 Antrea Authors
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

package updateparam

import (
	"net"
	"net/http"
	"strconv"
	"time"

	"antrea.io/antrea/pkg/flowaggregator/querier"
	"antrea.io/antrea/pkg/util/flowexport"
)

const (
	defaultExternalFlowCollectorTransport = "tcp"
	defaultExternalFlowCollectorPort      = "4739"
)

// HandleFunc returns the function which can handle the /updateparam API request.
func HandleFunc(faq querier.FlowAggregatorQuerier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logTickerStr := r.URL.Query().Get("logticker")
		if logTickerStr != "" {
			logTickerDuration, err := time.ParseDuration(logTickerStr)
			if err != nil {
				http.Error(w, "Error when parsing logticker: "+err.Error(), http.StatusNotFound)
			}
			faq.SetLogTicker(logTickerDuration)
		}
		includePodLabelsStr := r.URL.Query().Get("podlabels")
		if includePodLabelsStr != "" {
			includePodLabels, err := strconv.ParseBool(includePodLabelsStr)
			if err != nil {
				http.Error(w, "Error when parsing podlabels: "+err.Error(), http.StatusNotFound)
			}
			faq.SetIncludePodLabels(includePodLabels)
		}
		externalFlowCollectorAddr := r.URL.Query().Get("externalflowcollectoraddr")
		if externalFlowCollectorAddr != "" {
			host, port, proto, err := flowexport.ParseFlowCollectorAddr(externalFlowCollectorAddr, defaultExternalFlowCollectorPort, defaultExternalFlowCollectorTransport)
			if err != nil {
				http.Error(w, "Error when parsing externalFlowCollectorAddr: "+err.Error(), http.StatusNotFound)
			}
			faq.SetExternalFlowCollectorAddr(querier.ExternalFlowCollectorAddr{
				Address:  net.JoinHostPort(host, port),
				Protocol: proto,
			})
		}
		activeFlowRecordTimeoutStr := r.URL.Query().Get("activeflowrecordtimeout")
		if activeFlowRecordTimeoutStr != "" {
			activeFlowRecordTimeout, err := time.ParseDuration(activeFlowRecordTimeoutStr)
			if err != nil {
				http.Error(w, "Error when parsing activeflowrecordtimeout: "+err.Error(), http.StatusNotFound)
			}
			faq.SetActiveFlowRecordTimeout(activeFlowRecordTimeout)
		}
		inactiveFlowRecordTimeoutStr := r.URL.Query().Get("inactiveflowrecordtimeout")
		if inactiveFlowRecordTimeoutStr != "" {
			inactiveFlowRecordTimeout, err := time.ParseDuration(inactiveFlowRecordTimeoutStr)
			if err != nil {
				http.Error(w, "Error when parsing inactiveflowrecordtimeout: "+err.Error(), http.StatusNotFound)
			}
			faq.SetInactiveFlowRecordTimeout(inactiveFlowRecordTimeout)
		}
	}
}
