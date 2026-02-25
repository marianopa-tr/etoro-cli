package output

import (
	"encoding/json"
	"testing"
)

func TestPrintTradeResultTable(t *testing.T) {
	result := TradeResult{
		Action: "Buy", Symbol: "AAPL", IsDemo: false,
		Response: json.RawMessage(`{"orderId":123}`),
	}
	PrintTradeResult(result, Table)
}

func TestPrintTradeResultDemo(t *testing.T) {
	result := TradeResult{
		Action: "Buy", Symbol: "AAPL", IsDemo: true,
		Response: json.RawMessage(`{"orderId":123}`),
	}
	PrintTradeResult(result, Table)
}

func TestPrintTradeResultJSON(t *testing.T) {
	result := TradeResult{
		Action: "Sell", Symbol: "TSLA", IsDemo: false,
		Response: json.RawMessage(`{"orderId":456}`),
	}

	out := captureStdout(t, func() {
		PrintTradeResult(result, JSON)
	})

	var parsed map[string]any
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if parsed["action"] != "Sell" {
		t.Errorf("action = %v", parsed["action"])
	}
}

func TestPrintOrderCancelledTable(t *testing.T) {
	PrintOrderCancelled(123, Table)
}

func TestPrintOrderCancelledJSON(t *testing.T) {
	out := captureStdout(t, func() {
		PrintOrderCancelled(123, JSON)
	})

	var parsed map[string]any
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if parsed["orderId"] != float64(123) {
		t.Errorf("orderId = %v", parsed["orderId"])
	}
}
