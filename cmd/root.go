/*
Copyright © 2021 muratgu <mgungora@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	//"os"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "voipms",
	Short: "Command line interface to your Voip.ms account",
	Long: `
voipms is a CLI tool for your Voip.ms account.
`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./app.env)")

	rootCmd.PersistentFlags().StringP("did", "", viper.GetString("VOIPMS_DID"), "DID Number")

	rootCmd.Version = "0.0.1"
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("app")
		viper.SetConfigType("env")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
}

func Post(method string, formData map[string]string) (*string, error) {
	uid := viper.GetString("VOIPMS_API_UID")
	if uid == "" {
		log.Fatal("VOIPMS_API_UID undefined")
	}

	pwd := viper.GetString("VOIPMS_API_PWD")
	if pwd == "" {
		log.Fatal("VOIPMS_API_PWD undefined")
	}
	if formData == nil {
		formData = map[string]string{}
	}
	formData["api_username"] = uid
	formData["api_password"] = pwd
	formData["method"] = method
	
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	for key, element := range formData {
		_ = writer.WriteField(key, element)
	}
	if err := writer.Close(); err != nil {
		fmt.Println(err)
		return nil, err
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", "https://voip.ms/api/v1/rest.php", payload)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	data := fmt.Sprintf("%s", body)

	return &data, nil
}

func IfSetElse(value bool, whenSet string, whenNotSet string) string {
	if value {
		return whenSet
	} else {
		return whenNotSet
	}
}

func Println(data *string, err error) {
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(*data)
	}
}
