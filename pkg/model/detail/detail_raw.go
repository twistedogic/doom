// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    detail, err := UnmarshalDetail(bytes)
//    bytes, err = detail.Marshal()

package detail

import "bytes"
import "errors"
import "encoding/json"

func UnmarshalDetail(data []byte) (Detail, error) {
	var r Detail
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Detail) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Detail struct {
	Doc      []Doc  `json:"doc"`
	QueryURL string `json:"queryUrl"`
}

type Doc struct {
	Dob    int64  `json:"_dob"`
	Maxage int64  `json:"_maxage"`
	Data   Data   `json:"data"`
	Event  string `json:"event"`
}

type Data struct {
	Doc     string  `json:"_doc"`
	Matchid int64   `json:"_matchid"`
	Index   []Index `json:"index"`
	Teams   Teams   `json:"teams"`
	Types   Types   `json:"types"`
	Values  Values  `json:"values"`
}

type Teams struct {
	Away string `json:"away"`
	Home string `json:"home"`
}

type Types struct {
	The40                     string `json:"40"`
	The45                     string `json:"45"`
	The50                     string `json:"50"`
	The60                     string `json:"60"`
	The110                    string `json:"110"`
	The120                    string `json:"120"`
	The121                    string `json:"121"`
	The122                    string `json:"122"`
	The123                    string `json:"123"`
	The124                    string `json:"124"`
	The125                    string `json:"125"`
	The126                    string `json:"126"`
	The127                    string `json:"127"`
	The129                    string `json:"129"`
	The158                    string `json:"158"`
	The161                    string `json:"161"`
	The171                    string `json:"171"`
	The1029                   string `json:"1029"`
	The1030                   string `json:"1030"`
	The1126                   string `json:"1126"`
	Attackpercentage          string `json:"attackpercentage"`
	Ballsafepercentage        string `json:"ballsafepercentage"`
	Dangerousattackpercentage string `json:"dangerousattackpercentage"`
	Goalattempts              string `json:"goalattempts"`
}

type Values struct {
	The40                     TartuGecko `json:"40"`
	The60                     TartuGecko `json:"60"`
	The110                    TartuGecko `json:"110"`
	The120                    TartuGecko `json:"120"`
	The121                    TartuGecko `json:"121"`
	The122                    TartuGecko `json:"122"`
	The123                    TartuGecko `json:"123"`
	The124                    TartuGecko `json:"124"`
	The125                    TartuGecko `json:"125"`
	The126                    TartuGecko `json:"126"`
	The127                    TartuGecko `json:"127"`
	The129                    TartuGecko `json:"129"`
	The158                    TartuGecko `json:"158"`
	The161                    The161     `json:"161"`
	The171                    TartuGecko `json:"171"`
	The1029                   TartuGecko `json:"1029"`
	The1030                   TartuGecko `json:"1030"`
	The1126                   TartuGecko `json:"1126"`
	Attackpercentage          TartuGecko `json:"attackpercentage"`
	Ballsafepercentage        TartuGecko `json:"ballsafepercentage"`
	Dangerousattackpercentage TartuGecko `json:"dangerousattackpercentage"`
	Goalattempts              TartuGecko `json:"goalattempts"`
}

type TartuGecko struct {
	Name  string `json:"name"`
	Value Value  `json:"value"`
}

type Value struct {
	Away int64 `json:"away"`
	Home int64 `json:"home"`
}

type The161 struct {
	Name  string `json:"name"`
	Value Teams  `json:"value"`
}

type Index struct {
	Integer *int64
	String  *string
}

func (x *Index) UnmarshalJSON(data []byte) error {
	object, err := unmarshalUnion(data, &x.Integer, nil, nil, &x.String, false, nil, false, nil, false, nil, false, nil, false)
	if err != nil {
		return err
	}
	if object {
	}
	return nil
}

func (x *Index) MarshalJSON() ([]byte, error) {
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
