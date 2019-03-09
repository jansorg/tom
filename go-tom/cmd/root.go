package cmd

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/text/message"

	_config "github.com/jansorg/tom/go-tom/cmd/config"
	"github.com/jansorg/tom/go-tom/cmd/edit"
	"github.com/jansorg/tom/go-tom/cmd/frames"
	"github.com/jansorg/tom/go-tom/cmd/import"
	"github.com/jansorg/tom/go-tom/cmd/project"
	"github.com/jansorg/tom/go-tom/cmd/remove"
	"github.com/jansorg/tom/go-tom/cmd/report"
	"github.com/jansorg/tom/go-tom/cmd/status"
	"github.com/jansorg/tom/go-tom/config"
	"github.com/jansorg/tom/go-tom/context"
	"github.com/jansorg/tom/go-tom/i18n"
	"github.com/jansorg/tom/go-tom/query"
	"github.com/jansorg/tom/go-tom/store"
	"github.com/jansorg/tom/go-tom/storeHelper"
	"github.com/jansorg/tom/go-tom/util"
)

var ctx context.TomContext
var configFile string

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().String("data-dir", "", "data directory (default is $HOME/.tom)")
	RootCmd.PersistentFlags().String("backup-dir", "", "backup directory (default is $HOME/.tom/backup)")
	RootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file (default is $HOME/.tom/tom.yaml)")

	RootCmd.PersistentFlags().String("cpu-profile", "", "create a cpu profile for performance measurement")
	RootCmd.Flag("cpu-profile").Hidden = true
	RootCmd.Flags().String("mem-profile", "", "create a memory profile for performance measurement")
	RootCmd.Flag("mem-profile").Hidden = true

	project.NewCommand(&ctx, RootCmd)
	newTagsCommand(&ctx, RootCmd)
	frames.NewCommand(&ctx, RootCmd)
	newCreateCommand(&ctx, RootCmd)
	remove.NewCommand(&ctx, RootCmd)
	newRenameCommand(&ctx, RootCmd)
	newStartCommand(&ctx, RootCmd)
	newStopCommand(&ctx, RootCmd)
	newCancelCommand(&ctx, RootCmd)
	edit.NewEditCommand(&ctx, RootCmd)
	report.NewCommand(&ctx, RootCmd)
	imports.NewCommand(&ctx, RootCmd)
	status.NewCommand(&ctx, RootCmd)
	newInvoiceCommand(&ctx, RootCmd)
	_config.NewCommand(&ctx, RootCmd)
	// hidden command
	newCompletionCommand(&ctx, RootCmd)

	if err := viper.BindPFlag(config.KeyDataDir, RootCmd.PersistentFlags().Lookup("data-dir")); err != nil {
		util.Fatal(err)
	}
	if err := viper.BindPFlag(config.KeyBackupDir, RootCmd.PersistentFlags().Lookup("backup-dir")); err != nil {
		util.Fatal(err)
	}
}

func initConfig() {
	config.SetDefaults()
	if configFile != "" {
		viper.SetConfigFile(configFile)
	}

	// setup config dir if it doesn't exist
	dataDir := viper.GetString(config.KeyDataDir)
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dataDir, 0700); err != nil {
			util.Fatal(err)
		}
	}

	// setup backup dir if it doesn't exist
	backupDir := viper.GetString(config.KeyBackupDir)
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dataDir, 0700); err != nil {
			util.Fatal(err)
		}
	}

	_ = viper.ReadInConfig()

	dataStore, err := store.NewStore(dataDir, backupDir, viper.GetInt(config.KeyMaxBackups))
	if err != nil {
		util.Fatal(err)
	}

	ctx.Store = dataStore
	ctx.StoreHelper = storeHelper.NewStoreHelper(dataStore)
	ctx.Query = query.NewStoreQuery(dataStore)
	ctx.Language = i18n.FindPreferredLanguages()
	ctx.LocalePrinter = message.NewPrinter(ctx.Language)
	ctx.Locale = i18n.FindLocale(ctx.Language)
	ctx.DurationPrinter = i18n.NewDurationPrinter(ctx.Language)
	ctx.DecimalDurationPrinter = i18n.NewDecimalDurationPrinter(ctx.Language)
	ctx.DateTimePrinter = i18n.NewDateTimePrinter(ctx.Language)
}

const (
	//language=BASH
	bash_completion_func = `__gotime_projects_get()
{
    local -a projects
    readarray -t COMPREPLY < <(gotime projects 2>/dev/null | grep "$cur" | sed -e 's/ /\\ /g')
}

__gotime_get_projects()
{
    if [[ ${#nouns[@]} -eq 0 ]]; then
        __gotime_projects_get ""
	else
	    __gotime_projects_get ${nouns[${#nouns[@]} -1]}
    fi
    if [[ $? -eq 0 ]]; then
        return 0
    fi
}

__custom_func() {
    case ${last_command} in
        gotime_start)
            __gotime_get_projects
            return
            ;;
        *)
            ;;
    esac
}
`
)

var RootCmd = &cobra.Command{
	Use:                    "tom",
	Short:                  "tom is a command line application to track time.",
	Version:                "unknown",
	BashCompletionFunction: bash_completion_func,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cpuProfile, _ := cmd.Flags().GetString("cpu-profile")
		if cpuProfile != "" {
			log.Println("creating a cpu profile...")
			f, err := os.Create(cpuProfile)
			if err != nil {
				log.Fatal("could not create CPU profile: ", err)
			}
			if err := pprof.StartCPUProfile(f); err != nil {
				log.Fatal("could not start CPU profile: ", err)
			}
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		cpuProfile, _ := cmd.Flags().GetString("cpu-profile")
		if cpuProfile != "" {
			pprof.StopCPUProfile()
		}

		memProfile, _ := cmd.Flags().GetString("mem-profile")
		if memProfile != "" {
			log.Println("creating a mem profile...")
			f, err := os.Create(memProfile)
			if err != nil {
				log.Fatal("could not create memory profile: ", err)
			}
			runtime.GC() // get up-to-date statistics
			if err := pprof.WriteHeapProfile(f); err != nil {
				log.Fatal("could not write memory profile: ", err)
			}
			f.Close()
		}
	},
}

func Execute(version, commit, date string) {
	RootCmd.Version = version

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
