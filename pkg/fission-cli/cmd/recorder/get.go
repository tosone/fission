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
	"os"
	"text/tabwriter"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/fission/fission/pkg/controller/client"
	"github.com/fission/fission/pkg/fission-cli/cliwrapper/cli"
	"github.com/fission/fission/pkg/fission-cli/cmd"
)

type GetSubCommand struct {
	client *client.Client
	name   string
}

func Get(flags cli.Input) error {
	opts := GetSubCommand{
		client: cmd.GetServer(flags),
	}
	return opts.do(flags)
}

func (opts *GetSubCommand) do(flags cli.Input) error {
	err := opts.complete(flags)
	if err != nil {
		return err
	}
	return opts.run(flags)
}

func (opts *GetSubCommand) complete(flags cli.Input) error {
	opts.name = flags.String("name")

	if len(opts.name) <= 0 {
		return errors.New("need a recorder name, use --name")
	}
	return nil
}

func (opts *GetSubCommand) run(flags cli.Input) error {
	recorder, err := opts.client.RecorderGet(&metav1.ObjectMeta{
		Name:      opts.name,
		Namespace: "default",
	})
	if err != nil {
		return errors.Wrap(err, "error getting recorder")
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

	fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\n",
		"NAME", "ENABLED", "FUNCTION", "TRIGGERS", "RETENTION_POLICY", "EVICTION_POLICY")
	fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\n",
		recorder.Metadata.Name, recorder.Spec.Enabled, recorder.Spec.Function, recorder.Spec.Triggers, recorder.Spec.RetentionPolicy, recorder.Spec.EvictionPolicy)
	w.Flush()
	return nil
}