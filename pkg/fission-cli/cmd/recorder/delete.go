/*
Copyright 2019 The Fission Authors.

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

package recorder

import (
	"fmt"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/fission/fission/pkg/controller/client"
	"github.com/fission/fission/pkg/fission-cli/cliwrapper/cli"
	"github.com/fission/fission/pkg/fission-cli/cmd"
)

type DeleteSubCommand struct {
	client   *client.Client
	metadata *metav1.ObjectMeta
}

func Delete(flags cli.Input) error {
	opts := DeleteSubCommand{
		client: cmd.GetServer(flags),
	}
	return opts.do(flags)
}

func (opts *DeleteSubCommand) do(flags cli.Input) error {
	err := opts.complete(flags)
	if err != nil {
		return err
	}
	return opts.run(flags)
}

func (opts *DeleteSubCommand) complete(flags cli.Input) error {
	m, err := cmd.GetMetadata("name", "recorderns", flags)
	if err != nil {
		return err
	}
	opts.metadata = m
	return nil
}

func (opts *DeleteSubCommand) run(flags cli.Input) error {
	err := opts.client.RecorderDelete(opts.metadata)
	if err != nil {
		return errors.Wrap(err, "error deleting recorder")
	}
	fmt.Printf("recorder '%v' deleted\n", opts.metadata.Name)
	return nil
}