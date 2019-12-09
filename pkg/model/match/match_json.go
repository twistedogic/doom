package match

type Feed struct {
	Doc      []Doc  `json:"doc"`
	QueryURL string `json:"queryUrl"`
}

type Doc struct {
	Dob    int     `json:"_dob"`
	Maxage int     `json:"_maxage"`
	Data   []Datum `json:"data"`
	Event  string  `json:"event"`
}

type Datum struct {
	ID             int            `json:"_id"`
	Sid            int            `json:"_sid"`
	Sk             int            `json:"_sk"`
	Live           bool           `json:"live"`
	Name           string         `json:"name"`
	Realcategories []Realcategory `json:"realcategories"`
}

type Realcategory struct {
	ID          int          `json:"_id"`
	Rcid        int          `json:"_rcid"`
	Sid         int          `json:"_sid"`
	Sk          bool         `json:"_sk"`
	Name        string       `json:"name"`
	Tournaments []Tournament `json:"tournaments"`
}

type Tournament struct {
	Doc                  string      `json:"_doc"`
	ID                   int         `json:"_id"`
	Isk                  int         `json:"_isk"`
	Rcid                 int         `json:"_rcid"`
	Sid                  int         `json:"_sid"`
	Sk                   bool        `json:"_sk"`
	Tid                  int         `json:"_tid"`
	Utid                 int         `json:"_utid"`
	Abbr                 string      `json:"abbr"`
	Cuprosterid          *string     `json:"cuprosterid"`
	Currentseason        int         `json:"currentseason"`
	Friendly             bool        `json:"friendly"`
	Ground               interface{} `json:"ground"`
	Matches              []Match     `json:"matches"`
	Name                 string      `json:"name"`
	Outdated             bool        `json:"outdated"`
	Roundbyround         bool        `json:"roundbyround"`
	Seasonid             int         `json:"seasonid"`
	Seasontype           string      `json:"seasontype"`
	Seasontypename       string      `json:"seasontypename"`
	Seasontypeunique     *string     `json:"seasontypeunique"`
	Tournamentlevelname  *string     `json:"tournamentlevelname"`
	Tournamentlevelorder *int        `json:"tournamentlevelorder"`
	Year                 string      `json:"year"`
	Groupname            *string     `json:"groupname,omitempty"`
}

type Match struct {
	Dt                      Dt          `json:"_dt"`
	ID                      int         `json:"_id"`
	Mclink                  *bool       `json:"_mclink,omitempty"`
	Rcid                    int         `json:"_rcid"`
	Seasonid                int         `json:"_seasonid"`
	Sid                     int         `json:"_sid"`
	Sk                      bool        `json:"_sk"`
	Tid                     int         `json:"_tid"`
	Utid                    int         `json:"_utid"`
	Cancelled               bool        `json:"cancelled"`
	Coverage                Coverage    `json:"coverage"`
	Distance                *int        `json:"distance,omitempty"`
	Facts                   bool        `json:"facts"`
	HF                      float64     `json:"hf"`
	Localderby              bool        `json:"localderby"`
	Numberofperiods         int         `json:"numberofperiods"`
	Overtimelength          int         `json:"overtimelength"`
	P                       string      `json:"p"`
	Periodlength            int         `json:"periodlength"`
	Periods                 *Periods    `json:"periods"`
	Postponed               bool        `json:"postponed"`
	Removed                 bool        `json:"removed"`
	Result                  Result      `json:"result"`
	Round                   *int        `json:"round,omitempty"`
	Roundname               *Roundname  `json:"roundname,omitempty"`
	Stadiumid               int         `json:"stadiumid"`
	Status                  Status      `json:"status"`
	Teams                   Teams       `json:"teams"`
	Timeinfo                *Timeinfo   `json:"timeinfo"`
	Tobeannounced           bool        `json:"tobeannounced"`
	UpdatedUts              int         `json:"updated_uts"`
	Walkover                bool        `json:"walkover"`
	Week                    int         `json:"week"`
	Cards                   *Cards      `json:"cards,omitempty"`
	Pitchcondition          *int        `json:"pitchcondition,omitempty"`
	Temperature             interface{} `json:"temperature"`
	Weather                 *int        `json:"weather,omitempty"`
	Wind                    interface{} `json:"wind"`
	Windadvantage           *int        `json:"windadvantage,omitempty"`
	Cuproundmatchnumber     *int        `json:"cuproundmatchnumber,omitempty"`
	Cuproundnumberofmatches *int        `json:"cuproundnumberofmatches,omitempty"`
}

