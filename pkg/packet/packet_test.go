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
	"crypto/rand"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	t.Parallel()

	p := Get()

	assert.IsType(t, new(Packet), p)
	assert.NotNil(t, p.Metadata)
	assert.Equal(t, uint16(0), p.Metadata.Id)
	assert.Equal(t, uint16(0), p.Metadata.Operation)
	assert.Equal(t, uint32(0), p.Metadata.ContentLength)
	assert.Equal(t, []byte{}, p.Content.B)

	Put(p)
}

func TestRecycle(t *testing.T) {
	pool := NewPool()

	p := pool.Get()

	p.Metadata.Id = 32
	p.Metadata.Operation = 64
	p.Metadata.ContentLength = 128

	pool.Put(p)
	p = pool.Get()

	testData := make([]byte, p.Content.Cap()*2)
	_, err := rand.Read(testData)
	assert.NoError(t, err)
	for {
		assert.NotNil(t, p.Metadata)
		assert.Equal(t, uint16(0), p.Metadata.Id)
		assert.Equal(t, uint16(0), p.Metadata.Operation)
		assert.Equal(t, uint32(0), p.Metadata.ContentLength)
		assert.Equal(t, []byte{}, p.Content.B)

		p.Content.Write(testData)
		assert.Equal(t, len(testData), p.Content.Len())
		assert.GreaterOrEqual(t, p.Content.Cap(), len(testData))

		pool.Put(p)
		p = pool.Get()

		assert.NotNil(t, p.Metadata)
		assert.Equal(t, uint16(0), p.Metadata.Id)
		assert.Equal(t, uint16(0), p.Metadata.Operation)
		assert.Equal(t, uint32(0), p.Metadata.ContentLength)

		if p.Content.Cap() < len(testData) {
			continue
		} else {
			assert.Equal(t, 0, p.Content.Len())
			assert.GreaterOrEqual(t, p.Content.Cap(), len(testData))
			break
		}
	}

	pool.Put(p)
}

func TestWrite(t *testing.T) {
	t.Parallel()

	p := Get()

	b := make([]byte, 32)
	_, err := rand.Read(b)
	assert.NoError(t, err)

	p.Content.Write(b)
	assert.Equal(t, b, p.Content.B)

	p.Reset()
	assert.NotEqual(t, b, p.Content.B)
	assert.Equal(t, 0, p.Content.Len())
	assert.Equal(t, 512, p.Content.Cap())

	b = make([]byte, 1024)
	_, err = rand.Read(b)
	assert.NoError(t, err)

	p.Content.Write(b)

	assert.Equal(t, b, p.Content.B)
	assert.Equal(t, 1024, p.Content.Len())
	assert.GreaterOrEqual(t, p.Content.Cap(), 1024)

}
