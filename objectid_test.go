package objectid_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"go.neonxp.dev/objectid"
)

func TestID_String(t *testing.T) {
	objectid.Seed()
	id1 := objectid.New()
	id1String := id1.String()
	id2, err := objectid.FromString(id1String)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(id1, id2) {
		t.Errorf("ID2 = %v, want %v", id2, id1)
	}
	wantErr := "is invalid object id (length=3)"
	if _, err := objectid.FromString("test"); !strings.HasSuffix(err.Error(), wantErr) {
		t.Errorf("Err = %v, want suffix %s", err, wantErr)
	}
}

func TestID_JSON(t *testing.T) {
	objectid.Seed()
	type testStruct struct {
		ID objectid.ID `json:"oid"`
	}

	t1 := testStruct{
		ID: objectid.New(),
	}

	b1, err := json.Marshal(t1)
	if err != nil {
		t.Error(err)
	}

	t2 := new(testStruct)

	if err := json.Unmarshal(b1, t2); err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(t1, *t2) {
		t.Errorf("T2 = %v, want %v", t2, t1)
	}

	id := objectid.New()
	wantErr := fmt.Errorf(" is invalid object id (length=0)")
	if err := id.UnmarshalJSON(nil); err.Error() != wantErr.Error() {
		t.Errorf("UnmarshalJSON(nil) = %v, want %v", err, wantErr)
	}
}

func TestID_Time(t *testing.T) {
	objectid.Seed()
	testTime, _ := time.Parse(time.RFC3339, "2020-12-20T10:11:50Z")
	id := objectid.FromTime(testTime.UTC())
	if got := id.Time(); !reflect.DeepEqual(got.UTC(), testTime) {
		t.Errorf("ID.Time() = %v, want %v", got.String(), testTime.String())
	}
}

func TestID_Counter(t *testing.T) {
	objectid.Seed()
	for i := 1; i <= 10; i++ {
		id := objectid.New()
		if id.Counter() != uint64(i) {
			t.Errorf("id.Counter = %d, want %d", i, id.Counter())
		}
	}
}

func TestID_Less(t *testing.T) {
	objectid.Seed()
	type args struct {
		i2 objectid.ID
	}
	nowTime := time.Now()
	tests := []struct {
		name string
		i    objectid.ID
		args args
		want bool
	}{
		{
			name: "by time",
			i:    objectid.New(),
			args: args{
				i2: objectid.New(),
			},
			want: true,
		},
		{
			name: "by counter",
			i:    objectid.FromTime(nowTime),
			args: args{
				i2: objectid.FromTime(nowTime),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.Less(tt.args.i2); got != tt.want {
				t.Errorf("ID.Less() = %v, want %v", got, tt.want)
			}
		})
	}
}
