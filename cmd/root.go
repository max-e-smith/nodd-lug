package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

var verbose bool
var check bool
var background bool
var parallel int
var source string
var manifest string
var direct string
var try bool

var RootCmd = &cobra.Command{
	Use:   "clug",
	Short: "A domain driven retrieval tool for cruise-based datasets",
	Long: `A CLI library for downloading ocean data using domain driven criteria using
	direct access options such as the NOAA Open Data Dissemination (NODD) cloud.

	mb, will download all multibeam bathymetry data files when given a survey name argument(s), 
	path (prefix), or file manifest.

	csb, will download all crowdsourced bathymetry data files when given a survey name argument(s), 
	path (prefix), or file manifest.

	wcd, will download all water column data files when given a survey name argument(s), 
	path (prefix), or file manifest.

	help, provides usage information for each subcommand.
	`,
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Display more verbose output in console output. (default: false)")
	RootCmd.PersistentFlags().IntVarP(&parallel, "parallel", "p", 3, "Number of parallel downloads. (default: 3, max: 100)")
	RootCmd.PersistentFlags().BoolVarP(&check, "check", "c", false, "Check local disk space before downloading. (default: true)")
	RootCmd.PersistentFlags().BoolVarP(&background, "background", "b", false, "Run in background mode. (default: false)")
	RootCmd.PersistentFlags().StringVarP(&source, "source", "s", "", "Define direct data access source. (default: NODD)")
	RootCmd.PersistentFlags().StringVarP(&manifest, "manifest", "m", "", "Direct file download by providing a valid manifest. (default: none")
	RootCmd.PersistentFlags().StringVarP(&direct, "direct", "d", "", "Direct file download by providing a valid direct path. (default: none)")
	RootCmd.PersistentFlags().BoolVarP(&try, "try", "t", false, "Perform a dry run of command without actually downloading anything. (default: false)")

	vErr := viper.BindPFlag("verbose", RootCmd.PersistentFlags().Lookup("verbose"))
	if vErr != nil {
		log.Fatal(vErr)
	}

	pErr := viper.BindPFlag("parallel", RootCmd.PersistentFlags().Lookup("parallel"))
	if pErr != nil {
		log.Fatal(pErr)
	}

	cErr := viper.BindPFlag("check", RootCmd.PersistentFlags().Lookup("check"))
	if cErr != nil {
		log.Fatal(cErr)
	}

	bErr := viper.BindPFlag("background", RootCmd.PersistentFlags().Lookup("background"))
	if bErr != nil {
		log.Fatal(bErr)
	}

	sErr := viper.BindPFlag("source", RootCmd.PersistentFlags().Lookup("source"))
	if sErr != nil {
		log.Fatal(sErr)
	}

	mErr := viper.BindPFlag("manifest", RootCmd.PersistentFlags().Lookup("manifest"))
	if mErr != nil {
		log.Fatal(mErr)
	}

	dErr := viper.BindPFlag("direct", RootCmd.PersistentFlags().Lookup("direct"))
	if dErr != nil {
		log.Fatal(dErr)
	}

	tErr := viper.BindPFlag("try", RootCmd.PersistentFlags().Lookup("try"))
	if tErr != nil {
		log.Fatal(tErr)
	}

}
