package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/iwanbk/gosqlbencher/plan"
	"github.com/pkg/profile"
)

func main() {
	var (
		planFile string
		profMode string // profiling mode
		profDir  string // profiling directory
	)

	flag.StringVar(&planFile, "plan", "plan.yaml", "gosqlbencher plan file")
	flag.StringVar(&profMode, "prof-mode", "", "profiling mode:block,cpu, mem, mutex")
	flag.StringVar(&profDir, "prof-dir", "prof", "dir where profiling files are written")
	flag.Parse()

	// profiling
	if profMode != "" {
		options, err := getProfOpts(profMode, profDir)
		if err != nil {
			log.Fatalf("failed to get prof options: %v", err)
		}
		defer profile.Start(options...).Stop()
	}

	// read config
	pl, err := plan.Read(planFile)
	if err != nil {
		log.Fatalf("failed to read plan: %v", err)
	}

	// print some nice message
	log.Printf("Benchmarking\ndsn: %v\nNumWorker:%v\nmax open conns:%v\n\n",
		pl.DataSourceName, pl.NumWorker, pl.MaxOpenConns)

	// open DB
	db, err := openDB(pl.DriverName, pl.DataSourceName)
	if err != nil {
		log.Fatalf("failed to open database:%v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(pl.MaxOpenConns)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// run the benchmarks
	for i, query := range pl.Queries {
		err = benchmarQuery(ctx, db, pl.NumWorker, query)
		if err != nil {
			log.Fatalf("benchmarck query #%v failed: %v", i, err)
		}
	}
}

func getProfOpts(profMode, profDir string) ([]func(*profile.Profile), error) {
	if profDir == "" {
		return nil, errEmptyProfDir
	}

	var (
		options    []func(*profile.Profile)
		modeOption func(*profile.Profile)
	)

	switch profMode {
	case "block":
		modeOption = profile.BlockProfile
	case "cpu":
		modeOption = profile.CPUProfile
	case "mem":
		modeOption = profile.MemProfile
	case "mutex":
		modeOption = profile.MutexProfile
	default:
		return nil, fmt.Errorf("invalid profiling mode: %s", profMode)
	}
	options = append(options, modeOption, profile.ProfilePath(profDir))
	return options, nil
}

var (
	errEmptyProfDir = errors.New("empty profile directory")
)
