package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
)

type record []struct {
	Fullname        string `json:"fullname"`
	Group           string `json:"group"`
	ID              string `json:"id"`
	LastModifiedGmt string `json:"last_modified_gmt"`
	LastTouch       string `json:"last_touch"`
	Name            string `json:"name"`
	Note            string `json:"note"`
	Password        string `json:"password"`
	Share           string `json:"share"`
	URL             string `json:"url"`
	Username        string `json:"username"`
}

// ResourceRecord describes our lastpass record resource
func ResourceRecord() *schema.Resource {
	return &schema.Resource{
		Create: ResourceRecordCreate,
		Read:   ResourceRecordRead,
		Update: ResourceRecordUpdate,
		Delete: ResourceRecordDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"note": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func getTemplate(d *schema.ResourceData) string {
	name := d.Get("name").(string)
	url := d.Get("url").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	note := d.Get("note").(string)
	template := fmt.Sprintf(`
	Name: %s
	URL: %s
	Username: %s 
	Password: %s
	Notes:    # Add notes below this line.
	%s
	`, name, url, username, password, note)
	return template
}

// ResourceRecordCreate is used to create a new resource and generate ID.
func ResourceRecordCreate(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	template := getTemplate(d)
	cmd := exec.Command("lpass", "add", name, "--non-interactive", "--sync=now")
	var inbuf, errbuf bytes.Buffer
	inbuf.Write([]byte(template))
	cmd.Stdin = &inbuf
	cmd.Stderr = &errbuf
	err := cmd.Run()
	if err != nil {
		var err = errors.New(errbuf.String())
		return err
	}
	time.Sleep(time.Second * 5) // Need to finish sync with upstream/lastpass before we get actual ID.
	cmd = exec.Command("lpass", "show", name, "--json", "-x")
	var outbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	err = cmd.Run()
	if err != nil {
		var err = errors.New(errbuf.String())
		return err
	}
	var records record
	err = json.Unmarshal(outbuf.Bytes(), &records)
	if err != nil {
		return err
	}
	if len(records) > 1 {
		err := errors.New("more than one record with same name")
		return err
	}
	if records[0].ID == "0" {
		err := errors.New("got invalid ID 0, possible problem with lpass sync")
		return err
	}
	d.SetId(records[0].ID)
	return ResourceRecordRead(d, m)
}

// ResourceRecordRead is used to sync the local state with the actual state (upstream/lastpass)
func ResourceRecordRead(d *schema.ResourceData, m interface{}) error {
	cmd := exec.Command("lpass", "show", d.Id(), "--json", "-x")
	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	err := cmd.Run()
	if err != nil {
		// Make sure the record is not removed manually.
		if strings.Contains(errbuf.String(), "Could not find specified account") {
			d.SetId("")
			return nil
		}
		var err = errors.New(errbuf.String())
		return err
	}
	var records record
	err = json.Unmarshal(outbuf.Bytes(), &records)
	if err != nil {
		return err
	}
	d.Set("name", records[0].Name)
	d.Set("url", records[0].URL)
	d.Set("username", records[0].Username)
	d.Set("password", records[0].Password)
	d.Set("note", records[0].Note)
	return nil
}

// ResourceRecordUpdate is used to update our existing resource
func ResourceRecordUpdate(d *schema.ResourceData, m interface{}) error {
	template := getTemplate(d)
	cmd := exec.Command("lpass", "edit", d.Id(), "--non-interactive", "--sync=now")
	var inbuf, errbuf bytes.Buffer
	inbuf.Write([]byte(template))
	cmd.Stdin = &inbuf
	cmd.Stderr = &errbuf
	err := cmd.Run()
	if err != nil {
		var err = errors.New(errbuf.String())
		return err
	}
	return ResourceRecordRead(d, m)
}

// ResourceRecordDelete is exactly what it sounds like - it is called to destroy the resource.
func ResourceRecordDelete(d *schema.ResourceData, m interface{}) error {
	var errbuf bytes.Buffer
	cmd := exec.Command("lpass", "rm", d.Id(), "--sync=now")
	cmd.Stderr = &errbuf
	err := cmd.Run()
	if err != nil {
		// Make sure the record is not removed manually.
		if strings.Contains(errbuf.String(), "Could not find specified account") {
			return nil
		}
		var err = errors.New(errbuf.String())
		return err
	}
	return nil
}
