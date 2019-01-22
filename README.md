# gosqlbencher

[![Build Status](https://travis-ci.org/iwanbk/gosqlbencher.svg?branch=master)](https://travis-ci.org/iwanbk/gosqlbencher)
[![Go Report Card](https://goreportcard.com/badge/github.com/iwanbk/gosqlbencher)](https://goreportcard.com/report/github.com/iwanbk/gosqlbencher)
[![codecov](https://codecov.io/gh/iwanbk/gosqlbencher/branch/master/graph/badge.svg)](https://codecov.io/gh/iwanbk/gosqlbencher)

Tool to benchmark your Go SQL ecosystem: query, driver, and the server.

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
- [x] github.com/mattn/go-sqlite3
- [x] github.com/ziutek/mymysql
- [ ] github.com/go-sql-driver/mysql

#### Database Settings

- see how [SetMaxOpenConns](golang.org/pkg/database/sql/#DB.SetMaxOpenConns) affect the performance

#### Number of concurrent goroutines/requests

`num_worker` option set the number of goroutines which execute the queries.
In web application, it simulates the number of concurrent requests you have at any given time.

#### Go Profiling Data

It could generate Go profiling data with `prof-mode` option, so we could have better insight about how SQL package work under various scenarios.

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

## Configuration

The best docs for the configuration right now is the `godoc` page for [`Plan`](https://godoc.org/github.com/iwanbk/gosqlbencher/plan#Plan).

There are also benchmark plan examples in [examples](./examples) directory.

## TODO

- add support for sqlite
- plan file validation