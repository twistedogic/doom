// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    match, err := UnmarshalMatch(bytes)
//    bytes, err = match.Marshal()

package match

import "bytes"
import "errors"
import "encoding/json"

func UnmarshalMatch(data []byte) (Match, error) {
	var r Match
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Match) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Match struct {
	Doc      []Doc  `json:"doc"`
	QueryURL string `json:"queryUrl"`
}

type Doc struct {
	Dob    int64   `json:"_dob"`
	Maxage int64   `json:"_maxage"`
	Data   []Datum `json:"data"`
	Event  string  `json:"event"`
}

type Datum struct {
	Doc            string         `json:"_doc"`
	ID             int64          `json:"_id"`
	Sid            int64          `json:"_sid"`
	Sk             int64          `json:"_sk"`
	Live           bool           `json:"live"`
	Name           string         `json:"name"`
	Realcategories []Realcategory `json:"realcategories"`
}

type Realcategory struct {
	Doc         string       `json:"_doc"`
	ID          int64        `json:"_id"`
	Rcid        int64        `json:"_rcid"`
	Sid         int64        `json:"_sid"`
	Sk          bool         `json:"_sk"`
	Cc          *Cc          `json:"cc"`
	Name        string       `json:"name"`
	Tournaments []Tournament `json:"tournaments"`
}

type Cc struct {
	Doc         string `json:"_doc"`
	ID          int64  `json:"_id"`
	A2          string `json:"a2"`
	A3          string `json:"a3"`
	Continent   string `json:"continent"`
	Continentid int64  `json:"continentid"`
	Ioc         string `json:"ioc"`
	Name        string `json:"name"`
	Population  int64  `json:"population"`
}

type Tournament struct {
	Doc                  string         `json:"_doc"`
	ID                   int64          `json:"_id"`
	Isk                  int64          `json:"_isk"`
	Rcid                 int64          `json:"_rcid"`
	Sid                  int64          `json:"_sid"`
	Sk                   bool           `json:"_sk"`
	Tid                  int64          `json:"_tid"`
	Utid                 int64          `json:"_utid"`
	Abbr                 string         `json:"abbr"`
	Cuprosterid          *string        `json:"cuprosterid"`
	Currentseason        int64          `json:"currentseason"`
	Friendly             bool           `json:"friendly"`
	Ground               interface{}    `json:"ground"`
	Livetable            *Livetable     `json:"livetable"`
	Matches              []MatchElement `json:"matches"`
	Name                 string         `json:"name"`
	Outdated             bool           `json:"outdated"`
	Roundbyround         bool           `json:"roundbyround"`
	Seasonid             int64          `json:"seasonid"`
	Seasontype           string         `json:"seasontype"`
	Seasontypename       string         `json:"seasontypename"`
	Seasontypeunique     *string        `json:"seasontypeunique"`
	Tournamentlevelname  *string        `json:"tournamentlevelname"`
	Tournamentlevelorder *int64         `json:"tournamentlevelorder"`
	Year                 string         `json:"year"`
	Groupname            *string        `json:"groupname,omitempty"`
}

type MatchElement struct {
	Doc                     string        `json:"_doc"`
	Dt                      Dt            `json:"_dt"`
	ID                      int64         `json:"_id"`
	Mclink                  *bool         `json:"_mclink,omitempty"`
	Rcid                    int64         `json:"_rcid"`
	Seasonid                int64         `json:"_seasonid"`
	Sid                     int64         `json:"_sid"`
	Sk                      bool          `json:"_sk"`
	Tid                     int64         `json:"_tid"`
	Utid                    int64         `json:"_utid"`
	Cancelled               bool          `json:"cancelled"`
	Coverage                Coverage      `json:"coverage"`
	Distance                *int64        `json:"distance,omitempty"`
	Facts                   bool          `json:"facts"`
	HF                      float64       `json:"hf"`
	Localderby              bool          `json:"localderby"`
	Matchstatus             *Matchstatus  `json:"matchstatus"`
	Numberofperiods         int64         `json:"numberofperiods"`
	Overtimelength          int64         `json:"overtimelength"`
	P                       string        `json:"p"`
	Periodlength            int64         `json:"periodlength"`
	Periods                 *PeriodsUnion `json:"periods"`
	Postponed               bool          `json:"postponed"`
	Ptime                   *Livetable    `json:"ptime"`
	Removed                 bool          `json:"removed"`
	Result                  Result        `json:"result"`
	Round                   *int64        `json:"round,omitempty"`
	Roundname               *Roundname    `json:"roundname,omitempty"`
	Stadiumid               int64         `json:"stadiumid"`
	Status                  Status        `json:"status"`
	Teams                   Teams         `json:"teams"`
	Timeinfo                *Timeinfo     `json:"timeinfo"`
	Tobeannounced           bool          `json:"tobeannounced"`
	UpdatedUts              int64         `json:"updated_uts"`
	Walkover                bool          `json:"walkover"`
	Week                    int64         `json:"week"`
	Cards                   *Cards        `json:"cards,omitempty"`
	Pitchcondition          *int64        `json:"pitchcondition,omitempty"`
	Temperature             interface{}   `json:"temperature"`
	Weather                 *int64        `json:"weather,omitempty"`
	Wind                    interface{}   `json:"wind"`
	Windadvantage           *int64        `json:"windadvantage,omitempty"`
	Cuproundmatchnumber     *int64        `json:"cuproundmatchnumber,omitempty"`
	Cuproundnumberofmatches *int64        `json:"cuproundnumberofmatches,omitempty"`
}

