package model

import (
	"strconv"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/iancoleman/strcase"
	json "github.com/json-iterator/go"
	"github.com/mitchellh/mapstructure"
	"github.com/twistedogic/doom/pkg/helper/tag"
	"github.com/twistedogic/jsonpath"
)

type Team struct {
	Name string `jsonpath:"$.name"`
	ID   int    `jsonpath:"$._id"`
}

type Score struct {
	Home int `jsonpath:"$.home"`
	Away int `jsonpath:"$.away"`
}

type Match struct {
	ID       int   `jsonpath:"$._id"`
	Home     Team  `jsonpath:"$.teams.home"`
	Away     Team  `jsonpath:"$.teams.away"`
	Date     int   `jsonpath:"$._dt.uts"`
	Result   Score `jsonpath:"$.result"`
	FullTime Score `jsonpath:"$.periods.ft,omitempty"`
	HalfTime Score `jsonpath:"$.periods.p1,omitempty"`
	OverTime Score `jsonpath:"$.periods.ot,omitempty"`
	PK       Score `jsonpath:"$.periods.ap,omitempty"`
}

func (m Match) IsFinish() bool {
	return time.Now().After(time.Unix(int64(m.Date), 0))
}

func IsNumber(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

type Value struct {
	Home int
	Away int
}

func ParseValue(i map[string]string, v *Value) error {
	o := make(map[string]int)
	for k, s := range i {
		o[k] = 0
		for _, sval := range strings.Split(s, "/") {
			val, err := strconv.Atoi(sval)
			if err != nil {
				return err
			}
			o[k] += val
		}
	}
	return mapstructure.WeakDecode(o, v)
}

type Detail struct {
	BallPossession  Value
	GoalKicks       Value
	FreeKicks       Value
	Offsides        Value
	CornerKicks     Value
	ShotsOnTarget   Value
	ShotsOffTarget  Value
	ShotsBlocked    Value
	Saves           Value
	Fouls           Value
	Injuries        Value
	BallSafe        Value
	Substitutions   Value
	DangerousAttack Value
	Attack          Value
	GoalAttempts    Value
}

func (d *Detail) Parse(i interface{}) error {
	m := make(map[string]struct {
		Name  string
		Value map[string]string
	})
	if err := mapstructure.WeakDecode(i, &m); err != nil {
		return err
	}
	for _, name := range structs.Names(d) {
		snake := strcase.ToSnake(name)
		start := strings.Replace(strings.Title(snake), "_", " ", -1)
		title := strings.Title(strings.Replace(snake, "_", " ", -1))
		for key, val := range m {
			if (val.Name == title || val.Name == start) && (IsNumber(key) || key == "goalattempts") {
				var value Value
				if err := ParseValue(val.Value, &value); err != nil {
					return err
				}
				if err := tag.SetField(d, name, value); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (d *Detail) UnmarshalJSONPath(b []byte) error {
	var in interface{}
	if err := json.Unmarshal(b, &in); err != nil {
		return err
	}
	value, err := jsonpath.Lookup("$.values", in)
	if err != nil {
		return nil
	}
	return d.Parse(value)
}

type MatchDetail struct {
	Match
	Detail
}
