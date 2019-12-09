package odd

type Odd []OddElement

type OddElement struct {
	Matches []Match `json:"matches,omitempty"`
}

type Match struct {
	MatchID           string             `json:"matchID"`
	MatchIDinofficial string             `json:"matchIDinofficial"`
	MatchNum          string             `json:"matchNum"`
	MatchDate         string             `json:"matchDate"`
	MatchDay          string             `json:"matchDay"`
	Coupon            Coupon             `json:"coupon"`
	League            *League            `json:"league,omitempty"`
	HomeTeam          Team               `json:"homeTeam"`
	AwayTeam          Team               `json:"awayTeam"`
	MatchStatus       *string            `json:"matchStatus,omitempty"`
	MatchTime         string             `json:"matchTime"`
	Statuslastupdated *string            `json:"statuslastupdated,omitempty"`
	Inplaydelay       *string            `json:"inplaydelay,omitempty"`
	Channel           []Channel          `json:"channel"`
	LiveEvent         *LiveEvent         `json:"liveEvent,omitempty"`
	Accumulatedscore  []Accumulatedscore `json:"accumulatedscore"`
	Livescore         *Livescore         `json:"livescore,omitempty"`
	Cornerresult      *string            `json:"cornerresult,omitempty"`
	Cur               string             `json:"Cur"`
	IsDef             *string            `json:"isDef,omitempty"`
	HasWebTV          bool               `json:"hasWebTV"`
	Hadodds           *HadoddsClass      `json:"hadodds,omitempty"`
	Fhaodds           *HadoddsClass      `json:"fhaodds,omitempty"`
	Crsodds           *Crsodds           `json:"crsodds,omitempty"`
	Ooeodds           *ChloddsClass      `json:"ooeodds,omitempty"`
	Ttgodds           *Ttgodds           `json:"ttgodds,omitempty"`
	Hftodds           *Hftodds           `json:"hftodds,omitempty"`
	Hhaodds           *HadoddsClass      `json:"hhaodds,omitempty"`
	Hdcodds           *HadoddsClass      `json:"hdcodds,omitempty"`
	Hilodds           *ChloddsClass      `json:"hilodds,omitempty"`
	Chlodds           *ChloddsClass      `json:"chlodds,omitempty"`
	HasExtraTimePools bool               `json:"hasExtraTimePools"`
	Results           Results            `json:"results"`
	DefinedPools      []string           `json:"definedPools"`
	InplayPools       []string           `json:"inplayPools"`
}

type Accumulatedscore struct {
	Periodvalue  string `json:"periodvalue"`
	Periodstatus string `json:"periodstatus"`
	Home         string `json:"home"`
	Away         string `json:"away"`
}

type Team struct {
	TeamID     string `json:"teamID"`
	TeamNameCH string `json:"teamNameCH"`
	TeamNameEN string `json:"teamNameEN"`
}

type Channel struct {
	Order         int64  `json:"order"`
	ChannelID     string `json:"channelID"`
	ChannelNameCH string `json:"channelNameCH"`
	ChannelNameEN string `json:"channelNameEN"`
}

type ChloddsClass struct {
	Linelist   []Linelist `json:"LINELIST"`
	ID         string     `json:"ID"`
	Poolstatus string     `json:"POOLSTATUS"`
	Et         string     `json:"ET"`
	Inplay     string     `json:"INPLAY"`
	Allup      string     `json:"ALLUP"`
	Cur        string     `json:"Cur"`
	O          string     `json:"O,omitempty"`
	E          string     `json:"E,omitempty"`
}

type Linelist struct {
	Linenum    string `json:"LINENUM"`
	Mainline   string `json:"MAINLINE"`
	Linestatus string `json:"LINESTATUS"`
	Lineorder  string `json:"LINEORDER"`
	Line       string `json:"LINE"`
	H          string `json:"H"`
	L          string `json:"L"`
}

type Coupon struct {
	CouponID        string `json:"couponID"`
	CouponShortName string `json:"couponShortName"`
	CouponNameCH    string `json:"couponNameCH"`
	CouponNameEN    string `json:"couponNameEN"`
}

