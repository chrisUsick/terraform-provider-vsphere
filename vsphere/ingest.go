package vsphere

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
)

// "github.com/hashicorp/terraform-provider-vsphere/internal/helpers/virtualmachine"

func Ingest() {
	p := Provider()

	diag := p.Configure(context.Background(), &terraform.ResourceConfig{
		Config: map[string]interface{}{
			"vcenter_server":       "lv1vcenter01.dm.nfl.com",
			"allow_unverified_ssl": true,
			"user":                 os.Getenv("vsphere_user"),
			"password":             os.Getenv("vsphere_password"),
		},
	})
	if diag.HasError() {
		log.Fatalf("Error configuring provider %v", diag)
	}

	var client *Client
	client = p.Meta().(*Client)
	fmt.Printf("Client: %v\n", client)

	dcs, err := findDatacenters(client.vimClient)
	if err != nil {
		log.Fatalf("Failed to load datacenters %v", err)
	}

	for _, dc := range dcs {
		ingestDatacenter(client, p.ResourcesMap["vsphere_datacenter"], dc)
	}

}

func findDatacenters(client *govmomi.Client) ([]*object.Datacenter, error) {
	ctx := context.TODO()
	finder := find.NewFinder(client.Client, false)
	return finder.DatacenterList(ctx, "/*")
}

func ingestDatacenter(client *Client, resource *schema.Resource, dc *object.Datacenter) {
	resourceData := resource.TestResourceData()
	resourceData.SetId(dc.Name())
	err := resource.Read(resourceData, client)
	if err != nil {
		log.Fatalf("Error reading datacenter properties for `%s`. %v", dc.Name(), err)
	}
	fmt.Printf("Got DC\n %+v", resourceData)
}
