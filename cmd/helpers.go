package cmd

import (
	"fmt"
	"strconv"

	"github.com/etoro/etoro-cli/internal/api"
)

func parseInt(s string) (int, error) {
	return strconv.Atoi(s)
}

func errorf(format string, args ...any) error {
	return fmt.Errorf(format, args...)
}

func resolveInstrumentSymbols(client *api.Client, ids []int) map[int]string {
	symbolMap := make(map[int]string)
	if len(ids) == 0 {
		return symbolMap
	}

	unique := make(map[int]bool)
	for _, id := range ids {
		unique[id] = true
	}
	deduped := make([]int, 0, len(unique))
	for id := range unique {
		deduped = append(deduped, id)
	}

	resp, err := client.GetInstruments(deduped)
	if err != nil {
		for _, id := range deduped {
			symbolMap[id] = fmt.Sprintf("#%d", id)
		}
		return symbolMap
	}

	for _, inst := range resp.InstrumentDisplayDatas {
		if inst.Symbol != "" {
			symbolMap[inst.InstrumentID] = inst.Symbol
		} else {
			symbolMap[inst.InstrumentID] = inst.DisplayName
		}
	}

	for _, id := range deduped {
		if _, ok := symbolMap[id]; !ok {
			symbolMap[id] = fmt.Sprintf("#%d", id)
		}
	}

	return symbolMap
}
