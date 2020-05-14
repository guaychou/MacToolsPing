package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	Server struct {
		Address string `yaml:"address"`
	}
}

func readConfig()(*Config,error){
	config := &Config{}
	Dir,err:=os.Executable()
	log.Print(Dir)
	if err != nil {
		log.Fatal(err)
	}
	configPath:=getConfigDirectory()+"/"+configFileName
	file,err:=os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}
	return config, nil
}

func getConfigDirectory()string{
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func writeConfigState(address string){
	var config Config
	config.Server.Address=address
	d, err := yaml.Marshal(&config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = ioutil.WriteFile(getConfigDirectory()+"/"+configFileName, d, 0644)
	if err != nil {
		log.Fatal(err)
	}
}