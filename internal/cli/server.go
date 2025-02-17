package cli

import (
	"github.com/euiko/go-fullstack-boilerplate/internal/core/webapp"
	"github.com/spf13/cobra"
)

func Server(app *webapp.App) func(settings *webapp.Settings) webapp.Module {
	return func(settings *webapp.Settings) webapp.Module {
		return webapp.NewModule(webapp.WithCLI(func(cmd *cobra.Command) {
			startCmd := cobra.Command{
				Use:   "start",
				Short: "Start the web application",
				RunE: func(cmd *cobra.Command, args []string) error {
					return app.Start(cmd.Context())
				},
			}
			cmd.AddCommand(&startCmd)
		}))
	}
}
