# staticfish-cli

`staticfish-cli` is a command-line interface tool designed to collect business data. It uses the Google Maps API to search for specified business types in a location, scrapes their websites for contact information like emails, and stores the collected data into a PostgreSQL database.

## Features

- Searches Google Maps for businesses based on a query.
- Fetches detailed information for each business, including website URL and phone number.
- Scrapes the business's website to find an email address.
- Detects the technology used to build the website (e.g., WordPress, React, Shopify).
- Stores all collected information into a PostgreSQL database, avoiding duplicate entries.
- Dumps collected data to either a `.csv` or `.sql` file.
- Automated builds for Linux and Windows via GitHub Actions.

## Downloads

Pre-built binaries for Linux and Windows are automatically created for every push to the `main` branch. You can download them from the "Actions" tab in the GitHub repository.

1.  Navigate to the **Actions** tab in your GitHub repository.
2.  Click on the latest successful workflow run under the **Build** workflow.
3.  At the bottom of the workflow summary page, you will find the **Artifacts** section.
4.  Download the artifact for your operating system (`staticfish-cli-linux` or `staticfish-cli-windows`).

## Local Development

### Prerequisites

- Go (version 1.18 or later)
- A running PostgreSQL database
- A Google Maps API key

### Building the CLI

To build the application locally, run the provided shell script. This will compile the source code and create an executable file named `staticfish-cli` in the project root.

```bash
./build.sh
```

## Usage

The CLI has two primary subcommands: `google-search` and `dump`.

### `google-search`

This command executes the search and data collection process.

#### Flags

- `--api`: **(Required)** Your Google Maps API key.
- `--username`: **(Required)** The username for your PostgreSQL database.
- `--password`: **(Required)** The password for your PostgreSQL database.
- `--query`: **(Required)** The search query (e.g., "Bookstores in Manchester").
- `--pages`: (Optional) The number of pages to search. Defaults to 1.

#### Example

```bash
./staticfish-cli google-search \
  --api="AIzaSy...YOUR_GOOGLE_MAPS_API_KEY" \
  --username="your_db_user" \
  --password="your_db_password" \
  --query="Bookstores in Manchester" \
  --pages=5
```

### `dump`

This command dumps the data from the `businesses` table in the database to a file.

#### Flags

- `--username`: **(Required)** The username for your PostgreSQL database.
- `--password`: **(Required)** The password for your PostgreSQL database.
- `--format`: (Optional) The output format (`csv` or `sql`). Defaults to `csv`.

#### Example

```bash
./staticfish-cli dump \
  --username="your_db_user" \
  --password="your_db_password" \
  --format="sql"
```

