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
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return listAllDestinations(app), cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			destination := args[0]
			message := args[1]
			return app.client.SendMessage(destination, message, alias)
		},
	}
	cmd.Flags().StringVar(&alias, "alias", alias, "Name under which the message appears.")
	return cmd
}

func listAllDestinations(app *app) []string {
	destinations := []string{}
	users, err := app.client.ListUsers()
	if err == nil {
		for _, u := range users {
			destinations = append(destinations, "@"+u.Username)
		}
	}

	channels, err := app.client.ListChannels()
	if err == nil {
		for _, c := range channels {
			destinations = append(destinations, "#"+c.Name)
		}
	}
	return destinations
}
