package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newCompletionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script",
		Long: `To load completions:

Bash:

$ source <(rocket-chat-client completion bash)

# To load completions for each session, execute once:
Linux:
rocket-chat-client completion bash > /etc/bash_completion.d/rocket-chat-client
MacOS:
$ rocket-chat-client completion bash > /usr/local/etc/bash_completion.d/rocket-chat-client

Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ rocket-chat-client completion zsh > "${fpath[1]}/_rocket-chat-client"

# You will need to start a new shell for this setup to take effect.

Fish:

$ rocket-chat-client completion fish | source

# To load completions for each session, execute once:
$ rocket-chat-client completion fish > ~/.config/fish/completions/rocket-chat-client.fish
`,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactArgs(1),
		Hidden:                true,
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			shell := args[0]
			switch shell {
			case "bash":
				err = newRootCmd().GenBashCompletion(os.Stdout)
			case "zsh":
				err = newRootCmd().GenZshCompletion(os.Stdout)
			case "fish":
				err = newRootCmd().GenFishCompletion(os.Stdout, true)
			case "powershell":
				err = newRootCmd().GenPowerShellCompletion(os.Stdout)
			default:
				err = fmt.Errorf("unknown shell: %s", shell)
			}
			return err
		},
	}
	return cmd
}
