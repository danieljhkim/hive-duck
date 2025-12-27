package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/danieljhkim/hive-duck/internal/config"
	"github.com/danieljhkim/hive-duck/internal/engine"
	"github.com/danieljhkim/hive-duck/internal/output"
	"github.com/danieljhkim/hive-duck/internal/preprocess"
)

func Execute() {
	if err := newRootCmd().Execute(); err != nil {
		os.Exit(2)
	}
}

func newRootCmd() *cobra.Command {
	var (
		expr              string
		file              string
		dbPath            string
		configPath        string
		extsCSV           string
		silent            bool
		strict            bool
		dryRun            bool
		failOnUnsupported bool
		outputFormat      string
		hiveconf          []string
		hivevar           []string
	)

	cmd := &cobra.Command{
		Use:   "hive-duck",
		Short: "Hive-compatible wrapper for -e/-f using DuckDB",
		RunE: func(cmd *cobra.Command, args []string) error {
			if (expr == "" && file == "") || (expr != "" && file != "") {
				return fmt.Errorf("exactly one of -e or -f must be provided")
			}

			// Parse output format
			outFmt, err := output.ParseFormat(outputFormat)
			if err != nil {
				return err
			}

			// Load database mapping config if provided
			var dbMap *config.DatabaseMap
			if configPath != "" {
				dbMap, err = config.LoadDatabaseMap(configPath)
				if err != nil {
					return fmt.Errorf("load config: %w", err)
				}
			}

			cfg, err := config.FromFlags(hiveconf, hivevar)
			if err != nil {
				return err
			}
			cfg.StrictVars = strict

			// Load SQL input
			var sqlText string
			if expr != "" {
				sqlText = expr
			} else {
				b, err := os.ReadFile(file)
				if err != nil {
					return fmt.Errorf("read -f file: %w", err)
				}
				sqlText = string(b)
			}

			// Substitute vars
			sqlText, err = preprocess.Substitute(sqlText, cfg)
			if err != nil {
				return err
			}

			// Split statements
			stmts, err := preprocess.SplitStatements(sqlText)
			if err != nil {
				return err
			}

			// Check for unsupported Hive statements
			unsupported := preprocess.DetectUnsupported(stmts)
			if len(unsupported) > 0 {
				// Print warnings to stderr
				for _, u := range unsupported {
					fmt.Fprintf(os.Stderr, "WARNING: Unsupported Hive statement: %s\n", u.Keyword)
					fmt.Fprintf(os.Stderr, "  Statement: %s\n", u.Statement)
					fmt.Fprintf(os.Stderr, "  Reason: %s\n\n", u.Reason)
				}

				if failOnUnsupported {
					return fmt.Errorf("found %d unsupported Hive statement(s); use --fail-on-unsupported=false to continue anyway", len(unsupported))
				}
			}

			// Rewrite Hive statements to DuckDB equivalents
			rewriteOpts := &preprocess.RewriteOptions{
				DatabaseMap: dbMap,
			}
			rewriteResult, err := preprocess.Rewrite(stmts, rewriteOpts)
			if err != nil {
				return err
			}

			// Dry-run mode: print rewritten SQL and exit
			if dryRun {
				for _, stmt := range rewriteResult.Statements {
					fmt.Println(stmt + ";")
				}
				return nil
			}

			// Connect + run
			exts := []string{}
			if strings.TrimSpace(extsCSV) != "" {
				for _, e := range strings.Split(extsCSV, ",") {
					e = strings.TrimSpace(e)
					if e != "" {
						exts = append(exts, e)
					}
				}
			}

			r := engine.Runner{
				DBPath:       dbPath,
				Exts:         exts,
				Silent:       silent,
				OutputFormat: outFmt,
				DatabaseMap:  dbMap,
			}
			return r.Run(rewriteResult.Statements)
		},
	}

	cmd.Flags().StringVarP(&expr, "execute", "e", "", "SQL string to execute (Hive: -e)")
	cmd.Flags().StringVarP(&file, "file", "f", "", "SQL file to execute (Hive: -f)")
	cmd.Flags().StringVar(&dbPath, "database", ":memory:", "DuckDB database path or :memory:")
	cmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to databases.yaml config file for DB mapping")
	cmd.Flags().StringVar(&extsCSV, "ext", "", "Comma-separated DuckDB extensions to INSTALL/LOAD (e.g. avro,httpfs,json)")
	cmd.Flags().BoolVarP(&silent, "silent", "S", false, "Suppress non-result output")
	cmd.Flags().BoolVar(&strict, "strict-vars", true, "Fail if a referenced hiveconf/hivevar/env var is missing")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print rewritten SQL without executing")
	cmd.Flags().BoolVar(&failOnUnsupported, "fail-on-unsupported", false, "Fail if unsupported Hive statements are detected")
	cmd.Flags().StringVarP(&outputFormat, "output", "o", "table", "Output format: table, csv, tsv, json")

	cmd.Flags().StringArrayVar(&hiveconf, "hiveconf", nil, "Hive conf var k=v (repeatable)")
	cmd.Flags().StringArrayVar(&hivevar, "hivevar", nil, "Hive var name=v (repeatable)")

	return cmd
}
