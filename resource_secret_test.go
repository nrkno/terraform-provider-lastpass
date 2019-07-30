package main

import (
	"fmt"
	"testing"

	"github.com/nrkno/terraform-provider-lastpass/lastpass"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccLastpassSecret_Basic(t *testing.T) {
	var secret lastpass.Secret

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLastpassSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLastpassSecretConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLastpassSecretExists("lastpass_secret.foobar", &secret),
					testAccCheckLastpassSecretAttributes(&secret),
					resource.TestCheckResourceAttr(
						"lastpass_secret.foobar", "name", "terraform-provider-lastpass basic test"),
					resource.TestCheckResourceAttr(
						"lastpass_secret.foobar", "username", "gopher"),
					resource.TestCheckResourceAttr(
						"lastpass_secret.foobar", "password", "hunter2"),
					resource.TestCheckResourceAttr(
						"lastpass_secret.foobar", "note", "secret note"),
				),
			},
		},
	})
}

func testAccCheckLastpassSecretDestroy(s *terraform.State) error {
	provider := testAccProvider.Meta().(*lastpass.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "lastpass_secret" {
			continue
		}

		_, err := provider.Read(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Secret still exists")
		}
	}
	return nil
}

func testAccCheckLastpassSecretExists(n string, secret *lastpass.Secret) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Secret ID is set")
		}
		provider := testAccProvider.Meta().(*lastpass.Client)
		secrets, err := provider.Read(rs.Primary.ID)
		if err != nil {
			return err
		}
		if len(secrets) != 1 && secrets[0].ID != rs.Primary.ID {
			return fmt.Errorf("Secret not found")
		}
		*secret = secrets[0]
		return nil
	}
}

func testAccCheckLastpassSecretAttributes(secret *lastpass.Secret) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if secret.Name != "terraform-provider-lastpass basic test" {
			return fmt.Errorf("Bad content: %s", secret.Name)
		}
		return nil
	}
}

const testAccCheckLastpassSecretConfig_basic = `
resource "lastpass_secret" "foobar" {
    name = "terraform-provider-lastpass basic test"
    username = "gopher"
    password = "hunter2"
    note = "secret note"
}`
