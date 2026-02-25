package cmd

import (
	"fmt"

	"github.com/etoro/etoro-cli/internal/api"
	"github.com/etoro/etoro-cli/internal/output"
	"github.com/etoro/etoro-cli/internal/resolver"
	"github.com/spf13/cobra"
)

var tradeCmd = &cobra.Command{
	Use:   "trade",
	Short: "Open, close, or manage trades",
	Long: `Execute trades on the eToro platform. Use --demo for virtual trading.

Examples:
  etoro trade open AAPL --amount 500 --leverage 2
  etoro trade open BTC --units 0.1 --demo
  etoro trade close 12345
  etoro trade limit TSLA --price 200 --amount 1000`,
}

var tradeOpenCmd = &cobra.Command{
	Use:   "open <symbol>",
	Short: "Open a market position",
	Long: `Open a new position at market price. Specify investment by --amount (USD)
or --units (number of units/shares).

Examples:
  etoro trade open AAPL --amount 500
  etoro trade open AAPL --amount 500 --leverage 5 --sl 150 --tp 200
  etoro trade open BTC --units 0.5 --demo
  etoro trade open EURUSD --amount 1000 --short`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := mustLoadConfig()
		if err := cfg.RequireAuth(); err != nil {
			return err
		}

		client := api.NewClient(cfg, flagDemo)
		res := resolver.New(client)

		instrumentID, symbol, err := res.Resolve(args[0])
		if err != nil {
			return err
		}

		amount, _ := cmd.Flags().GetFloat64("amount")
		units, _ := cmd.Flags().GetFloat64("units")
		leverage, _ := cmd.Flags().GetInt("leverage")
		sl, _ := cmd.Flags().GetFloat64("sl")
		tp, _ := cmd.Flags().GetFloat64("tp")
		short, _ := cmd.Flags().GetBool("short")
		tsl, _ := cmd.Flags().GetBool("tsl")

		isBuy := !short

		if amount == 0 && units == 0 {
			return errorf("specify --amount or --units")
		}

		if !flagDemo && !flagYes {
			msg := fmt.Sprintf("⚠  REAL TRADE: %s %s", direction(isBuy), symbol)
			if amount > 0 {
				msg += fmt.Sprintf(" for $%.2f", amount)
			} else {
				msg += fmt.Sprintf(" for %.4f units", units)
			}
			if !output.ConfirmDanger(msg) {
				output.Infof("Trade cancelled.")
				return nil
			}
		}

		var slPtr, tpPtr *float64
		var tslPtr *bool
		if sl > 0 {
			slPtr = &sl
		}
		if tp > 0 {
			tpPtr = &tp
		}
		if tsl {
			tslPtr = &tsl
		}

		var resp []byte
		if units > 0 {
			resp, err = client.OpenPositionByUnits(&api.OpenByUnitsRequest{
				InstrumentID:   instrumentID,
				IsBuy:          isBuy,
				Leverage:       leverage,
				Units:          units,
				StopLossRate:   slPtr,
				TakeProfitRate: tpPtr,
				IsTslEnabled:   tslPtr,
			})
		} else {
			resp, err = client.OpenPositionByAmount(&api.OpenByAmountRequest{
				InstrumentID:   instrumentID,
				IsBuy:          isBuy,
				Leverage:       leverage,
				Amount:         amount,
				StopLossRate:   slPtr,
				TakeProfitRate: tpPtr,
				IsTslEnabled:   tslPtr,
			})
		}

		if err != nil {
			return err
		}

		output.PrintTradeResult(output.TradeResult{
			Action:   direction(isBuy),
			Symbol:   symbol,
			IsDemo:   flagDemo,
			Response: resp,
		}, output.GetFormat())
		return nil
	},
}

var tradeCloseCmd = &cobra.Command{
	Use:   "close <positionId>",
	Short: "Close an open position",
	Long: `Close a position by its ID. Optionally close partially with --percent.

Examples:
  etoro trade close 12345
  etoro trade close 12345 --percent 50`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := mustLoadConfig()
		if err := cfg.RequireAuth(); err != nil {
			return err
		}

		client := api.NewClient(cfg, flagDemo)

		positionID, err := parseInt(args[0])
		if err != nil {
			return errorf("invalid position ID: %s", args[0])
		}

		if !flagDemo && !flagYes {
			msg := fmt.Sprintf("⚠  REAL TRADE: Close position %d", positionID)
			if !output.ConfirmDanger(msg) {
				output.Infof("Trade cancelled.")
				return nil
			}
		}

		req := &api.ClosePositionRequest{}

		resp, err := client.ClosePosition(positionID, req)
		if err != nil {
			return err
		}

		output.PrintTradeResult(output.TradeResult{
			Action:   "Close",
			Symbol:   fmt.Sprintf("position #%d", positionID),
			IsDemo:   flagDemo,
			Response: resp,
		}, output.GetFormat())
		return nil
	},
}

