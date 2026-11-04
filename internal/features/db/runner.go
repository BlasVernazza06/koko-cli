package db

import (
	"context"
	"fmt"
	"io"
	"os"
)

// Runner orquesta la detección, resolución y ejecución de operaciones DB.
type Runner struct {
	Executor Executor
}

// NewRunner crea un nuevo Runner configurado con el OSExecutor por defecto.
func NewRunner() *Runner {
	return &Runner{
		Executor: NewOSExecutor(),
	}
}

// Run ejecuta el flujo completo de resolución y despacho de la acción DB.
func (r *Runner) Run(ctx context.Context, opts Options) error {
	targetDir := opts.Dir
	if targetDir == "" {
		targetDir = "."
	}

	pCtx, err := DetectProjectDBContext(targetDir)
	if err != nil {
		return err
	}

	resCmd, err := ResolveCommand(pCtx, opts.Action, opts.Args)
	if err != nil {
		return err
	}

	stdout := opts.Stdout
	if stdout == nil {
		stdout = os.Stdout
	}
	stderr := opts.Stderr
	if stderr == nil {
		stderr = os.Stderr
	}

	cmdSummary := FormatCommandSummary(resCmd)

	if opts.DryRun {
		fmt.Fprintf(stdout, "\n\033[90m┌\033[0m  \033[1mKoko DB · Dry Run Analysis\033[0m\n")
		fmt.Fprintf(stdout, "\033[90m│\033[0m\n")
		fmt.Fprintf(stdout, "\033[90m│\033[0m  \033[38;2;0;255;127m◆\033[0m  \033[1mTarget Dir\033[0m        \033[90m·\033[0m  \033[38;2;167;139;250m%s\033[0m\n", pCtx.ProjectDir)
		fmt.Fprintf(stdout, "\033[90m│\033[0m  \033[38;2;0;255;127m◆\033[0m  \033[1mDatabase Provider\033[0m \033[90m·\033[0m  \033[38;2;167;139;250m%s\033[0m\n", pCtx.Provider)
		fmt.Fprintf(stdout, "\033[90m│\033[0m  \033[38;2;0;255;127m◆\033[0m  \033[1mORM Engine\033[0m        \033[90m·\033[0m  \033[38;2;167;139;250m%s\033[0m\n", pCtx.ORM)
		fmt.Fprintf(stdout, "\033[90m│\033[0m  \033[38;2;0;255;127m◆\033[0m  \033[1mLayout\033[0m            \033[90m·\033[0m  \033[38;2;167;139;250m%s\033[0m\n", pCtx.Layout)
		fmt.Fprintf(stdout, "\033[90m│\033[0m  \033[38;2;0;255;127m◆\033[0m  \033[1mPackage Manager\033[0m   \033[90m·\033[0m  \033[38;2;167;139;250m%s\033[0m\n", pCtx.PackageManager)
		fmt.Fprintf(stdout, "\033[90m│\033[0m  \033[38;2;0;255;127m◆\033[0m  \033[1mResolution\033[0m        \033[90m·\033[0m  \033[90m%s\033[0m\n", resCmd.SourceDesc)
		fmt.Fprintf(stdout, "\033[90m│\033[0m\n")
		fmt.Fprintf(stdout, "\033[90m│\033[0m  \033[1mResolved Command:\033[0m \033[38;2;255;184;108m%s\033[0m (in %s)\n", cmdSummary, resCmd.WorkDir)
		fmt.Fprintf(stdout, "\033[90m│\033[0m\n")
		fmt.Fprintf(stdout, "\033[90m└\033[0m  \033[33mDry run completed: no processes were executed.\033[0m\n\n")
		return nil
	}

	fmt.Fprintf(stdout, "\n\033[90m┌\033[0m  \033[1mKoko DB · %s\033[0m \033[90m(%s · %s)\033[0m\n", opts.Action, pCtx.ORM, pCtx.Provider)
	fmt.Fprintf(stdout, "\033[90m│\033[0m  Executing: \033[38;2;167;139;250m%s\033[0m \033[90m[%s]\033[0m\n", cmdSummary, resCmd.SourceDesc)
	fmt.Fprintf(stdout, "\033[90m└\033[0m\n\n")

	executor := r.Executor
	if executor == nil {
		executor = NewOSExecutor()
	}

	return executor.Execute(ctx, *resCmd, opts.Stdin, stdout, stderr)
}

