package hetznerrobot

import (
	"context"
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"regexp"
)

func resourceSshKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSshKeyCreate,
		ReadContext:   resourceSshKeyRead,
		UpdateContext: resourceSshKeyUpdate,
		DeleteContext: resourceSshKeyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceSshKeyImportState,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Key name",
			},
			"data": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Key data in OpenSSH or SSH2 format",
				ForceNew:    true,
			},
			"fingerprint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Key fingerprint",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Key algorithm type",
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Key size in bits",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation date",
			},
		},
	}
}

func resourceSshKeyImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	c := meta.(HetznerRobotClient)

	keyFingerprint := d.Id()
	match, err := regexp.Match(`^[0-9a-f]{2}(:[0-9a-f]{2}){15}$`, []byte(keyFingerprint))
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, errors.New("invalid key fingerprint format")
	}

	key, err := c.getSshKey(ctx, keyFingerprint)
	if err != nil {
		return nil, err
	}

	d.Set("name", key.Name)
	d.Set("data", key.Data)
	d.Set("fingerprint", key.Fingerprint)
	d.Set("type", key.Type)
	d.Set("size", key.Size)
	d.Set("created_at", key.CreatedAt)

	return []*schema.ResourceData{d}, nil
}

func resourceSshKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)

	name := d.Get("name").(string)
	data := d.Get("data").(string)

	key, err := c.createSshKey(ctx, name, data)
	if err != nil {
		return diag.Errorf("Unable to create SSH key %q:\n\t %q", name, err)
	}

	d.Set("fingerprint", key.Fingerprint)
	d.Set("type", key.Type)
	d.Set("size", key.Size)
	d.Set("created_at", key.CreatedAt)
	d.SetId(key.Fingerprint)

	return diag.Diagnostics{}
}

func resourceSshKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)

	keyFingerprint := d.Id()

	key, err := c.getSshKey(ctx, keyFingerprint)
	if err != nil {
		return diag.Errorf("Unable to find SSH key %q:\n\t %q", keyFingerprint, err)
	}

	d.Set("name", key.Name)
	d.Set("data", key.Data)
	d.Set("fingerprint", key.Fingerprint)
	d.Set("type", key.Type)
	d.Set("size", key.Size)
	d.Set("created_at", key.CreatedAt)

	return diag.Diagnostics{}
}

func resourceSshKeyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)

	keyFingerprint := d.Id()
	name := d.Get("name").(string)

	key, err := c.updateSshKey(ctx, keyFingerprint, name)
	if err != nil {
		return diag.Errorf("Unable to update SSH key %q:\n\t %q", keyFingerprint, err)
	}

	d.Set("name", key.Name)

	return diag.Diagnostics{}
}

func resourceSshKeyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(HetznerRobotClient)

	keyFingerprint := d.Id()

	err := c.deleteSshKey(ctx, keyFingerprint)
	if err != nil {
		return diag.Errorf("Unable to delete SSH key %q:\n\t %q", keyFingerprint, err)
	}

	return diag.Diagnostics{}
}
