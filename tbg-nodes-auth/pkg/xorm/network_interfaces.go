package xorm

import (
	"context"
	stdsql "database/sql"

	"github.com/unchain/tbg-nodes-auth/pkg/3p/sql"

	"github.com/unchain/tbg-nodes-auth/gen/orm"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

func GetNetworkInterface(db *sql.DB, networkUUID string, protocol string) (network *orm.PublicEthereumNetwork, ni *orm.PublicEthereumNetworksNetworkExternalInterface, creds []*BasicAuthCreds, err error) {
	err = db.WrapTx(func(ctx context.Context, tx *stdsql.Tx) error {
		var err error

		network, err = orm.PublicEthereumNetworks(orm.PublicEthereumNetworkWhere.UUID.EQ(networkUUID)).One(ctx, tx)
		if err != nil {
			return err
		}

		ni, err = network.NetworkPublicEthereumNetworksNetworkExternalInterfaces(
			orm.PublicEthereumNetworksNetworkExternalInterfaceWhere.Protocol.EQ(null.NewString(protocol, true)),
			qm.Load(orm.PublicEthereumNetworksNetworkExternalInterfaceRels.PublicEthereumNetworksBasicauthCreds),
		).One(ctx, tx)
		if err != nil {
			return err
		}

		for _, cred := range ni.R.PublicEthereumNetworksBasicauthCreds {
			creds = append(creds, &BasicAuthCreds{
				Username: cred.Username.String,
				Password: cred.Password.String,
			})
		}

		return nil
	})

	if err != nil {
		return nil, nil, nil, err
	}

	return network, ni, creds, nil
}

type BasicAuthCreds struct {
	Username string
	Password string
}
