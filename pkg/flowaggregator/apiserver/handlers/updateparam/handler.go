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
	"net/http"
	"strconv"
	"time"

	"antrea.io/antrea/pkg/flowaggregator/querier"
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
			faq.UpdateLogTicker(logTickerDuration)
		}
		includePodLabelsStr := r.URL.Query().Get("podlabels")
		if includePodLabelsStr != "" {
			includePodLabels, err := strconv.ParseBool(includePodLabelsStr)
			if err != nil {
				http.Error(w, "Error when parsing podlabels: "+err.Error(), http.StatusNotFound)
			}
			faq.UpdateIncludePodLabels(includePodLabels)
		}
	}
}
