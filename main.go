package main

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func main() {
	var yamlconf string
	flag.StringVar(&yamlconf, "config", "./config.yaml", "")
	flag.Parse()

	buf, err := ioutil.ReadFile(yamlconf)
	if err != nil {
		logrus.Fatalf("reading config.yaml got error: %v", err)
		os.Exit(1)
	}
	config := appConfig{
		CacheDir:  os.TempDir(),
		Kindlegen: "kindlegen",
	}
	if err := yaml.Unmarshal(buf, &config); err != nil {
		logrus.Fatal(err)
		os.Exit(1)
	}
	if config.Verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}
	srv := newServer(&config)
	srv.run()
	logrus.Info("bye-bye")
}
