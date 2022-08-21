// Package config
/*
   Copyright 2022 Jellyfish message broker
   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at
       http://www.apache.org/licenses/LICENSE-2.0
   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/
package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Addr   string   `yaml:"addr"`
	Slaves []string `yaml:"slaves"`
}

func New(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "open config by path")
	}

	cnf := Config{}
	err = yaml.NewDecoder(file).Decode(&cnf)
	if err != nil {
		return nil, errors.Wrap(
			yaml.NewDecoder(file).Decode(&cnf),
			"decode config information",
		)
	}

	return &cnf, nil
}
