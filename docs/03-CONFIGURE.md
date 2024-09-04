# Configuration

## Server settings

* env:`TOOKHOOK_DATABASE_URI`, `-d`, `--database-uri`: database connection string. The rules store in the database. Currently, only YAML file is supported for storing rules.  (default "file:///hooks.yml").
* env:`TOOKHOOK_LISTEN`, `-s`, `--listen`: listen address and port in format [<host>]:<port> (default "localhost:8080").
* env:`TOOKHOOK_LOG_LEVEL`, `-l`, `--log-level`: log level (default "info").
* env:`TOOKHOOK_PLUGINS`, `-p`, `--plugins`: comma separated list of plugins. Example, `TOOKHOOK_PLUGINS=/app/plugin1,/app/plugin2`.

## Rules
