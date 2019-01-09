package alienvault

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/form3tech-oss/alienvault"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

const testAccSensorKeyConfig = `
	resource "alienvault_sensor_key" "test-e2e-sensor-key" {}`

func TestAccResourceSensorKey(t *testing.T) {
	var key alienvault.SensorKey

	if ok, err := testAccProvider.Meta().(*alienvault.Client).HasSensorKeyAvailability(); err != nil {
		t.Fatalf("Failed to check for sensor key availability: %s", err)
	} else if !ok {
		t.Skip("Your alienvault license does not allow for the creation of any more sensors/keys")
	}

	refreshName := "alienvault_sensor_key.test-e2e-sensor-key"

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: refreshName,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSensorKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSensorKeyConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSensorKeyExists(refreshName, &key),
				),
			},
		},
	})
}

func testAccCheckSensorKeyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*alienvault.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alienvault_sensor_key" {
			continue
		}

		key, err := client.GetSensorKey(rs.Primary.ID)

		if err == nil && key != nil && !key.Consumed {
			return fmt.Errorf("key %q still exists", rs.Primary.ID)
		}

		if err != nil && !strings.Contains(err.Error(), "not found") {
			return fmt.Errorf("Unexpected error when checking for existence of key: %s", err)
		}
	}

	return nil
}

func testAccCheckSensorKeyExists(n string, res *alienvault.SensorKey) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no job ID is set")
		}

		client := testAccProvider.Meta().(*alienvault.Client)

		key, err := client.GetSensorKey(rs.Primary.ID)
		if err != nil {
			return err
		}

		*res = *key
		return nil
	}
}
