package history

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
	"github.com/twistedogic/doom/pkg/model"
)

const (
	csvDateFormat = "02/01/06"
	keyDateFormat = "20060201"
	tag           = "csv"

	Type model.Type = "history"
)

func ToSnakeCase(s string) string {
	lower := strings.ToLower(s)
	return strings.ReplaceAll(lower, " ", "_")
}

type Model struct {
	LeagueDivision            string  `csv:"Div"`
	MatchDate                 string  `csv:"Date"`
	Home                      string  `csv:"HomeTeam"`
	Away                      string  `csv:"AwayTeam"`
	HomeGoals                 int     `csv:"HG"`
	AwayGoals                 int     `csv:"AG"`
	FullTimeResult            string  `csv:"FTR"`
	HalfTimeResult            string  `csv:"HTR"`
	FullTimeHomeGoals         int     `csv:"FTHG"`
	FullTimeAwayGoals         int     `csv:"FTAG"`
	HalfTimeHomeGoals         int     `csv:"HTHG"`
	HalfTimeAwayGoals         int     `csv:"HTAG"`
	HomeTeamShots             int     `csv:"HS"`
	AwayTeamShots             int     `csv:"AS"`
	HomeTeamShotsonTarget     int     `csv:"HST"`
	AwayTeamShotsonTarget     int     `csv:"AST"`
	HomeTeamHitWoodwork       int     `csv:"HHW"`
	AwayTeamHitWoodwork       int     `csv:"AHW"`
	HomeTeamCorners           int     `csv:"HC"`
	AwayTeamCorners           int     `csv:"AC"`
	HomeTeamFoulsCommitted    int     `csv:"HF"`
	AwayTeamFoulsCommitted    int     `csv:"AF"`
	HomeTeamFreeKicksConceded int     `csv:"HFKC"`
	AwayTeamFreeKicksConceded int     `csv:"AFKC"`
	HomeTeamOffsides          int     `csv:"HO"`
	AwayTeamOffsides          int     `csv:"AO"`
	HomeTeamYellowCards       int     `csv:"HY"`
	AwayTeamYellowCards       int     `csv:"AY"`
	HomeTeamRedCards          int     `csv:"HR"`
	AwayTeamRedCards          int     `csv:"AR"`
	HomeTeamBookingsPoints    int     `csv:"HBP"`
	AwayTeamBookingsPoints    int     `csv:"ABP"`
	MarketMaximumHome         float64 `csv:"MaxH"`
	MarketMaximumDraw         float64 `csv:"MaxD"`
	MarketMaximumAway         float64 `csv:"MaxA"`
	MarketAverageHome         float64 `csv:"AvgH"`
	MarketAverageDraw         float64 `csv:"AvgD"`
	MarketAverageAway         float64 `csv:"AvgA"`
	BetMaximumHome            float64 `csv:"BbMxH"`
	BetMaximumDraw            float64 `csv:"BbMxD"`
	BetMaximumAway            float64 `csv:"BbMxA"`
	BetAverageHome            float64 `csv:"BbAvH"`
	BetAverageDraw            float64 `csv:"BbAvD"`
	BetAverageAway            float64 `csv:"BbAvA"`
}

func (m Model) Item(i *model.Item) error {
	b, err := json.Marshal(&m)
	if err != nil {
		return err
	}
	date, err := time.Parse(csvDateFormat, m.MatchDate)
	if err != nil {
		fmt.Println(m.MatchDate)
		return err
	}
	key := fmt.Sprintf(
		"%s|%s|%s",
		date.Format(keyDateFormat),
		ToSnakeCase(m.Home),
		ToSnakeCase(m.Away),
	)
	i.Key, i.Type, i.Data = key, Type, b
	return nil
}

func Transform(r io.Reader, encoder model.Encoder) error {
	reader := csv.NewReader(r)
	indexMap := make(map[string]int)
	fieldMap := make(map[string]string)
	for _, field := range structs.Fields(Model{}) {
		name := field.Name()
		fieldMap[name] = field.Tag(tag)
	}
	for {
		isHeader := false
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if len(indexMap) == 0 {
			for i, col := range row {
				indexMap[col] = i
			}
			continue
		}
		for col, i := range indexMap {
			if len(row) > i && col == row[i] {
				isHeader = true
			}
		}
		if isHeader {
			continue
		}
		container := make(map[string]string)
		for name, tag := range fieldMap {
			if i, ok := indexMap[tag]; len(row) > i && ok {
				container[name] = row[i]
			}
		}
		model := new(Model)
		if err := mapstructure.WeakDecode(container, model); err != nil {
			return err
		}
		if err := encoder.Encode(model); err != nil {
			return err
		}
	}
	return nil
}