type Cards struct {
	Away CardsAway `json:"away"`
	Home CardsAway `json:"home"`
}

type CardsAway struct {
	RedCount    int `json:"red_count"`
	YellowCount int `json:"yellow_count"`
}

type Coverage struct {
	Advantage           interface{} `json:"advantage"`
	Ballspotting        bool        `json:"ballspotting"`
	Basiclineup         bool        `json:"basiclineup"`
	Cornersonly         bool        `json:"cornersonly"`
	Deepercoverage      bool        `json:"deepercoverage"`
	Formations          int         `json:"formations"`
	Hasstats            bool        `json:"hasstats"`
	Injuries            int         `json:"injuries"`
	Inlivescore         bool        `json:"inlivescore"`
	Lineup              int         `json:"lineup"`
	Liveodds            bool        `json:"liveodds"`
	Livetable           int         `json:"livetable"`
	Lmtsupport          int         `json:"lmtsupport"`
	Matchdatacomplete   bool        `json:"matchdatacomplete"`
	Mediacoverage       bool        `json:"mediacoverage"`
	Multicast           bool        `json:"multicast"`
	Penaltyshootout     int         `json:"penaltyshootout"`
	Scoutconnected      bool        `json:"scoutconnected"`
	Scoutcoveragestatus int         `json:"scoutcoveragestatus"`
	Scoutmatch          int         `json:"scoutmatch"`
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
	Tzoffset int    `json:"tzoffset"`
	Uts      int    `json:"uts"`
}

type Periods struct {
	Ft *Ft `json:"ft,omitempty"`
	P1 Ft  `json:"p1"`
	Ot *Ft `json:"ot,omitempty"`
}

type Ft struct {
	Away int `json:"away"`
	Home int `json:"home"`
}

type Result struct {
	Away   *int    `json:"away"`
	Home   *int    `json:"home"`
	Winner *string `json:"winner"`
}

type Roundname struct {
	Doc                 string      `json:"_doc"`
	ID                  int         `json:"_id"`
	Cuproundnumber      interface{} `json:"cuproundnumber"`
	Displaynumber       interface{} `json:"displaynumber"`
	Shortname           *string     `json:"shortname,omitempty"`
	Statisticssortorder *int        `json:"statisticssortorder,omitempty"`
}

type Status struct {
	Doc  string `json:"_doc"`
	ID   int    `json:"_id"`
	Name string `json:"name"`
}

type Teams struct {
	Away Team `json:"away"`
	Home Team `json:"home"`
}

type Team struct {
	Doc        string  `json:"_doc"`
	ID         int     `json:"_id"`
	Abbr       string  `json:"abbr"`
	Haslogo    bool    `json:"haslogo"`
	Iscountry  bool    `json:"iscountry"`
	Mediumname string  `json:"mediumname"`
	Name       string  `json:"name"`
	Nickname   *string `json:"nickname"`
	Uid        int     `json:"uid"`
	Virtual    bool    `json:"virtual"`
}

type Timeinfo struct {
	Ended      *string     `json:"ended"`
	Injurytime interface{} `json:"injurytime"`
	Played     interface{} `json:"played"`
	Remaining  interface{} `json:"remaining"`
	Running    bool        `json:"running"`
	Started    interface{} `json:"started"`
}
