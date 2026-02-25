package cmd

import (
	"fmt"

	"github.com/etoro/etoro-cli/internal/api"
	"github.com/etoro/etoro-cli/internal/output"
	"github.com/spf13/cobra"
)

var ordersCmd = &cobra.Command{
	Use:   "orders",
	Short: "List and manage pending orders",
	Long: `View and cancel pending orders.

Examples:
  etoro orders list
  etoro orders cancel 12345
  etoro orders cancel-all`,
}

var ordersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all pending orders",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := mustLoadConfig()
		if err := cfg.RequireAuth(); err != nil {
			return err
		}

		client := api.NewClient(cfg, flagDemo)

		resp, err := client.GetPortfolio()
		if err != nil {
			return err
		}

		p := resp.ClientPortfolio

		var allIDs []int
		for _, o := range p.Orders {
			allIDs = append(allIDs, o.InstrumentID)
		}
		for _, o := range p.OrdersForOpen {
			allIDs = append(allIDs, o.InstrumentID)
		}
		symbols := resolveInstrumentSymbols(client, allIDs)

		rows := make([]output.OrderRow, 0)

		for _, o := range p.Orders {
			dir := "Long"
			if !o.IsBuy {
				dir = "Short"
			}
			rows = append(rows, output.OrderRow{
				OrderID:   o.OrderID,
				Symbol:    symbols[o.InstrumentID],
				Direction: dir,
				Amount:    o.Amount,
				Rate:      o.Rate,
				Leverage:  o.Leverage,
				SL:        o.StopLossRate,
				TP:        o.TakeProfitRate,
				OpenDate:  formatTime(o.OpenDateTime),
			})
		}

		for _, o := range p.OrdersForOpen {
			dir := "Long"
			if !o.IsBuy {
				dir = "Short"
			}
			rows = append(rows, output.OrderRow{
				OrderID:   o.OrderID,
				Symbol:    symbols[o.InstrumentID],
				Direction: dir,
				Amount:    o.Amount,
				Leverage:  o.Leverage,
				SL:        o.StopLossRate,
				TP:        o.TakeProfitRate,
				OpenDate:  formatTime(o.OpenDateTime),
			})
		}

		output.PrintOrders(rows, output.GetFormat())
		return nil
	},
}

var ordersCancelCmd = &cobra.Command{
	Use:   "cancel <orderId>",
	Short: "Cancel a pending order",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := mustLoadConfig()
		if err := cfg.RequireAuth(); err != nil {
			return err
		}

		client := api.NewClient(cfg, flagDemo)

		orderID, err := parseInt(args[0])
		if err != nil {
			return errorf("invalid order ID: %s", args[0])
		}

		if !flagYes {
			if !output.Confirm(fmt.Sprintf("Cancel order %d?", orderID)) {
				output.Infof("Cancelled.")
				return nil
			}
		}

		_, err = client.CancelLimitOrder(orderID)
		if err != nil {
			_, err = client.CancelOrder(orderID)
		}
		if err != nil {
			return err
		}

		output.PrintOrderCancelled(orderID, output.GetFormat())
		return nil
	},
}

var ordersCancelAllCmd = &cobra.Command{
	Use:   "cancel-all",
	Short: "Cancel all pending orders",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := mustLoadConfig()
		if err := cfg.RequireAuth(); err != nil {
			return err
		}

		client := api.NewClient(cfg, flagDemo)

		resp, err := client.GetPortfolio()
		if err != nil {
			return err
		}

		p := resp.ClientPortfolio
		total := len(p.Orders) + len(p.OrdersForOpen)

		if total == 0 {
			output.Infof("No pending orders to cancel.")
			return nil
		}

		if !flagYes {
			if !output.ConfirmDanger(fmt.Sprintf("Cancel ALL %d pending orders?", total)) {
				output.Infof("Cancelled.")
				return nil
			}
		}

		cancelled := 0
		for _, o := range p.Orders {
			if _, err := client.CancelLimitOrder(o.OrderID); err != nil {
				if _, err2 := client.CancelOrder(o.OrderID); err2 != nil {
					output.Errorf("failed to cancel order %d: %s", o.OrderID, err2)
					continue
				}
			}
			cancelled++
		}
		for _, o := range p.OrdersForOpen {
			if _, err := client.CancelOrder(o.OrderID); err != nil {
				output.Errorf("failed to cancel order %d: %s", o.OrderID, err)
				continue
			}
			cancelled++
		}

		if output.GetFormat() == output.JSON {
			output.PrintJSON(map[string]any{
				"cancelled": cancelled,
				"total":     total,
			})
		} else {
			output.Successf("Cancelled %d of %d orders.", cancelled, total)
		}
		return nil
	},
}

func init() {
	ordersCmd.AddCommand(ordersListCmd)
	ordersCmd.AddCommand(ordersCancelCmd)
	ordersCmd.AddCommand(ordersCancelAllCmd)
	rootCmd.AddCommand(ordersCmd)
}
