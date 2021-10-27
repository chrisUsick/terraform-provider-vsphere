package vsphere

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
)

const (
	objectTypeVirtualMachine string = "virtualMachine"
	objectTypeHost           string = "host"
)

func dataSourceVSphereVirtualMachines() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVsphereVirtualMachinesRead,

		Schema: map[string]*schema.Schema{
			"datacenter": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "name of datacenter. This can be its name or path",
			},
			"paths": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "list of absolute paths to VMs",
				Elem:        &schema.Schema{Type: schema.TypeList},
			},
		},
	}
}

func dataSourceVsphereVirtualMachinesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client).vimClient
	datacenterName := d.Get("datacenter").(string)
	paths, err := datacenterVMs(client, datacenterName)
	if err != nil {
		return err
	}
	d.Set("paths", paths)
	d.SetId(fmt.Sprintf("%s_vms", datacenterName))
	return nil
}

func datacenterVMs(c *govmomi.Client, name string) ([]string, error) {
	finder := find.NewFinder(c.Client, true)
	dc, err := finder.Datacenter(context.TODO(), name)
	if err != nil {
		return nil, err
	}
	vms, err := finder.VirtualMachineList(context.TODO(), fmt.Sprintf("%v/...", dc.InventoryPath))
	if err != nil {
		return nil, err
	}
	out := make([]string, len(vms))
	for i, v := range vms {
		out[i] = v.InventoryPath
	}
	return out, nil
}
