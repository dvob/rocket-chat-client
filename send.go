package main

import (
	"github.com/spf13/cobra"
)

func newSendMsgCmd(app *app) *cobra.Command {
	var (
		alias string
	)
	cmd := &cobra.Command{
		Use:   "send <user|channel> <message>",
		Short: "Send a message to a user or a channel.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			destination := args[0]
			message := args[1]
			return app.client.SendMessage(destination, message, alias)
		},
	}
	cmd.Flags().StringVar(&alias, "alias", alias, "Name under which the message appears.")
	return cmd
}
