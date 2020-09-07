package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newListCmd(app *app) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List users or channels",
	}
	cmd.AddCommand(
		newListUsersCmd(app),
		newListChannelsCmd(app),
	)
	return cmd
}

func newListUsersCmd(app *app) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "List users.",
		RunE: func(cmd *cobra.Command, args []string) error {
			users, err := app.client.ListUsers()
			if err != nil {
				return err
			}

			for _, user := range users {
				fmt.Println(user.Username)
			}
			return nil
		},
	}
	return cmd
}

func newListChannelsCmd(app *app) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "channel",
		Short: "List channels.",
		RunE: func(cmd *cobra.Command, args []string) error {
			channels, err := app.client.ListChannels()
			if err != nil {
				return err
			}

			for _, channel := range channels {
				fmt.Println(channel.Name)
			}
			return nil
		},
	}
	return cmd
}
