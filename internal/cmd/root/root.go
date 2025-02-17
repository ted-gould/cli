package root

import (
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/axiomhq/pkg/version"
	"github.com/spf13/cobra"

	"github.com/axiomhq/cli/internal/cmdutil"
	"github.com/axiomhq/cli/internal/config"

	// Core commands
	ingestCmd "github.com/axiomhq/cli/internal/cmd/ingest"
	queryCmd "github.com/axiomhq/cli/internal/cmd/query"
	streamCmd "github.com/axiomhq/cli/internal/cmd/stream"

	// Management commands
	configCmd "github.com/axiomhq/cli/internal/cmd/config"
	datasetCmd "github.com/axiomhq/cli/internal/cmd/dataset"
	organizationCmd "github.com/axiomhq/cli/internal/cmd/organization"
	tokenCmd "github.com/axiomhq/cli/internal/cmd/token"

	// Additional commands
	authCmd "github.com/axiomhq/cli/internal/cmd/auth"
	completionCmd "github.com/axiomhq/cli/internal/cmd/completion"
	versionCmd "github.com/axiomhq/cli/internal/cmd/version"
)

// NewRootCmd creates and returns the root command.
func NewRootCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "axiom <command> <subcommand>",
		Short: "Axiom CLI",
		Long:  "The power of Axiom on the command-line.",

		SilenceErrors: true,
		SilenceUsage:  true,

		Example: heredoc.Doc(`
			$ axiom auth login
			$ axiom version
			$ cat /var/log/nginx/*.log | axiom ingest nginx-logs
		`),

		Annotations: map[string]string{
			"help:credentials": heredoc.Doc(`
				See 'axiom help credentials' for help and guidance on authentication.
			`),
			"help:environment": heredoc.Doc(`
				See 'axiom help environment' for the list of supported environment variables.
			`),
		},

		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			if fl := cmd.Flag("config"); fl.Changed {
				if f.Config, err = config.Load(fl.Value.String()); err != nil {
					return err
				}
			}

			if fl := cmd.Flag("deployment"); fl.Changed {
				f.Config.ActiveDeployment = fl.Value.String()
			}
			if fl := cmd.Flag("org-id"); fl.Changed {
				f.Config.OrganizationIDOverride = fl.Value.String()
			}
			if fl := cmd.Flag("token"); fl.Changed {
				f.Config.TokenOverride = fl.Value.String()
			}
			if fl := cmd.Flag("url"); fl.Changed {
				f.Config.URLOverride = fl.Value.String()
			}

			f.Config.Insecure = cmd.Flag("insecure").Changed
			f.Config.ForceCloud = cmd.Flag("force-cloud").Changed
			f.IO.EnableActivityIndicator(!cmd.Flag("no-spinner").Changed)

			return nil
		},
	}

	// IO
	cmd.SetIn(f.IO.In())
	cmd.SetOut(f.IO.Out())
	cmd.SetErr(f.IO.ErrOut())

	// Version
	cmd.Flags().BoolP("version", "v", false, "Show axiom version")
	cmd.SetVersionTemplate("{{ .Short }} version {{ .Version }}\n")
	cmd.Version = version.Release()

	// Help & usage
	cmd.PersistentFlags().BoolP("help", "h", false, "Show help for command")
	cmd.SetHelpFunc(rootHelpFunc(f.IO))
	cmd.SetUsageFunc(rootUsageFunc)
	cmd.SetFlagErrorFunc(rootFlagErrrorFunc)

	// Overrides
	cmd.PersistentFlags().StringP("config", "C", "", "Path to configuration file to use")
	cmd.PersistentFlags().StringP("deployment", "D", "", "Deployment to use")
	cmd.PersistentFlags().StringP("org-id", "O", os.Getenv("AXIOM_ORG_ID"), "Organization ID to use (only valid for Axiom Cloud)")
	cmd.PersistentFlags().StringP("token", "T", os.Getenv("AXIOM_TOKEN"), "Token to use")
	cmd.PersistentFlags().StringP("url", "U", os.Getenv("AXIOM_URL"), "Url to use")
	cmd.PersistentFlags().BoolP("insecure", "I", false, "Bypass certificate validation")
	cmd.PersistentFlags().BoolP("force-cloud", "F", false, "Treat deployment as Axiom Cloud")
	cmd.PersistentFlags().Bool("no-spinner", false, "Disable the activity indicator")

	// Core commands
	cmd.AddCommand(ingestCmd.NewIngestCmd(f))
	cmd.AddCommand(queryCmd.NewQueryCmd(f))
	cmd.AddCommand(streamCmd.NewStreamCmd(f))

	// Management commands
	cmd.AddCommand(configCmd.NewConfigCmd(f))
	cmd.AddCommand(datasetCmd.NewDatasetCmd(f))
	cmd.AddCommand(organizationCmd.NewOrganizationCmd(f))
	cmd.AddCommand(tokenCmd.NewTokenCmd(f))

	// Additional commands
	cmd.AddCommand(authCmd.NewAuthCmd(f))
	cmd.AddCommand(completionCmd.NewCompletionCmd(f))
	cmd.AddCommand(versionCmd.NewVersionCmd(f, version.Print("Axiom CLI")))

	// Help topics
	cmd.AddCommand(newHelpTopic(f.IO, "credentials"))
	cmd.AddCommand(newHelpTopic(f.IO, "environment"))

	return cmd
}
