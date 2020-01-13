package humio

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	humio "github.com/humio/cli/api"
)

func resourceParser() *schema.Resource {
	return &schema.Resource{
		Create: resourceParserCreate,
		Read:   resourceParserRead,
		Update: resourceParserUpdate,
		Delete: resourceParserDelete,

		Schema: map[string]*schema.Schema{
			"repository": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"tag_fields": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"test_data": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"parser_script": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
		},
	}
}

func resourceParserCreate(d *schema.ResourceData, client interface{}) error {
	parser, err := parserFromResourceData(d, client)
	if err != nil {
		return fmt.Errorf("could not obtain parser from resource data: %v", err)
	}

	err = client.(*humio.Client).Parsers().Add(d.Get("repository").(string), &parser, false)
	if err != nil {
		return fmt.Errorf("could not create parser: %v", err)
	}
	d.SetId(parser.Name)

	return resourceParserRead(d, client)
}

func resourceParserRead(d *schema.ResourceData, client interface{}) error {
	parser, err := client.(*humio.Client).Parsers().Get(d.Get("repository").(string), d.Get("name").(string))
	if err != nil {
		return fmt.Errorf("could not get parser: %v", err)
	}
	resourceDataFromParser(parser, d)
	return nil
}

func resourceDataFromParser(a *humio.Parser, d *schema.ResourceData) error {
	d.Set("name", a.Name)
	d.Set("script", a.Script)
	d.Set("tag_fields", a.TagFields)
	d.Set("tests", a.Tests)
	d.SetId(d.Id())
	return nil
}

func resourceParserUpdate(d *schema.ResourceData, client interface{}) error {
	parser, err := parserFromResourceData(d, client)
	if err != nil {
		return fmt.Errorf("could not obtain parser from resource data: %v", err)
	}

	err = client.(*humio.Client).Parsers().Add(d.Get("repository").(string), &parser, true)
	if err != nil {
		return fmt.Errorf("could not create parser: %v", err)
	}
	d.SetId(parser.Name)
	return resourceParserRead(d, client)
}

func resourceParserDelete(d *schema.ResourceData, client interface{}) error {
	if err := client.(*humio.Client).Parsers().Remove(d.Get("repository").(string), d.Get("name").(string)); err != nil {
		return fmt.Errorf("could not delete parser: %v", err)
	}
	return nil
}

func parserFromResourceData(d *schema.ResourceData, client interface{}) (humio.Parser, error) {
	return humio.Parser{
		Name:      d.Get("name").(string),
		Script:    d.Get("parser_script").(string),
		TagFields: convertInterfaceListToStringSlice(d.Get("tag_fields").([]interface{})),
		Tests:     convertInterfaceListToParserTestCases(d.Get("test_data").([]interface{})),
	}, nil
}

func convertInterfaceListToParserTestCases(s []interface{}) []humio.ParserTestCase {
	var element []humio.ParserTestCase
	for _, item := range s {
		value, _ := item.(string)
		element = append(element, humio.ParserTestCase{Input: value})
	}
	return element
}
