/*
Copyright 2019 The Crossplane Authors.

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

package core

import (
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"gopkg.in/alecthomas/kingpin.v2"
	crds "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane/apis"
	"github.com/crossplane/crossplane/pkg/controller/apiextensions"
	"github.com/crossplane/crossplane/pkg/controller/pkg"
	"github.com/crossplane/crossplane/pkg/xpkg"
)

// Command configuration for the core Crossplane controllers.
type Command struct {
	Name      string
	Namespace string
	CacheDir  string
	Sync      time.Duration
}

// FromKingpin produces the core Crossplane command from a Kingpin command.
func FromKingpin(cmd *kingpin.CmdClause) *Command {
	c := &Command{Name: cmd.FullCommand()}
	cmd.Flag("namespace", "Namespace used to unpack and run packages.").Short('n').Default("crossplane-system").OverrideDefaultFromEnvar("POD_NAMESPACE").StringVar(&c.Namespace)
	cmd.Flag("cache-dir", "Directory used for caching package images.").Short('c').Default("/cache").OverrideDefaultFromEnvar("CACHE_DIR").ExistingDirVar(&c.CacheDir)
	cmd.Flag("sync", "Controller manager sync period duration such as 300ms, 1.5h or 2h45m").Short('s').Default("1h").DurationVar(&c.Sync)
	return c
}

// Run core Crossplane controllers.
func (c *Command) Run(log logging.Logger) error {
	log.Debug("Starting", "sync-period", c.Sync.String())

	cfg, err := ctrl.GetConfig()
	if err != nil {
		return errors.Wrap(err, "Cannot get config")
	}

	mgr, err := ctrl.NewManager(cfg, ctrl.Options{SyncPeriod: &c.Sync})
	if err != nil {
		return errors.Wrap(err, "Cannot create manager")
	}

	if err := crds.AddToScheme(mgr.GetScheme()); err != nil {
		return errors.Wrap(err, "Cannot add CustomResourceDefinition API to scheme")
	}

	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		return errors.Wrap(err, "Cannot add core Crossplane APIs to scheme")
	}

	if err := apiextensions.Setup(mgr, log); err != nil {
		return errors.Wrap(err, "Cannot setup API extension controllers")
	}

	pkgCache := xpkg.NewImageCache(c.CacheDir, afero.NewOsFs())

	if err := pkg.Setup(mgr, log, pkgCache, c.Namespace); err != nil {
		return errors.Wrap(err, "Cannot add packages controllers to manager")
	}

	return errors.Wrap(mgr.Start(ctrl.SetupSignalHandler()), "Cannot start controller manager")
}
