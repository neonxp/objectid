package objectid

/**
  ObjectID generator
  Copyright (C) 2022  Alexaner Kiryukhin <i@neonxp.dev>

  This program is free software: you can redistribute it and/or modify
  it under the terms of the GNU General Public License as published by
  the Free Software Foundation, either version 3 of the License, or
  (at your option) any later version.

  This program is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU General Public License for more details.

  You should have received a copy of the GNU General Public License
  along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"sync/atomic"
	"time"
)

var (
	iterator     uint64 = 0
	iteratorSeed uint64
)

// Seed необходимо вызвать в начале сессии
func Seed() {
	rand.Seed(time.Now().UnixMicro())
	iteratorSeed = rand.Uint64()
	iterator = 0
}

// New возвращает новый идентификатор
func New() ID {
	return FromTime(time.Now())
}

// FromTime возвращает идентификатор на основе переданного времени
func FromTime(t time.Time) ID {
	p1 := uint64(t.UnixMicro())
	p2 := atomic.AddUint64(&iterator, 1) + iteratorSeed
	p3 := rand.Uint64()
	r := make([]byte, 0, 24)
	r = binary.BigEndian.AppendUint64(r, p1)
	r = binary.BigEndian.AppendUint64(r, p2)
	r = binary.BigEndian.AppendUint64(r, p3)
	log.Println(r)
	return r
}

// FromString возвращает идентификатор из base64 представления
func FromString(s string) (ID, error) {
	b, err := base64.RawStdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	if len(b) != 24 {
		return nil, fmt.Errorf("%s is invalid object id (length=%d)", b, len(b))
	}
	return ID(b), nil
}

// ID тип представляющий собой идентификатор
type ID []byte

// String возвращает base64 представление идентификатора
func (i ID) String() string {
	return base64.RawStdEncoding.EncodeToString(i)
}

// MarshalJSON формирует json представление идентификатора
func (i ID) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, i.String())), nil
}

// UnmarshalJSON парсит идентификатор из json
func (i *ID) UnmarshalJSON(b []byte) error {
	if len(b) != 34 {
		return fmt.Errorf("%s is invalid object id (length=%d)", b, len(b))
	}
	id, err := FromString(string(b[1 : len(b)-1]))
	*i = id
	return err
}

// Time возвращает время создания идентификатора
func (i ID) Time() time.Time {
	if len(i) < 8 {
		return time.Time{}
	}
	t := i[:8]
	return time.UnixMicro(int64(binary.BigEndian.Uint64(t)))
}

// Counter возвращает порядковый номер идентификатора в сессии
func (i ID) Counter() uint64 {
	if len(i) < 16 {
		return 0
	}
	c := i[8:16]
	return binary.BigEndian.Uint64(c) - iteratorSeed
}

// Less возвращает true если i2 > i
func (i ID) Less(i2 ID) bool {
	if i.Time().Equal(i2.Time()) {
		return i.Counter() < i2.Counter()
	}
	return i.Time().Before(i2.Time())
}
