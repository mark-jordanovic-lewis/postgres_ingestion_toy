- Run go build in src directory, this builds the binaries to run the project.
# Optimise injection of data into postgres DB
Author : Mark Jordanovic-Lewis

## WARNING
Setting up this project may effect unwanted changes in your postgres DB, please use a fresh DB or ensure that the function which will be applied to the DB to  update the ts column does not interfere with any other procedures or triggers you have already setup.

You have been warned.


### Setup

- If you don't have go use asdf to get it and setup your PATH approprately.
  (see - https://github.com/asdf-vm/asdf)
	$ asdf plugin-add golang https://github.com/kennyp/asdf-golang.git
	$ asdf install golang 1.8
- Run the setup script `setup.sh` which installs the required go packages.
	$ ./setup.sh
- Run `workspace` which sets go to use the current dir as your workspace.
	$ ./workspace
- If you don't have ruby install it with asdf.
	$ asdf plugin-add ruby https://github.com/asdf-vm/asdf-ruby.git
	$ asdf install ruby 2.4.0
	$ asdf local ruby 2.4.0
- Install necessary gems
	$ gem install bundler
	$ bundle install
- Run the ruby script to insert a function and trigger into the postgres DB:
	$ ruby build_DB.rb
- Run `go build` in src directory, this builds the binaries to run the project.
	$ cd src
	$ go build

(u cud not has bourne again shell so maybe change /bin/bash in above scripts to /bin/zsh)

### Test
#### Go
Run go test in src directory (all the source is in there).
#### Ruby
Run rspec in the top level of the project dir


### Running the project

main <db_uname> <db_name> <db_pwd> <db_tablename>


# Task

*) Generate postgres DB with single table `swarmtest`
	Rows: ts    - timestamp (measurement of row copy time)
	      src   - bigint
	      dst   - bigint
	      flags - bigint
   Add update timestamp function to DB.

1) Generate M*3 random bigints (max M : 100,000)
2) Use COPY FROM to ingest data
	COPY `tablename` [ ( column [, ...] ) ]
		FROM { 'filename' }
3) Perform this N times (max N : 10,000)
4) Read ts column and calculate rows/s for benchmarking

	done at this point - make faster

5) Lock row after updating index and benchmark
6) Make concurrent, ensure safety and benchmark

	really done here - containerise

7) Make into a docker container that pulls this from github
   installs everything and builds the DB and project.
