package main

import (
  "fmt"
  "github.com/imdario/mergo"
  l "github.com/launchdarkly/foundation/logger"
  "github.com/launchdarkly/gcfg"
  "os"
)

var Version = "DEV"

type Config struct {

}

var Global Config

func init() {
  var configFile string
  var override Config
  mode := os.Getenv("$name;format="upper,snake,word"$_MODE")

  isDev := false

  if mode == "production" {
    configFile = "$name$.prod.conf"
  } else if mode == "staging" {
    configFile = "$name$.stg.conf"
  } else if mode == "dogfood" {
    configFile = "$name$.dog.conf"
  } else {
    configFile = "$name$.local.conf"
    isDev = true
  }

  err := gcfg.ReadFileInto(&Global, configFile)

  if err != nil {
    fmt.Println(err)
    panic("Unable to read configuration file from " + configFile)
  }

  if isDev {
    err = gcfg.ReadFileInto(&override, "$name$.override.conf")
    if err != nil {
      l.Warn.Printf("Failed to load override configuration file: %+v", err)
    } else {
      err = mergo.MergeWithOverwrite(&Global, override)
      if err != nil {
        l.Warn.Printf("Failed to merge override configuration: %+v", err)
      }
    }
  }
}

func main() {
	fmt.Println("Starting $name$ version " + Version)
}
