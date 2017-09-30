require_relative '../lib/pg_setup'
require 'pg'

RSpec.describe "postgresql DB build script for swarm64 codetest" do
  let(:cols) {[
    ['ts', 'timestamp'],
    ['src', 'bigint'],
    ['dst', 'bigint'],
    ['flags', 'bigint']
  ]}
  context "When creating the database", order: :defined do

    before(:all) { PgSetup::build_db }
    after(:all) { PgSetup::tear_down }

    it 'builds a DB with name "swarmtest"' do
      connection = PG.connect(dbname: 'postgres')
      dbs = connection.exec("SELECT datname FROM pg_database")
                      .map{|db| db['datname']}
      expect(dbs.include? 'swarmtest').to be true
      connection.close
    end

    it 'identifies connections to a database' do
      connection = PG.connect(dbname: 'postgres')
      expect(PgSetup.send(:has_connections, connection, 'postgres')).to be true
    end

    it 'agressively kills connections to a database' do
      pg_connection = PG.connect(dbname: 'postgres')
      st_connection = PG.connect(dbname: 'swarmtest')
      expect(PgSetup.send(:has_connections, pg_connection, 'swarmtest')).to be true
      PgSetup.send(:kill_connections, pg_connection, 'swarmtest')
      expect(PgSetup.send(:has_connections, pg_connection, 'swarmtest')).to be false
      pg_connection.close
      st_connection.close
    end

    it 'builds ingestion_test table' do
      PgSetup::build_table(columns: cols)
      connection = PG.connect(dbname: 'swarmtest')
      tables = connection.exec("SELECT table_name FROM information_schema.tables")
                         .map{|db| db['table_name']}
      expect(tables.include? 'ingestion_test').to be true
      connection.close
    end

    it 'builds the columns correctly' do
      connection = PG.connect(dbname: 'swarmtest')
      cols = connection.exec(%(SELECT column_name, data_type
                               FROM information_schema.columns
                               WHERE table_name='ingestion_test'))
                       .map{|col| [col['column_name'], col['data_type']] }
      expect(cols).to eq cols
      connection.close
    end

    it 'adds timestamping function correctly' do
      PgSetup::insert_ts_function
      connection = PG.connect(dbname: 'swarmtest')
      procs = connection.exec(%(SELECT p.proname FROM pg_proc as p))
      expect(procs.values.include? ["update_ts_column"] ).to be true
    end

    it 'inserts timestamping trigger correctly' do
      PgSetup::insert_ts_trigger
      connection = PG.connect(dbname: 'swarmtest')
      trigs = connection.exec(%(SELECT t.tgname FROM pg_trigger t))
      expect(trigs.values.include? ["update_timestamp"] ).to be true
    end
  end

  context 'after DB is built', context: :defined do
    before(:all) do
      PgSetup::build_db
      PgSetup::build_table(columns: [
        ['ts', 'timestamp'],
        ['src', 'bigint'],
        ['dst', 'bigint'],
        ['flags', 'bigint']
      ])
      PgSetup::insert_ts_function
      PgSetup::insert_ts_trigger
    end
    after(:all) { PgSetup::tear_down }

    it 'sets ts on insert of data fields' do
      connection = PG.connect(dbname: 'swarmtest')
      new_row = connection.exec(%(INSERT INTO ingestion_test (src, dst, flags)
                                  VALUES (1,2,3)
                                  RETURNING *)).values.flatten
      expect(new_row[1,3].map(&:to_i)).to eq [1,2,3]
      @timestamp = new_row[0]
      expect(new_row[0]).not_to be nil
    end

    # outside scope of test but left in for completeness of DB build
    it 'updates the timestamp on updating data fields' do
      connection = PG.connect(dbname: 'swarmtest')
      new_vals = connection.exec(%(UPDATE ingestion_test
                                   SET (src, dst, flags) = (4,5,6)
                                   WHERE src = 1 AND dst = 2 AND flags = 3
                                   RETURNING *)).values.flatten
      expect(new_vals[1,3].map(&:to_i)).to eq [4,5,6]
      expect(new_vals[0]).not_to eq @timestamp
    end

  end

end
