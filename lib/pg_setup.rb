# insert_ts_updater and insert_ts_trigger were shamelessly taken from:
#   http://www.revsys.com/blog/2006/aug/04/automatically-updating-a-timestamp-column-in-postgresql/

# To make this better I would make a connect_to function which returns the connection
# and queue all these up
# in prepared statements to be executed in a single transaction. I want to show
# where I got to with unit testing TDD, before a final refactor to make it easier
# to see my workflow. Plus this is just a DB setup script for now.

# Made to demo TDD practices, as discussed in stage one interview

require 'pg'

class PgSetup
  class << self
    def build_db(database:'swarmtest')
      conn = PG.connect(dbname: 'postgres')
      conn.exec(%(CREATE DATABASE #{database}))
      conn.close
    end

    def build_table(table: 'ingestion_test', database: 'swarmtest', columns: [])
      cols = columns.map do |col|
        column = col.join(' ')
        column += ' default current_timestamp' if col[1] == 'timestamp'
        column
      end.join(', ')
      conn = PG.connect(dbname: database)
      conn.exec(%(CREATE TABLE #{table}\( #{cols} \)))
      conn.close
    end

    # dangerous if you have an update_ts_column as a function in your DB
    # outside scope of test but left in for completeness of DB use
    # not used in actual build script
    def insert_ts_function(database:'swarmtest')
      conn = PG.connect(dbname: database)
      conn.exec(%(CREATE OR REPLACE FUNCTION update_ts_column()
                  RETURNS TRIGGER AS $$
                  BEGIN
                    NEW.ts = now();
                    RETURN NEW;
                  END;
                  $$ language 'plpgsql';))
      conn.close
    end

    # dangerous if you have update_timestamp as a trigger in your DB
    # outside scope of test but left in for completeness of DB use
    # not used in actual build script, I was just having fun.
    def insert_ts_trigger(database:'swarmtest', table:'ingestion_test')
      conn = PG.connect(dbname: database)
      conn.exec(%(CREATE TRIGGER update_timestamp
                  BEFORE UPDATE ON #{table}
                  FOR EACH ROW EXECUTE PROCEDURE update_ts_column()))
      conn.close
    end

    # used for testing only
    def tear_down(database:'swarmtest')
      conn = PG.connect(dbname: 'postgres')
      kill_connections(conn, database) if has_connections(conn, database)
      conn.exec(%(DROP DATABASE #{database}))
      conn.close
    end

    private

    def kill_connections(conn, database)
      conn.exec(%(SELECT pg_terminate_backend(pg_stat_activity.pid)
                  FROM pg_stat_activity
                  WHERE pg_stat_activity.datname = '#{database}'
                  AND pid <> pg_backend_pid())
      )
    end

    def has_connections(conn, database)
      conn.exec(%(SELECT act.datname FROM pg_stat_activity as act))
          .values.include? [database]
    end

  end
end
