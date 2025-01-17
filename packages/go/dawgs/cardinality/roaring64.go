// Copyright 2023 Specter Ops, Inc.
// 
// Licensed under the Apache License, Version 2.0
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
// 
// SPDX-License-Identifier: Apache-2.0

package cardinality

import (
	"github.com/RoaringBitmap/roaring/roaring64"
)

type bitmap64Iterator struct {
	iterator roaring64.IntPeekable64
}

func (s bitmap64Iterator) HasNext() bool {
	return s.iterator.HasNext()
}

func (s bitmap64Iterator) Next() uint64 {
	return s.iterator.Next()
}

type bitmap64 struct {
	bitmap *roaring64.Bitmap
}

func NewBitmap64() Duplex[uint64] {
	return bitmap64{
		bitmap: roaring64.New(),
	}
}

func (s bitmap64) Each(delegate func(nextValue uint64) (bool, error)) error {
	for itr := s.bitmap.Iterator(); itr.HasNext(); {
		if ok, err := delegate(itr.Next()); err != nil {
			return err
		} else if !ok {
			break
		}
	}

	return nil
}

func (s bitmap64) Iterator() Iterator[uint64] {
	return bitmap64Iterator{
		iterator: s.bitmap.Iterator(),
	}
}

func (s bitmap64) Slice() []uint64 {
	return s.bitmap.ToArray()
}

func (s bitmap64) Contains(value uint64) bool {
	return s.bitmap.Contains(value)
}

func (s bitmap64) CheckedAdd(value uint64) bool {
	return s.bitmap.CheckedAdd(value)
}

func (s bitmap64) Add(values ...uint64) {
	s.bitmap.AddMany(values)
}

func (s bitmap64) Remove(value uint64) {
	s.bitmap.Remove(value)
}

func (s bitmap64) And(provider Provider[uint64]) {
	switch typedProvider := provider.(type) {
	case bitmap64:
		s.bitmap.Or(typedProvider.bitmap)

	case Duplex[uint64]:
		s.Each(func(nextValue uint64) (bool, error) {
			if !typedProvider.Contains(nextValue) {
				s.Remove(nextValue)
			}

			return true, nil
		})
	}
}
func (s bitmap64) Or(provider Provider[uint64]) {
	switch typedProvider := provider.(type) {
	case bitmap64:
		s.bitmap.Or(typedProvider.bitmap)

	case Duplex[uint64]:
		typedProvider.Each(func(nextValue uint64) (bool, error) {
			s.Add(nextValue)
			return true, nil
		})
	}
}

func (s bitmap64) Cardinality() uint64 {
	return s.bitmap.GetCardinality()
}

func (s bitmap64) Clone() Duplex[uint64] {
	return bitmap64{
		bitmap: s.bitmap.Clone(),
	}
}
