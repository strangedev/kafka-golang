/* Copyright 2020 Noah Hummel
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

package maps

import (
	"github.com/google/uuid"
	"github.com/strangedev/kafka-golang/lib"
	"github.com/strangedev/kafka-golang/schema"
)

type aliasMapType map[schema.Alias]uuid.UUID

// SchemaMap is a KeyObservable map of Aliases to UUIDs
type AliasMap struct {
	lib.ConcurrentObservable
	Map aliasMapType
}

// Upsert inserts or updates an Alias, UUID pair into the map.
// All of the map's observers are notified of this change.
// It returns true, if the map entry did already exist and was overwritten.
func (m AliasMap) Insert(alias schema.Alias, schemaUUID uuid.UUID) bool {
	m.DataLock.Lock()
	_, overwritten := m.Map[alias]
	m.Map[alias] = schemaUUID
	m.DataLock.Unlock()
	m.Notify(alias)
	return overwritten
}

// NewAliasMap constructs an empty AliasMap with no observers.
func NewAliasMap() AliasMap {
	return AliasMap{
		ConcurrentObservable: lib.NewConcurrentObservable(),
		Map:                  make(aliasMapType),
	}
}
