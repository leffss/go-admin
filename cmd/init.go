package cmd

import (
	_ "github.com/leffss/go-admin/pkg/logging"
	_ "github.com/leffss/go-admin/pkg/setting"

	"github.com/leffss/go-admin/models"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use: "init",
	Short: "init database",
	Run: func(cmd *cobra.Command, args []string) {
		models.InitDatabase(password)
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		// 初始化依赖包
		initPkg()
	},
}

func init()  {
	initCmd.Flags().StringVarP(&password, "password", "p", "", "password for admin (required)")
	_ = initCmd.MarkFlagRequired("password")
	rootCmd.AddCommand(initCmd)
}
