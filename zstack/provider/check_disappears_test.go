// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"

	tfjson "github.com/hashicorp/terraform-json"
)

// stateResourceAtAddress looks up a resource in the Terraform state by its address.
func stateResourceAtAddress(state *tfjson.State, address string) (*tfjson.StateResource, error) {
	if state == nil || state.Values == nil || state.Values.RootModule == nil {
		return nil, fmt.Errorf("no state available")
	}
	for _, r := range state.Values.RootModule.Resources {
		if r.Address == address {
			return r, nil
		}
	}
	return nil, fmt.Errorf("not found in state: %s", address)
}

// disappearsCheck is a generic statecheck.StateCheck that deletes a resource
// via the ZStack SDK to simulate external deletion.
type disappearsCheck struct {
	resourceAddress string
	deleteFunc      func(cli *client.ZSClient, id string) error
}

func (d disappearsCheck) CheckState(ctx context.Context, req statecheck.CheckStateRequest, resp *statecheck.CheckStateResponse) {
	r, err := stateResourceAtAddress(req.State, d.resourceAddress)
	if err != nil {
		resp.Error = err
		return
	}

	id, _ := r.AttributeValues["uuid"].(string)
	if id == "" {
		id, _ = r.AttributeValues["id"].(string)
	}
	if id == "" {
		resp.Error = fmt.Errorf("no uuid or id found for %s", d.resourceAddress)
		return
	}

	cli := testAccClientLoggedIn()
	if err := d.deleteFunc(cli, id); err != nil {
		resp.Error = fmt.Errorf("failed to delete %s (%s): %w", d.resourceAddress, id, err)
	}
}

// stateCheckDisappears returns a StateCheck that deletes the resource externally.
func stateCheckDisappears(resourceAddress string, deleteFunc func(cli *client.ZSClient, id string) error) statecheck.StateCheck {
	return disappearsCheck{resourceAddress: resourceAddress, deleteFunc: deleteFunc}
}

func stateCheckZoneDisappears(resourceAddress string) statecheck.StateCheck {
	return stateCheckDisappears(resourceAddress, func(cli *client.ZSClient, id string) error {
		return cli.DeleteZone(id, param.DeleteModePermissive)
	})
}

func stateCheckAccountDisappears(resourceAddress string) statecheck.StateCheck {
	return stateCheckDisappears(resourceAddress, func(cli *client.ZSClient, id string) error {
		return cli.DeleteAccount(id, param.DeleteModePermissive)
	})
}

func stateCheckSshKeyPairDisappears(resourceAddress string) statecheck.StateCheck {
	return stateCheckDisappears(resourceAddress, func(cli *client.ZSClient, id string) error {
		return cli.DeleteSshKeyPair(id, param.DeleteModePermissive)
	})
}

func stateCheckAffinityGroupDisappears(resourceAddress string) statecheck.StateCheck {
	return stateCheckDisappears(resourceAddress, func(cli *client.ZSClient, id string) error {
		return cli.DeleteAffinityGroup(id, param.DeleteModePermissive)
	})
}

func stateCheckClusterDisappears(resourceAddress string) statecheck.StateCheck {
	return stateCheckDisappears(resourceAddress, func(cli *client.ZSClient, id string) error {
		return cli.DeleteCluster(id, param.DeleteModePermissive)
	})
}

func stateCheckIAM2ProjectDisappears(resourceAddress string) statecheck.StateCheck {
	return stateCheckDisappears(resourceAddress, func(cli *client.ZSClient, id string) error {
		return cli.DeleteIAM2Project(id, param.DeleteModePermissive)
	})
}
