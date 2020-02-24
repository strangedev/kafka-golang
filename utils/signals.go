/* Copyright 2020 Noah Hummel
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

package utils

import "os"

// SigAbort waits for a value on the ready channel.
// Waiting is be aborted, when a value is received on the abort channel first.
// It returns a channel that receives a value when either condition has been met,
// true if ready was first to receive, true otherwise.
func SigAbort(ready chan bool, abort chan os.Signal) chan bool {
	fin := make(chan bool)
	go (func() {
		select {
		case _ = <-abort:
			fin <- false
			close(fin)
			return
		case _ = <-ready:
			fin <- true
			close(fin)
			return
		}
	})()
	return fin
}
