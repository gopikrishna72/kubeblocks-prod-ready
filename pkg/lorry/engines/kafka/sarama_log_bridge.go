/*
Copyright 2021 The Dapr Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kafka

import (
	"fmt"

	"github.com/go-logr/logr"
)

type SaramaLogBridge struct {
	logger logr.Logger
}

func (b SaramaLogBridge) Print(v ...interface{}) {
	b.logger.Info(fmt.Sprint(v...))
}

func (b SaramaLogBridge) Printf(format string, v ...interface{}) {
	b.logger.Info(fmt.Sprintf(format, v...))
}

func (b SaramaLogBridge) Println(v ...interface{}) {
	b.Print(v...)
}
