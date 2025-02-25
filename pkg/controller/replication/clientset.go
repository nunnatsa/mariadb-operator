package replication

import (
	"context"
	"errors"
	"fmt"

	mariadbv1alpha1 "github.com/mariadb-operator/mariadb-operator/api/v1alpha1"
	"github.com/mariadb-operator/mariadb-operator/pkg/refresolver"
	sqlClient "github.com/mariadb-operator/mariadb-operator/pkg/sql"
	sqlClientSet "github.com/mariadb-operator/mariadb-operator/pkg/sqlset"
)

type replicationClientSet struct {
	*sqlClientSet.ClientSet
}

func newReplicationClientSet(mariadb *mariadbv1alpha1.MariaDB, refResolver *refresolver.RefResolver) (*replicationClientSet, error) {
	if !mariadb.Replication().Enabled {
		return nil, errors.New("'mariadb.spec.replication' is required to create a replicationClientSet")
	}
	return &replicationClientSet{
		ClientSet: sqlClientSet.NewClientSet(mariadb, refResolver),
	}, nil
}

func (c *replicationClientSet) close() error {
	return c.Close()
}

func (c *replicationClientSet) clientForIndex(ctx context.Context, index int) (*sqlClient.Client, error) {
	return c.ClientForIndex(ctx, index)
}

func (c *replicationClientSet) currentPrimaryClient(ctx context.Context) (*sqlClient.Client, error) {
	if c.Mariadb.Status.CurrentPrimaryPodIndex == nil {
		return nil, errors.New("'status.currentPrimaryPodIndex' must be set")
	}
	client, err := c.ClientForIndex(ctx, *c.Mariadb.Status.CurrentPrimaryPodIndex)
	if err != nil {
		return nil, fmt.Errorf("error getting current primary client: %v", err)
	}
	return client, nil
}

func (c *replicationClientSet) newPrimaryClient(ctx context.Context) (*sqlClient.Client, error) {
	client, err := c.ClientForIndex(ctx, *c.Mariadb.Replication().Primary.PodIndex)
	if err != nil {
		return nil, fmt.Errorf("error getting new primary client: %v", err)
	}
	return client, nil
}
