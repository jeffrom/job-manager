package commands

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/jeffrom/job-manager/mjob/client"
)

type CompletionCmd struct {
	*cobra.Command
}

func (c *CompletionCmd) Cmd() *cobra.Command { return c.Command }
func (c *CompletionCmd) Execute(ctx context.Context, cfg *client.Config, cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Usage()
	}
	switch args[0] {
	case "bash":
		cmd.Root().GenBashCompletion(os.Stdout)
	case "zsh":
		cmd.Root().GenZshCompletion(os.Stdout)
	case "fish":
		cmd.Root().GenFishCompletion(os.Stdout, true)
	case "powershell":
		cmd.Root().GenPowerShellCompletion(os.Stdout)
	}
	return nil
}

func newCompletionCmd(cfg *client.Config) *CompletionCmd {
	c := &CompletionCmd{
		Command: &cobra.Command{
			Use:   "completion [bash|zsh|fish|powershell]",
			Short: "generate completion script",
			Long: `To load completions:

Bash:

$ source <(yourprogram completion bash)

# To load completions for each session, execute once:
Linux:
  $ yourprogram completion bash > /etc/bash_completion.d/yourprogram
MacOS:
  $ yourprogram completion bash > /usr/local/etc/bash_completion.d/yourprogram

Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ yourprogram completion zsh > "${fpath[1]}/_yourprogram"

# You will need to start a new shell for this setup to take effect.

Fish:

$ yourprogram completion fish | source

# To load completions for each session, execute once:
$ yourprogram completion fish > ~/.config/fish/completions/yourprogram.fish
`,
			DisableFlagsInUseLine: true,
			ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
			Args:                  cobra.MaximumNArgs(1),
		},
	}

	return c
}