type Crsodds struct {
	S0003      string `json:"S0003"`
	S0104      string `json:"S0104"`
	S0100      string `json:"S0100"`
	S0401      string `json:"S0401"`
	S0400      string `json:"S0400"`
	Sm1Ma      string `json:"SM1MA"`
	S0302      string `json:"S0302"`
	S0102      string `json:"S0102"`
	S0200      string `json:"S0200"`
	S0001      string `json:"S0001"`
	S0000      string `json:"S0000"`
	Sm1Md      string `json:"SM1MD"`
	S0502      string `json:"S0502"`
	S0105      string `json:"S0105"`
	S0501      string `json:"S0501"`
	S0005      string `json:"S0005"`
	S0500      string `json:"S0500"`
	S0204      string `json:"S0204"`
	S0402      string `json:"S0402"`
	Sm1Mh      string `json:"SM1MH"`
	S0201      string `json:"S0201"`
	S0300      string `json:"S0300"`
	S0101      string `json:"S0101"`
	S0303      string `json:"S0303"`
	S0002      string `json:"S0002"`
	S0301      string `json:"S0301"`
	S0203      string `json:"S0203"`
	S0004      string `json:"S0004"`
	S0103      string `json:"S0103"`
	S0202      string `json:"S0202"`
	S0205      string `json:"S0205"`
	ID         string `json:"ID"`
	Poolstatus string `json:"POOLSTATUS"`
	Et         string `json:"ET"`
	Inplay     string `json:"INPLAY"`
	Allup      string `json:"ALLUP"`
	Cur        string `json:"Cur"`
}

type HadoddsClass struct {
	D          string `json:"D,omitempty"`
	H          string `json:"H"`
	A          string `json:"A"`
	ID         string `json:"ID"`
	Poolstatus string `json:"POOLSTATUS"`
	Et         string `json:"ET"`
	Inplay     string `json:"INPLAY"`
	Allup      string `json:"ALLUP"`
	Cur        string `json:"Cur"`
	Hg         string `json:"HG,omitempty"`
	Ag         string `json:"AG,omitempty"`
}

type Hftodds struct {
	Hh         string `json:"HH"`
	Dh         string `json:"DH"`
	HD         string `json:"HD"`
	DD         string `json:"DD"`
	Aa         string `json:"AA"`
	Ha         string `json:"HA"`
	Ad         string `json:"AD"`
	Da         string `json:"DA"`
	Ah         string `json:"AH"`
	ID         string `json:"ID"`
	Poolstatus string `json:"POOLSTATUS"`
	Et         string `json:"ET"`
	Inplay     string `json:"INPLAY"`
	Allup      string `json:"ALLUP"`
	Cur        string `json:"Cur"`
}

type League struct {
	LeagueID        string `json:"leagueID"`
	LeagueShortName string `json:"leagueShortName"`
	LeagueNameCH    string `json:"leagueNameCH"`
	LeagueNameEN    string `json:"leagueNameEN"`
}

type LiveEvent struct {
	IlcLiveDisplay  bool        `json:"ilcLiveDisplay"`
	HasLiveInfo     bool        `json:"hasLiveInfo"`
	IsIncomplete    bool        `json:"isIncomplete"`
	MatchIDbetradar string      `json:"matchIDbetradar"`
	Matchstate      string      `json:"matchstate"`
	StateTS         string      `json:"stateTS"`
	Liveevent       []Liveevent `json:"liveevent"`
}

type Liveevent struct {
	Order          int64  `json:"order"`
	MinutesElasped string `json:"minutesElasped"`
	ActionType     string `json:"actionType"`
	PlayerNameCH   string `json:"playerNameCH"`
	PlayerNameEN   string `json:"playerNameEN"`
	Homeaway       string `json:"homeaway"`
}

type Livescore struct {
	Home string `json:"home"`
	Away string `json:"away"`
}

type Results struct {
}

type Ttgodds struct {
	P0         string `json:"P0"`
	M7         string `json:"M7"`
	P5         string `json:"P5"`
	P2         string `json:"P2"`
	P1         string `json:"P1"`
	P4         string `json:"P4"`
	P6         string `json:"P6"`
	P3         string `json:"P3"`
	ID         string `json:"ID"`
	Poolstatus string `json:"POOLSTATUS"`
	Et         string `json:"ET"`
	Inplay     string `json:"INPLAY"`
	Allup      string `json:"ALLUP"`
	Cur        string `json:"Cur"`
}
