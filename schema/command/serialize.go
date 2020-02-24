/* Copyright 2020 Noah Hummel
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

// Package command contains commands for updating the schema stored in Kafka.
// Since the schema repository can't use itself to specify the Command/Message format
// for update actions, it uses plain JSON.
package command

import "github.com/google/uuid"

// UpdateRequest sets the given UUID to equal the given plain-text Avro spec.
type UpdateRequest struct {
	UUID uuid.UUID `json:"UUID"`
	Spec string    `json:"spec"`
}

// AliasRequest sets the given Alias to equal the given UUID.
type AliasRequest struct {
	UUID  uuid.UUID `json:"UUID"`
	Alias string    `json:"alias"`
}
