// Package main is hadron binary entrypoing
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
)

const (
	flagNamePlan     = "plan"
	goCommandRunVerb = "run"
)

var errPlanFileNotFound = errors.New("plan file not found")

func main() {
	// Configure zerolog
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cmd := &cli.Command{
		Name:  "hadron",
		Usage: "Declarative Docker deployment tool",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "log-level",
				Aliases: []string{"l"},
				Value:   "info",
				Usage:   "Log level (debug, info, warn, error)",
			},
		},
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			// Set log level
			level, err := zerolog.ParseLevel(cmd.String("log-level"))
			if err != nil {
				return ctx, fmt.Errorf("invalid log level: %w", err)
			}
			zerolog.SetGlobalLevel(level)

			return ctx, nil
		},
		Commands: []*cli.Command{
			{
				Name:  "deploy",
				Usage: "Deploy a plan to remote hosts",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     flagNamePlan,
						Aliases:  []string{"p"},
						Required: true,
						Usage:    "Path to the deployment plan (Go file)",
					},
					&cli.BoolFlag{
						Name:  "dry-run",
						Usage: "Show what would be deployed without executing",
					},
				},
				Action: deploy,
			},
			{
				Name:  "destroy",
				Usage: "Destroy resources defined in a plan",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     flagNamePlan,
						Aliases:  []string{"p"},
						Required: true,
						Usage:    "Path to the deployment plan (Go file)",
					},
				},
				Action: destroy,
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal().Err(err).Msg("Command failed")
	}
}

func deploy(_ context.Context, cmd *cli.Command) error {
	planPath := cmd.String(flagNamePlan)
	dryRun := cmd.Bool("dry-run")

	// Determine if planPath is a directory or file
	stat, err := os.Stat(planPath)
	if err != nil {
		return fmt.Errorf("%w: %s", errPlanFileNotFound, planPath)
	}

	var planDir string

	var args []string

	if stat.IsDir() {
		// Directory: go run .
		planDir = planPath
		args = []string{goCommandRunVerb, "."}
	} else {
		// File: go run basename
		planDir = filepath.Dir(planPath)
		args = []string{goCommandRunVerb, filepath.Base(planPath)}
	}

	log.Info().Str("plan", planPath).Bool("dry-run", dryRun).Msg("Deploying plan")

	// Execute go run on the plan
	//nolint:gosec
	execCmd := exec.Command("go", args...)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	execCmd.Env = append(os.Environ(), fmt.Sprintf("HADRON_DRY_RUN=%t", dryRun))
	execCmd.Dir = planDir

	if err := execCmd.Run(); err != nil {
		return fmt.Errorf("failed to execute plan: %w", err)
	}

	return nil
}

func destroy(_ context.Context, cmd *cli.Command) error {
	planPath := cmd.String(flagNamePlan)

	// Determine if planPath is a directory or file
	stat, err := os.Stat(planPath)
	if err != nil {
		return fmt.Errorf("%w: %s", errPlanFileNotFound, planPath)
	}

	var planDir string

	var args []string

	if stat.IsDir() {
		// Directory: go run .
		planDir = planPath
		args = []string{goCommandRunVerb, "."}
	} else {
		// File: go run basename
		planDir = filepath.Dir(planPath)
		args = []string{goCommandRunVerb, filepath.Base(planPath)}
	}

	log.Info().Str("plan", planPath).Msg("Destroying resources")

	// Execute go run on the plan with destroy mode
	//nolint:gosec
	execCmd := exec.Command("go", args...)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	execCmd.Env = append(os.Environ(), "HADRON_DESTROY=true")
	execCmd.Dir = planDir

	if err := execCmd.Run(); err != nil {
		return fmt.Errorf("failed to execute plan: %w", err)
	}

	return nil
}
