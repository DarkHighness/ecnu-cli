package net

import (
	"fmt"
	"github.com/DarkHighness/ecnu-cli/pkg/srun"
	"github.com/DarkHighness/ecnu-cli/pkg/table"
	"github.com/parnurzeal/gorequest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net/url"
	"time"
)

// Cmd represents the srun command
var Cmd = &cobra.Command{
	Use:   "net",
	Short: "i-hdu network auth cli",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		_, err := url.ParseRequestURI(viper.GetString("net.endpoint"))
		cobra.CheckErr(err)
		portalServer = srun.New(viper.GetString("net.endpoint"), viper.GetString("net.acid"))
	},
}

var portalServer *srun.PortalServer

func init() {
	Cmd.PersistentFlags().StringP("endpoint", "e", "", "endpoint host of srun")
	viper.SetDefault("net.endpoint", "https://login.ecnu.edu.cn")
	cobra.CheckErr(viper.BindPFlag("net.endpoint", Cmd.PersistentFlags().Lookup("endpoint")))

	Cmd.PersistentFlags().StringP("acid", "a", "", "ac_id of srun")
	resp, _, errs := gorequest.New().Get("http://www.baidu.com").End()
	if errs != nil {
		viper.SetDefault("net.acid", "0")
	} else if acid := resp.Request.URL.Query().Get("ac_id"); acid != "" {
		viper.SetDefault("net.acid", acid)
	} else {
		viper.SetDefault("net.acid", "1")
	}
	cobra.CheckErr(viper.BindPFlag("net.acid", Cmd.PersistentFlags().Lookup("acid")))

	loginCmd.Flags().StringP("username", "u", "", "username of srun")
	loginCmd.Flags().StringP("password", "p", "", "password of srun")
	loginCmd.Flags().BoolP("daemon", "d", false, "daemon mode")
	loginCmd.Flags().IntP("interval", "i", 60, "second interval of daemon mode")

	logoutCmd.Flags().StringP("username", "u", "", "username of srun")

	Cmd.AddCommand(infoCmd, loginCmd, logoutCmd, internetCmd)
}

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "show info of your i-hdu network",
	Run: func(cmd *cobra.Command, args []string) {
		info, err := portalServer.GetUserInfo()
		table.PrintStruct(info, "chinese")
		cobra.CheckErr(err)
	},
}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login i-hdu of the account",
	Run: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(viper.BindPFlag("net.auth.username", cmd.Flags().Lookup("username")))
		cobra.CheckErr(viper.BindPFlag("net.auth.password", cmd.Flags().Lookup("password")))

		cobra.CheckErr(portalServer.SetUsername(viper.GetString("net.auth.username")))
		cobra.CheckErr(portalServer.SetPassword(viper.GetString("net.auth.password")))
		info, err := portalServer.GetUserInfo()
		if viper.GetBool("verbose") {
			table.PrintStruct(info, "chinese")
		}
		cobra.CheckErr(err)

		challenge, err := portalServer.GetChallenge()
		if viper.GetBool("verbose") {
			table.PrintStruct(challenge, "chinese")
		}
		cobra.CheckErr(err)

		loginResponse, err := portalServer.PortalLogin()
		table.PrintStruct(loginResponse, "chinese")
		cobra.CheckErr(err)

		if v, err := cmd.Flags().GetBool("daemon"); err == nil && v {
			interval, _ := cmd.Flags().GetInt("interval")
			log.Printf("start daemon: check every %d seconds\n", interval)
			startTime := time.Now()
			for {
				//检测是否登录成功，如果登录过期则重新登录
				info, err := portalServer.GetUserInfo()
				cobra.CheckErr(err)
				if ok, _ := info.IsOK(); !ok {
					log.Println("check failed: start retry")
					_, err = portalServer.GetChallenge()
					cobra.CheckErr(err)
					_, err = portalServer.PortalLogin()
					cobra.CheckErr(err)

					startTime = time.Now()
				}
				log.Printf("check succed: live time %fs\n", time.Now().Sub(startTime).Seconds())

				//检测是否能访问互联网，如果不能，则退出登录并尝试重新登录
				if !portalServer.Internet() {
					log.Println("internet is not available, try to login again")
					if _, err := portalServer.PortalLogout(); err != nil {
						log.Println(err)
					}
				} else {
					time.Sleep(time.Second * time.Duration(interval))
				}
			}
		}

	},
}

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "logout i-hdu of the account",
	Run: func(cmd *cobra.Command, args []string) {
		info, err := portalServer.GetUserInfo()
		cobra.CheckErr(err)
		if ok, _ := info.IsOK(); !ok {
			fmt.Println("you are not login")
			return
		}
		if viper.GetBool("verbose") {
			table.PrintStruct(info, "chinese")
		}
		fmt.Printf("you are logout account %s\n", info.UserName)
		cobra.CheckErr(portalServer.SetUsername(info.UserName))
		logoutResponse, err := portalServer.PortalLogout()
		table.PrintStruct(logoutResponse, "chinese")
		cobra.CheckErr(err)
	},
}

// internetCmd represents the logout command
var internetCmd = &cobra.Command{
	Use:   "internet",
	Short: "check if connect to the internet",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(portalServer.Internet())
	},
}
