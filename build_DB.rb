require_relative 'lib/pg_setup'

PgSetup::build_db
PgSetup::build_table(columns: [
  ['ts', 'timestamp'],
  ['src', 'bigint'],
  ['dst', 'bigint'],
  ['flags', 'bigint']
])
