/* Copyright 2020 Noah Hummel
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

package kafka_golang

import "encoding/json"

// JSONProducer is capable of producing messages whose body is JSON encoded.
type JSONProducer interface {
	// ProduceJSONSync produces a message without headers and key by JSON-encoding the given value.
	ProduceJSONSync(topic string, partition int32, value interface{}) error
}

func (k *KafkaProducer) ProduceJSONSync(topic string, partition int32, value interface{}) error {
	marshaled, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return k.ProduceSimpleSync(topic, partition, marshaled)
}
