require_relative 'lib/pg_setup'

# lol, don't do this ;)

begin
# add swarm64 user
PgSetup::add_user
rescue
end

begin
# build test database
PgSetup::build_db
rescue
end

begin
PgSetup::build_table(columns: [
  ['ts', 'timestamp'],
  ['src', 'bigint'],
  ['dst', 'bigint'],
  ['flags', 'bigint']
])
rescue
end

begin
# build benchmark database
PgSetup::build_db(database:'swarm_benchmark')
rescue
end

begin
PgSetup::build_table(
  table:'ingestion_benchmark',
  database:'swarm_benchmark',
  columns: [
    ['ts', 'timestamp'],
    ['src', 'bigint'],
    ['dst', 'bigint'],
    ['flags', 'bigint']
])
rescue
end
