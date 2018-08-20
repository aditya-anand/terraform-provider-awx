// File : resource_inventory.go
package awx

import (
	"fmt"
	"strconv"

	awxgo "github.com/Colstuwjx/awx-go"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceInventoryObject() *schema.Resource {
	return &schema.Resource{
		Create: resourceInventoryCreate,
		Read:   resourceInventoryRead,
		Delete: resourceInventoryDelete,
		Update: resourceInventoryUpdate,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"organization": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"kind": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"host_filter": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"variables": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceInventoryCreate(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.InventoriesService

	_, res, _ := awxService.ListInventories(map[string]string{"name": d.Get("name").(string)})
	if len(res.Results) >= 1 {
		return fmt.Errorf("Inventory %s with id %d already exists", res.Results[0].Name, res.Results[0].ID)
	}

	result, err := awxService.CreateInventory(map[string]interface{}{
		"name":         d.Get("name").(string),
		"organization": d.Get("organization").(string),
		"kind":         d.Get("kind").(string),
		"host_filter":  d.Get("host_filter").(string),
		"variables":    d.Get("variables").(string),
	}, map[string]string{})
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceInventoryRead(d, m)

}

func resourceInventoryUpdate(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.InventoriesService
	_, res, _ := awxService.ListInventories(map[string]string{"name": d.Get("name").(string)})
	if len(res.Results) >= 1 {
		return fmt.Errorf("Inventory %s with id %d already exists", res.Results[0].Name, res.Results[0].ID)
	}
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	_, err = awxService.UpdateInventory(id, map[string]interface{}{
		"name":         d.Get("name").(string),
		"organization": d.Get("organization").(string),
		"kind":         d.Get("kind").(string),
		"host_filter":  d.Get("host_filter").(string),
		"variables":    d.Get("variables").(string),
	}, nil)
	if err != nil {
		return err
	}

	return resourceInventoryRead(d, m)

}

func resourceInventoryRead(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.InventoriesService
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Inventory %d not found", id)
	}
	r, err := awxService.GetInventory(id, map[string]string{})
	if err != nil {
		return err
	}
	d = setInventoryResourceData(d, r)
	return nil
}

func resourceInventoryDelete(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.InventoriesService
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	if _, err := awxService.DeleteInventory(id); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func setInventoryResourceData(d *schema.ResourceData, r *awxgo.Inventory) *schema.ResourceData {
	d.Set("name", r.Name)
	d.Set("description", r.Description)
	d.Set("organization", strconv.Itoa(r.Organization))
	d.Set("variables", r.Variables)
	return d
}