var tradeLimitCmd = &cobra.Command{
	Use:   "limit <symbol>",
	Short: "Place a limit order",
	Long: `Place a limit order to open a position when the price reaches a target.

Examples:
  etoro trade limit AAPL --price 170 --amount 500
  etoro trade limit BTC --price 50000 --amount 1000 --demo`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := mustLoadConfig()
		if err := cfg.RequireAuth(); err != nil {
			return err
		}

		client := api.NewClient(cfg, flagDemo)
		res := resolver.New(client)

		instrumentID, symbol, err := res.Resolve(args[0])
		if err != nil {
			return err
		}

		price, _ := cmd.Flags().GetFloat64("price")
		amount, _ := cmd.Flags().GetFloat64("amount")
		leverage, _ := cmd.Flags().GetInt("leverage")
		sl, _ := cmd.Flags().GetFloat64("sl")
		tp, _ := cmd.Flags().GetFloat64("tp")
		short, _ := cmd.Flags().GetBool("short")

		isBuy := !short

		if price == 0 {
			return errorf("--price is required for limit orders")
		}
		if amount == 0 {
			return errorf("--amount is required for limit orders")
		}

		if !flagDemo && !flagYes {
			msg := fmt.Sprintf("⚠  REAL LIMIT ORDER: %s %s at %.4f for $%.2f", direction(isBuy), symbol, price, amount)
			if !output.ConfirmDanger(msg) {
				output.Infof("Order cancelled.")
				return nil
			}
		}

		req := &api.LimitOrderRequest{
			InstrumentID: instrumentID,
			IsBuy:        isBuy,
			Leverage:     leverage,
			Amount:       amount,
			Rate:         price,
		}
		if sl <= 0 {
			sl = 0.01
		}
		req.StopLossRate = &sl
		if tp > 0 {
			req.TakeProfitRate = &tp
		} else {
			noTP := true
			req.IsNoTakeProfit = &noTP
		}

		resp, err := client.PlaceLimitOrder(req)
		if err != nil {
			return err
		}

		output.PrintTradeResult(output.TradeResult{
			Action:   "Limit " + direction(isBuy),
			Symbol:   symbol,
			IsDemo:   flagDemo,
			Response: resp,
		}, output.GetFormat())
		return nil
	},
}

func direction(isBuy bool) string {
	if isBuy {
		return "Buy"
	}
	return "Sell"
}

func init() {
	tradeOpenCmd.Flags().Float64("amount", 0, "investment amount in USD")
	tradeOpenCmd.Flags().Float64("units", 0, "number of units/shares")
	tradeOpenCmd.Flags().Int("leverage", 1, "leverage multiplier")
	tradeOpenCmd.Flags().Float64("sl", 0, "stop-loss rate")
	tradeOpenCmd.Flags().Float64("tp", 0, "take-profit rate")
	tradeOpenCmd.Flags().Bool("short", false, "open a short (sell) position")
	tradeOpenCmd.Flags().Bool("tsl", false, "enable trailing stop-loss")

	tradeLimitCmd.Flags().Float64("price", 0, "target price for the limit order")
	tradeLimitCmd.Flags().Float64("amount", 0, "investment amount in USD")
	tradeLimitCmd.Flags().Int("leverage", 1, "leverage multiplier")
	tradeLimitCmd.Flags().Float64("sl", 0, "stop-loss rate")
	tradeLimitCmd.Flags().Float64("tp", 0, "take-profit rate")
	tradeLimitCmd.Flags().Bool("short", false, "short (sell) limit order")

	tradeCmd.AddCommand(tradeOpenCmd)
	tradeCmd.AddCommand(tradeCloseCmd)
	tradeCmd.AddCommand(tradeLimitCmd)
	rootCmd.AddCommand(tradeCmd)
}