// PrintDBStackOverview imprime un resumen visual del stack de base de datos cuando se invoca `koko db` sin subcomando.
func PrintDBStackOverview(w io.Writer, targetDir string) error {
	pCtx, err := DetectProjectDBContext(targetDir)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "\n\033[90m┌\033[0m  \033[1mKoko DB · Database Stack Overview\033[0m\n")
	fmt.Fprintf(w, "\033[90m│\033[0m\n")
	fmt.Fprintf(w, "\033[90m│\033[0m  \033[38;2;0;255;127m◆\033[0m  \033[1mTarget Dir\033[0m        \033[90m·\033[0m  \033[38;2;167;139;250m%s\033[0m\n", pCtx.ProjectDir)
	if pCtx.ConfigExists {
		fmt.Fprintf(w, "\033[90m│\033[0m  \033[38;2;0;255;127m◆\033[0m  \033[1mManifest\033[0m          \033[90m·\033[0m  \033[38;2;167;139;250mkoko.config.json detected\033[0m\n")
	} else {
		fmt.Fprintf(w, "\033[90m│\033[0m  \033[38;2;255;184;108m!\033[0m  \033[1mManifest\033[0m          \033[90m·\033[0m  \033[38;2;255;184;108mInferred from files / package.json\033[0m\n")
	}
	fmt.Fprintf(w, "\033[90m│\033[0m  \033[38;2;0;255;127m◆\033[0m  \033[1mLayout\033[0m            \033[90m·\033[0m  \033[38;2;167;139;250m%s\033[0m\n", pCtx.Layout)
	fmt.Fprintf(w, "\033[90m│\033[0m  \033[38;2;0;255;127m◆\033[0m  \033[1mPackage Manager\033[0m   \033[90m·\033[0m  \033[38;2;167;139;250m%s\033[0m\n", pCtx.PackageManager)
	fmt.Fprintf(w, "\033[90m│\033[0m  \033[38;2;0;255;127m◆\033[0m  \033[1mDatabase Provider\033[0m \033[90m·\033[0m  \033[38;2;167;139;250m%s\033[0m\n", defaultStr(pCtx.Provider, "none"))
	fmt.Fprintf(w, "\033[90m│\033[0m  \033[38;2;0;255;127m◆\033[0m  \033[1mORM / Driver\033[0m      \033[90m·\033[0m  \033[38;2;167;139;250m%s\033[0m\n", defaultStr(pCtx.ORM, "none"))
	fmt.Fprintf(w, "\033[90m│\033[0m\n")
	fmt.Fprintf(w, "\033[90m│\033[0m  \033[1mAvailable Subcommands:\033[0m\n")
	fmt.Fprintf(w, "\033[90m│\033[0m    \033[32m•\033[0m \033[1mmigrate\033[0m   \033[90mRun pending database schema migrations\033[0m\n")
	fmt.Fprintf(w, "\033[90m│\033[0m    \033[32m•\033[0m \033[1mpush\033[0m      \033[90mPush schema state directly to database\033[0m\n")
	fmt.Fprintf(w, "\033[90m│\033[0m    \033[32m•\033[0m \033[1mgenerate\033[0m  \033[90mGenerate ORM client and type definitions\033[0m\n")
	fmt.Fprintf(w, "\033[90m│\033[0m    \033[32m•\033[0m \033[1mstudio\033[0m    \033[90mOpen interactive database visual UI / studio\033[0m\n")
	fmt.Fprintf(w, "\033[90m│\033[0m    \033[32m•\033[0m \033[1mseed\033[0m      \033[90mSeed database with test data\033[0m\n")
	fmt.Fprintf(w, "\033[90m│\033[0m    \033[32m•\033[0m \033[1mreset\033[0m     \033[90mReset database schema and drop tables\033[0m\n")
	fmt.Fprintf(w, "\033[90m│\033[0m\n")
	fmt.Fprintf(w, "\033[90m└\033[0m  \033[90mUse 'koko db <command> --help' for command-specific flags.\033[0m\n\n")

	return nil
}

func defaultStr(s string, def string) string {
	if s == "" {
		return def
	}
	return s
}