type Cards struct {
	Away CardsAway `json:"away"`
	Home CardsAway `json:"home"`
}

type CardsAway struct {
	RedCount    int64 `json:"red_count"`
	YellowCount int64 `json:"yellow_count"`
}

type Coverage struct {
	Advantage           interface{} `json:"advantage"`
	Ballspotting        bool        `json:"ballspotting"`
	Basiclineup         bool        `json:"basiclineup"`
	Cornersonly         bool        `json:"cornersonly"`
	Deepercoverage      bool        `json:"deepercoverage"`
	Formations          int64       `json:"formations"`
	Hasstats            bool        `json:"hasstats"`
	Injuries            int64       `json:"injuries"`
	Inlivescore         bool        `json:"inlivescore"`
	Lineup              int64       `json:"lineup"`
	Liveodds            bool        `json:"liveodds"`
	Livetable           int64       `json:"livetable"`
	Lmtsupport          int64       `json:"lmtsupport"`
	Matchdatacomplete   bool        `json:"matchdatacomplete"`
	Mediacoverage       bool        `json:"mediacoverage"`
	Multicast           bool        `json:"multicast"`
	Penaltyshootout     int64       `json:"penaltyshootout"`
	Scoutconnected      bool        `json:"scoutconnected"`
	Scoutcoveragestatus int64       `json:"scoutcoveragestatus"`
	Scoutmatch          int64       `json:"scoutmatch"`
	Scouttest           bool        `json:"scouttest"`
	Substitutions       bool        `json:"substitutions"`
	Tacticallineup      bool        `json:"tacticallineup"`
	Tiebreak            interface{} `json:"tiebreak"`
	Venue               bool        `json:"venue"`
}

type Dt struct {
	Doc      string `json:"_doc"`
	Date     string `json:"date"`
	Time     string `json:"time"`
	Tz       string `json:"tz"`
	Tzoffset int64  `json:"tzoffset"`
	Uts      int64  `json:"uts"`
}

type PeriodsClass struct {
	Ft *Ft `json:"ft,omitempty"`
	P1 Ft  `json:"p1"`
	Ot *Ft `json:"ot,omitempty"`
}

type Ft struct {
	Away int64 `json:"away"`
	Home int64 `json:"home"`
}

type Result struct {
	Away   int64  `json:"away"`
	Home   int64  `json:"home"`
	Winner string `json:"winner"`
}

type Roundname struct {
	Doc                 string      `json:"_doc"`
	ID                  int64       `json:"_id"`
	Name                *Name       `json:"name"`
	Cuproundnumber      interface{} `json:"cuproundnumber"`
	Displaynumber       interface{} `json:"displaynumber"`
	Shortname           *string     `json:"shortname,omitempty"`
	Statisticssortorder *int64      `json:"statisticssortorder,omitempty"`
}

type Status struct {
	Doc  string `json:"_doc"`
	ID   int64  `json:"_id"`
	Name string `json:"name"`
}

type Teams struct {
	Away Team `json:"away"`
	Home Team `json:"home"`
}

type Team struct {
	Doc        string  `json:"_doc"`
	ID         int64   `json:"_id"`
	Abbr       string  `json:"abbr"`
	Haslogo    bool    `json:"haslogo"`
	Iscountry  bool    `json:"iscountry"`
	Mediumname string  `json:"mediumname"`
	Name       string  `json:"name"`
	Nickname   *string `json:"nickname"`
	Uid        int64   `json:"uid"`
	Virtual    bool    `json:"virtual"`
	Cc         *Cc     `json:"cc,omitempty"`
}

type Timeinfo struct {
	Ended      *string     `json:"ended"`
	Injurytime interface{} `json:"injurytime"`
	Played     interface{} `json:"played"`
	Remaining  interface{} `json:"remaining"`
	Running    bool        `json:"running"`
	Started    interface{} `json:"started"`
}

type Livetable struct {
	Bool    *bool
	Integer *int64
}

func (x *Livetable) UnmarshalJSON(data []byte) error {
	object, err := unmarshalUnion(data, &x.Integer, nil, &x.Bool, nil, false, nil, false, nil, false, nil, false, nil, false)
	if err != nil {
		return err
	}
	if object {
	}
	return nil
}

