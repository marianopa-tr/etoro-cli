package output

import "encoding/json"

type TradeResult struct {
	Action   string          `json:"action"`
	Symbol   string          `json:"symbol"`
	IsDemo   bool            `json:"isDemo"`
	Response json.RawMessage `json:"response"`
}

func PrintTradeResult(result TradeResult, format Format) {
	if format == JSON {
		PrintJSON(result)
		return
	}

	label := ""
	if result.IsDemo {
		label = Yellow("[DEMO] ")
	}

	Successf("%s%s order for %s submitted successfully.", label, result.Action, Cyan(result.Symbol))
}

func PrintOrderCancelled(orderID int, format Format) {
	if format == JSON {
		PrintJSON(map[string]any{
			"action":  "cancelled",
			"orderId": orderID,
		})
		return
	}

	Successf("Order %d cancelled.", orderID)
}
