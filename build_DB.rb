require_relative 'lib/pg_setup'

# add swarm64 user
PgSetup::add_user

# build test database
PgSetup::build_db
PgSetup::build_table(columns: [
  ['ts', 'timestamp'],
  ['src', 'bigint'],
  ['dst', 'bigint'],
  ['flags', 'bigint']
])


# build benchmark database
PgSetup::build_db(database:'swarm_benchmark')
PgSetup::build_table(
  table:'ingestion_benchmark',
  database:'swarm_benchmark',
  columns: [
    ['ts', 'timestamp'],
    ['src', 'bigint'],
    ['dst', 'bigint'],
    ['flags', 'bigint']
])
