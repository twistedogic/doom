package model

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/fatih/structs"
	"github.com/iancoleman/strcase"
	"github.com/mitchellh/mapstructure"
	"github.com/twistedogic/doom/pkg/helper"
	"github.com/twistedogic/doom/pkg/jsonpath"
)

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

type Values struct {
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

func IsNumber(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func (v *Values) Parse(i interface{}) error {
	m := make(map[string]struct {
		Name  string
		Value map[string]string
	})
	if err := mapstructure.WeakDecode(i, &m); err != nil {
		return err
	}
	for _, name := range structs.Names(v) {
		snake := strcase.ToSnake(name)
		start := strings.Replace(strings.Title(snake), "_", " ", -1)
		title := strings.Title(strings.Replace(snake, "_", " ", -1))
		for key, val := range m {
			if (val.Name == title || val.Name == start) && (IsNumber(key) || key == "goalattempts") {
				var value Value
				if err := ParseValue(val.Value, &value); err != nil {
					return err
				}
				if err := helper.SetField(v, name, value); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

type Detail struct {
	ID     uint64 `jsonpath:"$._matchid"`
	Home   string `jsonpath:"$.teams.home" boltholdIndex:"Team"`
	Away   string `jsonpath:"$.teams.away" boltholdIndex:"Team"`
	Detail Values
}

func (d *Detail) UnmarshalJSON(b []byte) error {
	var in interface{}
	if err := json.Unmarshal(b, &in); err != nil {
		return err
	}
	if err := jsonpath.ParseJsonpath(in, d); err != nil {
		return err
	}
	value, err := jsonpath.Lookup("$.values", in)
	if err != nil {
		return err
	}
	return d.Detail.Parse(value)
}

type Team struct {
	Name string `jsonpath:"$.name" boltholdIndex:"Team"`
	ID   int    `jsonpath:"$._id"`
}

type Match struct {
	ID   int  `jsonpath:"$._id"`
	Home Team `jsonpath:"$.teams.home"`
	Away Team `jsonpath:"$.teams.away"`
	Date int  `jsonpath:"$._dt.uts"`
}
