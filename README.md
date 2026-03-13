# Bookkeeper

## Setup Config
Create and update the config file using the following yaml and place it somewhere (eg `~/.config/bookkeeper/config.yaml`). The file path will be supplied to each app via command line argument.
```yaml
---
database:
  host: localhost
  port: 5432
  name: bookkeeper
  user: someuser
  password: somepassword
google:
  auth:
    dir: ~/.config/googleauth/
  calendars:
    - name: PWHL Hockey
      id: something@group.calendar.google.com
server:
  host: localhost
  port: 4000
```

## Setup Data
See the [Test Data](testdata/) directory. Create and update a data file using [data.yaml](data.yaml) as a base. The file path will be supplied to the initializer app via command line argument.

Specify currencies, bank accounts (including paths to csv files to import), accounting accounts and assigners (rules for automatically categorizing transactions to accounting accounts).

Specify sports registrations, memberships and school invoices.

## Setup Database
```bash
# Reference: https://petereisentraut.blogspot.com/2010/03/running-sql-scripts-with-psql.html
sudo -u postgres psql -X -q -1 -v ON_ERROR_STOP=1 -d bookkeeper -f drop.sql -f bank.sql -f book.sql -f sports.sql -f school.sql
```

## Initialize Database and Import Bank Transactions
```bash
# using sample test data
go run ./cmd/initializer/main.go -config ~/.config/bookkeeper/config.yaml
# OR user provided data
go run ./cmd/initializer/main.go -config ~/.config/bookkeeper/config.yaml -data ./data/data.yaml
```

## Import Bank Transactions into Accounting Journal
```bash
go run ./cmd/journal-importer/main.go -config ~/.config/bookkeeper/config.yaml
```

## Run Web Server
```bash
go run ./cmd/web/main.go -config ~/.config/bookkeeper/config.yaml
## OR with logs outputted to files
go run ./cmd/web/main.go -config ~/.config/bookkeeper/config.yaml >>/tmp/info.log 2>>/tmp/error.log
```

Navigate to [http://localhost:4000](http://localhost:4000).
