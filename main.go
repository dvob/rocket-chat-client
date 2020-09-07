package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dsbrng25b/rocket-chat-client/pkg/rocketchat"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	envVarPrefix = "RC_"
)

func main() {
	err := newRootCmd().Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type app struct {
	client *rocketchat.Client
}

func newRootCmd() *cobra.Command {
	var (
		app       = &app{}
		serverURL = "https://localhost:3000"
		userID    string
		token     string
	)
	cmd := &cobra.Command{
		Use:   "rocket-chat-client",
		Short: "rocket-chat-client allows you to interact with Rocket.Chat",
		Long: `The rocket-chat-client allows you to list users and channels and send
messages to them. All options can also be set as environment variables
(e.g. --user-id becomse RC_USER_ID).`,
		TraverseChildren: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			cmd.Flags().VisitAll(func(f *pflag.Flag) {
				optName := strings.ToUpper(f.Name)
				optName = strings.ReplaceAll(optName, "-", "_")
				varName := envVarPrefix + optName
				if val, ok := os.LookupEnv(varName); ok && !f.Changed {
					err = f.Value.Set(val)
				}
			})

			app.client = rocketchat.NewClient(serverURL, userID, token)
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true
			return err
		},
	}
	cmd.PersistentFlags().StringVar(&serverURL, "url", serverURL, "URL of Rocket.Chat")
	cmd.PersistentFlags().StringVar(&userID, "user-id", userID, "The ID of the user")
	cmd.PersistentFlags().StringVar(&token, "token", token, "Authentication Token")
	cmd.AddCommand(
		newListCmd(app),
		newSendMsgCmd(app),
		newCompletionCmd(),
		newVersionCmd(),
	)
	return cmd
}
