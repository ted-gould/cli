package root

// Derived from https://github.com/cli/cli/blob/trunk/pkg/cmd/root/help_topic.go

import (
	"github.com/spf13/cobra"

	"github.com/axiomhq/cli/pkg/terminal"
)

var topics = map[string]string{
	"credentials": `
		The supplied access token can either be an API token, retrieved from the
		settings page of the Axiom deployment (Settings -> API Tokens) or a
		Personal Access token, retrieved from the users profile page
		(Settings -> Profile).

		Be aware, that API tokens are only valid for ingestion and querying,
		depending on their permissions! Using them with Axiom CLI is encouraged
		for ingest-only and/or query-only situations but renders the CLI unable
		to do anything else. Use a Personal Access Token to get full access to
		the deployment.
	`,

	"environment": `
		AXIOM_DEPLOYMENT: The deployment to use. Overwrittes the choice loaded
		from the configuration file.

		AXIOM_ORG_ID: The organization id of the organization the access token
		is valid for. Only valid for Axiom Cloud.
		
		AXIOM_PAGER, PAGER (in order of precedence): A terminal paging program
		to send standard output to, e.g. "less".

		AXIOM_TOKEN: Token The access token to use. Overwrittes the choice
		loaded from the configuration file.

		AXIOM_URL: The deployment url to use. Overwrittes the choice loaded from
		the configuration file.

		VISUAL, EDITOR (in order of precedence): The editor to use for authoring
		text.

		NO_COLOR: Set to any value to avoid printing ANSI escape sequences for
		color output.

		CLICOLOR: Set to "0" to disable printing ANSI colors in output.

		CLICOLOR_FORCE: Set to a value other than "0" to keep ANSI colors in
		output even when the output is piped.
	`,
}

func newHelpTopic(io *terminal.IO, topic string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  topic,
		Long: io.Doc(topics[topic]),

		Hidden: true,

		Run: helpTopicHelpFunc,
	}

	cmd.SetHelpFunc(helpTopicHelpFunc)
	cmd.SetUsageFunc(helpTopicUsageFunc)

	return cmd
}

func helpTopicHelpFunc(cmd *cobra.Command, args []string) {
	cmd.Print(cmd.Long)
}

func helpTopicUsageFunc(cmd *cobra.Command) error {
	cmd.Printf("Usage: axiom help %s", cmd.Use)
	return nil
}
