# gosqlbencher

[![Build Status](https://travis-ci.org/iwanbk/gosqlbencher.svg?branch=master)](https://travis-ci.org/iwanbk/gosqlbencher)
[![Go Report Card](https://goreportcard.com/badge/github.com/iwanbk/gosqlbencher)](https://goreportcard.com/report/github.com/iwanbk/gosqlbencher)
[![codecov](https://codecov.io/gh/iwanbk/gosqlbencher/branch/master/graph/badge.svg)](https://codecov.io/gh/iwanbk/gosqlbencher)

Tool to benchmark Go SQL ecosystem: query, [database/sql](https://golang.org/pkg/database/sql/)funcs usage, driver, and the server.

## Why using this tool 

Why using this tool instead of existing tools like pgbench and SysBench?

Go's [database/sql](https://golang.org/pkg/database/sql/) package has it's own behaviour that different with 
both `pgbench` and `SysBench`. It already has it's own connection pool and also has it's own way to manage
the database connection and prepared statements.

## Things to observe/benchmark

There are many things we could observe/benchmark with this tool

#### Database server

Play with your DB server settings and benchmark it againts your query.

#### Database query

##### Test various queries like exec, query, and query_row.

Supported DB queries:
- [x] Exec
- [x] ExecContext
- [x] Query
- [x] QueryContext
- [ ] QueryRow
- [ ] QueryRowContext
- [ ] Transaction

TODO : scan the result of `Query%` queries.

##### Test between using prepared statement, placeholder, and plain query.

Supported query modes:
- prepared statement:
  - create it once on the initialization
  - create it on each query
- using placeholder (e.g.: `$1` in postgresql)
- plain query (e.g. using Go's `%d` and `%s`)

#### Driver

Easily switch between supported drivers and see the difference.

Supported drivers:
- [x] github.com/lib/pq (driver_name : postgres)
- [x] github.com/jackc/pgx (driver_name : pgx)
- [x] github.com/ziutek/mymysql (driver name: mymysql)
- [x] github.com/go-sql-driver/mysql (driver name: mysql)
- [] github.com/mattn/go-sqlite3

#### Database Settings

- see how [SetMaxOpenConns](golang.org/pkg/database/sql/#DB.SetMaxOpenConns) affect the performance

#### Number of concurrent goroutines/requests

`num_worker` option set the number of goroutines which execute the queries.
In web application, it simulates the number of concurrent requests you have at any given time.

#### Go Profiling Data

It could generate Go profiling data with `prof-mode` option, so we could have better insight about how SQL package work under various scenarios.

## Configuration

The `gosqlbencher` configuration is called `plan` file. It is called `plan` instead of `config` because it is not a mere config, it contatins the benchmarking procedures to be executed by `gosqlbencher`.

The best docs for the configuration right now is the `godoc` page for [`Plan`](https://godoc.org/github.com/iwanbk/gosqlbencher/plan#Plan).

There are also benchmark plan examples in [examples](./examples) directory.

## How It Works

1. Query is benchmarked by allowing query to be specified in the plan file. It supports two mode on substituing the value in query:
- placholder: something like `$1` in postgresql or `?` in mysql is supported
- standard Go `fmt`: using `%d` for integer or `%s` for string

The value will then be generated `randomly` or `sequentially` on execution, it depends on the plan file.

2. [database/sql](https://golang.org/pkg/database/sql/) funcs usage is benchmarked by supporting below query execution:
- prepare on initialization: query will be executed in prepared statement, which will only be created on initialization
- prepare: prepared statement will be executed right before executing the query
- SQL placholder usage
- Plain query without SQL placeholder, possibly using Go formatting

3. Driver could be specified as well in the plan file, so we could easily see the performance between drivers.

See [driver](#driver) section to see the list of available drivers

4. Database Settings

Current database setting supported is only max open connections, which can be configured in the plan file.

## Quick Start

Build
```bash
go build -v
```

See the help

```bash
$ ./gosqlbencher -h
Usage of ./gosqlbencher:
  -plan string
        gosqlbencher plan file (default "plan.yaml")
  -prof-dir string
        dir where profiling files are written (default "prof")
  -prof-mode string
        profiling mode:block,cpu, mem, mutex

```

Start postgresql server using the provided docker compose
```
$ docker-compose -f examples/docker-compose-postgresql.yaml up -d --build
Creating network "examples_default" with the default driver
Creating postgres-gosqlbencher ... done
```

Create test table and execute `insert` test
```bash
$ ./gosqlbencher -plan=examples/insert.plan.postgre.yaml 
```

Execute various `select` tests
```bash
 ./gosqlbencher -plan=examples/query.plan.postgre.yaml 
```

Profiling the CPU for SQL query, with prepared statement created on each query
```bash
$ $ ./gosqlbencher -plan=examples/query.plan.postgre.profile.yaml -prof-mode=cpu -prof-dir=prof
```
It currently has limitation to only able to profile one query

Create pdf version of the profiling data
```
go tool pprof -pdf prof/cpu.pprof > ~/Desktop/cpuprof.pdf
```
Check the `pdf` in `~/Desktop/cpuprof.pd`

Try different scenario in the `plan` file and compare the profiling data & performance, for example:
- set different number of `max_open_conns`
- set `prepare_on_init` to `true`: it usually will improve the performance

## TODO

- add support for sqlite
- plan file validation