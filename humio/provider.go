package humio

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	humio "github.com/humio/cli/api"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ConfigureFunc: func(r *schema.ResourceData) (interface{}, error) {
			client, err := humio.NewClient(humio.Config{
				Address: r.Get("addr").(string),
				Token:   r.Get("api_token").(string),
			})
			if err != nil {
				panic(fmt.Sprintf("could not create humio client: %v", err))
			}
			return client, nil
		},
		ResourcesMap: map[string]*schema.Resource{
			"humio_alert":        resourceAlert(),
			"humio_ingest_token": resourceIngestToken(),
			"humio_parser":       resourceParser(),
			"humio_notifier":     resourceNotifier(),
		},
		Schema: map[string]*schema.Schema{
			"addr": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HUMIO_ADDR", "https://cloud.humio.com/"),
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if v[len(v)-1] != '/' {
						// TODO: determine if we really want to enforce this.
						return warns, append(errs, fmt.Errorf("error: address '%q' must contain a trailing '/', got: %s", key, v))
					}
					return validateURL(val, key)
				},
			},
			"api_token": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("HUMIO_API_TOKEN", nil),
			},
		},
	}
}

func validateURL(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)
	u, err := url.Parse(v)
	if err != nil {
		errs = append(errs, fmt.Errorf("error: %s is not a valid URL", v))
	} else if u.Scheme == "" || u.Host == "" {
		errs = append(errs, fmt.Errorf("error: %s must be an absolute URL", v))
	} else if u.Scheme != "http" && u.Scheme != "https" {
		errs = append(errs, fmt.Errorf("error: %s must begin with http or https", v))
	}
	return warns, errs
}
