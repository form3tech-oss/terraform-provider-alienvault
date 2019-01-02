package alienvault

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProviderConfigure(t *testing.T) {

	authCalled := false

	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authCalled = true
	}))
	defer ts.Close()

	resourceSchema := map[string]*schema.Schema{
		"fqdn": &schema.Schema{
			Type: schema.TypeString,
		},
		"username": &schema.Schema{
			Type: schema.TypeString,
		},
		"password": &schema.Schema{
			Type: schema.TypeString,
		},
	}
	resourceDataMap := map[string]interface{}{
		"fqdn":     strings.Replace(ts.URL, "https://", "", -1),
		"username": "something",
		"password": "something",
	}
	resourceLocalData := schema.TestResourceDataRaw(t, resourceSchema, resourceDataMap)

	provider, err := providerConfigure(resourceLocalData)
	require.Nil(t, err)

	_, ok := provider.(*Client)
	require.True(t, ok)

	assert.True(t, authCalled)

}
