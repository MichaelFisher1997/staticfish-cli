# staticfish-cli

`staticfish-cli` is a command-line interface tool designed to collect business data. It uses the Google Maps API to search for specified business types in a location, scrapes their websites for contact information like emails, and stores the collected data into a PostgreSQL database.

## Features

- Searches Google Maps for businesses based on a query.
- Fetches detailed information for each business, including website URL and phone number.
- Scrapes the business's website to find an email address.
- Detects the technology used to build the website (e.g., WordPress, React, Shopify).
- Stores all collected information into a PostgreSQL database, avoiding duplicate entries.

## Getting Started

### Prerequisites

- Go (version 1.18 or later)
- A running PostgreSQL database
- A Google Maps API key

### Building the CLI

To build the application, run the provided shell script. This will compile the source code and create an executable file named `staticfish-cli` in the project root.

```bash
./build.sh
```

## Usage

The primary functionality is handled by the `google-search` subcommand.

### `google-search`

This command executes the search and data collection process. It requires three flags to operate: your Google Maps API key and the credentials for your PostgreSQL database.

#### Flags

- `--api`: **(Required)** Your Google Maps API key.
- `--username`: **(Required)** The username for your PostgreSQL database.
- `--password`: **(Required)** The password for your PostgreSQL database.
- `--query`: **(Required)** The search query (e.g., "Bookstores in Manchester").
- `--pages`: (Optional) The number of pages to search. Defaults to 1.

#### Example

Here is an example of how to run the `google-search` command. Replace the placeholder values with your actual credentials.

```bash
./staticfish-cli google-search \
  --api="AIzaSy...YOUR_GOOGLE_MAPS_API_KEY" \
  --username="your_db_user" \
  --password="your_db_password" \
  --query="Bookstores in Manchester" \
  --pages=5
```

Upon execution, the tool will begin fetching data and logging its progress to the console, indicating which businesses are being processed.

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

