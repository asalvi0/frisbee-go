/*
	Copyright 2022 Loophole Labs

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

package packet

import (
	"github.com/loopholelabs/frisbee/internal/protocol"
	"sync"
)

// Packet is the structured frisbee data packet, and contains the following:
//
//	type Packet struct {
//		Message struct {
//			Id            uint16 // 2 Bytes
//			Operation     uint16 // 2 Bytes
//			ContentLength uint32 // 4 Bytes
//		}
//		Content []byte
//	}
//
// The ID field can be used however the user sees fit, however ContentLength must match the length of the content being
// delivered with the frisbee message (see the Async.WriteMessage function for more details), and the Operation field must be greater than uint16(9).
type Packet struct {
	Message *protocol.Message
	Content []byte
}

// Write efficiently copies the byte slice b into the packet, however it
// does *not* update the content length.
func (p *Packet) Write(b []byte) {
	if len(p.Content) < len(b) {
		p.Content = append(p.Content[0:], b...)
	} else {
		p.Content = p.Content[:copy(p.Content[0:], b)]
	}
}

func (p *Packet) Reset() {
	p.Message.Id = 0
	p.Message.Operation = 0
	p.Message.ContentLength = 0
	p.Content = p.Content[:0]
}

var pool sync.Pool

func Get() (s *Packet) {
	v := pool.Get()
	if v == nil {
		v = new(Packet)
	}

	s = v.(*Packet)
	if s.Message == nil {
		s.Message = new(protocol.Message)
	}
	return
}

func Put(p *Packet) {
	if p != nil {
		p.Reset()
		pool.Put(p)
	}
}
