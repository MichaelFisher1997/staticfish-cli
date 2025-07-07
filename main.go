package main

import (
	"fmt"
	"os"

	"github.com/micqdf/data-collection/staticfish-cli/google-search"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "staticfish-cli",
	Short: "A CLI tool for interacting with StaticFish services.",
	Long:  `A longer description that spans multiple lines and likely contains examples and usage of using your application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Default action when no subcommand is provided
		fmt.Println("Welcome to StaticFish CLI! Use 'help' for more information.")
	},
}

var googleSearchCmd = &cobra.Command{
	Use:   "google-search",
	Short: "Performs a Google search and scrapes business data.",
	Long: `
Performs a Google search for a specified query, scrapes the resulting business websites for contact information, and stores the data in a PostgreSQL database.

This command requires a Google Maps API key and database credentials.
`,
	Example: `./staticfish-cli google-search \
  --api="YOUR_GOOGLE_MAPS_API_KEY" \
  --username="your_db_user" \
  --password="your_db_password" \
  --query="Bookstores in Manchester" \
  --pages=5`,
	Run: func(cmd *cobra.Command, args []string) {
		api, _ := cmd.Flags().GetString("api")
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		query, _ := cmd.Flags().GetString("query")
		pages, _ := cmd.Flags().GetInt("pages")

		googlesearch.Search(api, username, password, query, pages)
	},
}

var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Dumps the collected business data.",
	Long: `
Dumps the data from the 'businesses' table in the PostgreSQL database to a file.
The output format can be specified as either CSV or SQL.
`,
	Example: `./staticfish-cli dump \
  --username="your_db_user" \
  --password="your_db_password" \
  --format="sql"`,
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		format, _ := cmd.Flags().GetString("format")

		googlesearch.Dump(username, password, format)
	},
}

func init() {
	rootCmd.AddCommand(googleSearchCmd)
	rootCmd.AddCommand(dumpCmd)

	googleSearchCmd.Flags().String("api", "", "Your StaticFish API key.")
	googleSearchCmd.Flags().String("username", "", "Your StaticFish username.")
	googleSearchCmd.Flags().String("password", "", "Your StaticFish password.")
	googleSearchCmd.Flags().String("query", "", "The search query.")
	googleSearchCmd.Flags().Int("pages", 1, "The number of pages to search.")
	googleSearchCmd.MarkFlagRequired("api")
	googleSearchCmd.MarkFlagRequired("username")
	googleSearchCmd.MarkFlagRequired("password")
	googleSearchCmd.MarkFlagRequired("query")

	dumpCmd.Flags().String("username", "", "Your StaticFish username.")
	dumpCmd.Flags().String("password", "", "Your StaticFish password.")
	dumpCmd.Flags().String("format", "csv", "The output format (csv or sql).")
	dumpCmd.MarkFlagRequired("username")
	dumpCmd.MarkFlagRequired("password")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}