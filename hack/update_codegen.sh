#!/usr/bin/env bash

# Copyright 2021 the Apache Pinot Kubernetes Operator authors.

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
set -o pipefail

CUSTOM_HEADER=${PWD}/hack/boilerplate.go.txt

GOBIN=${PWD}/bin "${PWD}"/hack/generate_groups.sh \
  client,lister,informer \
  github.com/spaghettifunk/pinot-operator/pkg/client \
  pinot:v1alpha1 \
  --go-header-file "${CUSTOM_HEADER}"