func (x *Livetable) MarshalJSON() ([]byte, error) {
	return marshalUnion(x.Integer, nil, x.Bool, nil, false, nil, false, nil, false, nil, false, nil, false)
}

type Matchstatus struct {
	Bool   *bool
	String *string
}

func (x *Matchstatus) UnmarshalJSON(data []byte) error {
	object, err := unmarshalUnion(data, nil, nil, &x.Bool, &x.String, false, nil, false, nil, false, nil, false, nil, false)
	if err != nil {
		return err
	}
	if object {
	}
	return nil
}

func (x *Matchstatus) MarshalJSON() ([]byte, error) {
	return marshalUnion(nil, nil, x.Bool, x.String, false, nil, false, nil, false, nil, false, nil, false)
}

type PeriodsUnion struct {
	AnythingArray []interface{}
	PeriodsClass  *PeriodsClass
}

func (x *PeriodsUnion) UnmarshalJSON(data []byte) error {
	x.AnythingArray = nil
	x.PeriodsClass = nil
	var c PeriodsClass
	object, err := unmarshalUnion(data, nil, nil, nil, nil, true, &x.AnythingArray, true, &c, false, nil, false, nil, false)
	if err != nil {
		return err
	}
	if object {
		x.PeriodsClass = &c
	}
	return nil
}

func (x *PeriodsUnion) MarshalJSON() ([]byte, error) {
	return marshalUnion(nil, nil, nil, nil, x.AnythingArray != nil, x.AnythingArray, x.PeriodsClass != nil, x.PeriodsClass, false, nil, false, nil, false)
}

type Name struct {
	Integer *int64
	String  *string
}

func (x *Name) UnmarshalJSON(data []byte) error {
	object, err := unmarshalUnion(data, &x.Integer, nil, nil, &x.String, false, nil, false, nil, false, nil, false, nil, false)
	if err != nil {
		return err
	}
	if object {
	}
	return nil
}

func (x *Name) MarshalJSON() ([]byte, error) {
	return marshalUnion(x.Integer, nil, nil, x.String, false, nil, false, nil, false, nil, false, nil, false)
}

func unmarshalUnion(data []byte, pi **int64, pf **float64, pb **bool, ps **string, haveArray bool, pa interface{}, haveObject bool, pc interface{}, haveMap bool, pm interface{}, haveEnum bool, pe interface{}, nullable bool) (bool, error) {
	if pi != nil {
		*pi = nil
	}
	if pf != nil {
		*pf = nil
	}
	if pb != nil {
		*pb = nil
	}
	if ps != nil {
		*ps = nil
	}

	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	tok, err := dec.Token()
	if err != nil {
		return false, err
	}

	switch v := tok.(type) {
	case json.Number:
		if pi != nil {
			i, err := v.Int64()
			if err == nil {
				*pi = &i
				return false, nil
			}
		}
		if pf != nil {
			f, err := v.Float64()
			if err == nil {
				*pf = &f
				return false, nil
			}
			return false, errors.New("Unparsable number")
		}
		return false, errors.New("Union does not contain number")
	case float64:
		return false, errors.New("Decoder should not return float64")
	case bool:
		if pb != nil {
			*pb = &v
			return false, nil
		}
		return false, errors.New("Union does not contain bool")
	case string:
		if haveEnum {
			return false, json.Unmarshal(data, pe)
		}
		if ps != nil {
			*ps = &v
			return false, nil
		}
		return false, errors.New("Union does not contain string")
	case nil:
		if nullable {
			return false, nil
		}
		return false, errors.New("Union does not contain null")
	case json.Delim:
		if v == '{' {
			if haveObject {
				return true, json.Unmarshal(data, pc)
			}
			if haveMap {
				return false, json.Unmarshal(data, pm)
			}
			return false, errors.New("Union does not contain object")
		}
		if v == '[' {
			if haveArray {
				return false, json.Unmarshal(data, pa)
			}
			return false, errors.New("Union does not contain array")
		}
		return false, errors.New("Cannot handle delimiter")
	}
	return false, errors.New("Cannot unmarshal union")

}

func marshalUnion(pi *int64, pf *float64, pb *bool, ps *string, haveArray bool, pa interface{}, haveObject bool, pc interface{}, haveMap bool, pm interface{}, haveEnum bool, pe interface{}, nullable bool) ([]byte, error) {
	if pi != nil {
		return json.Marshal(*pi)
	}
	if pf != nil {
		return json.Marshal(*pf)
	}
	if pb != nil {
		return json.Marshal(*pb)
	}
	if ps != nil {
		return json.Marshal(*ps)
	}
	if haveArray {
		return json.Marshal(pa)
	}
	if haveObject {
		return json.Marshal(pc)
	}
	if haveMap {
		return json.Marshal(pm)
	}
	if haveEnum {
		return json.Marshal(pe)
	}
	if nullable {
		return json.Marshal(nil)
	}
	return nil, errors.New("Union must not be null")
}
