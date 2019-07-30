package main

import (
	"fmt"
	"testing"

	"github.com/nrkno/terraform-provider-lastpass/lastpass"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccResourceSecret_Basic(t *testing.T) {
	var secret lastpass.Secret
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccResourceSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSecretConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccResourceSecretExists("lastpass_secret.foobar", &secret),
					resource.TestCheckResourceAttr(
						"lastpass_secret.foobar", "name", "terraform-provider-lastpass resource basic test"),
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

func TestAccResourceSecret_Generated(t *testing.T) {
	var secret lastpass.Secret
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccResourceSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSecretConfig_generated,
				Check: resource.ComposeTestCheckFunc(
					testAccResourceSecretExists("lastpass_secret.foobar", &secret),
					testAccResourceSecretGeneratedLength("lastpass_secret.foobar", 24),
					resource.TestCheckResourceAttr(
						"lastpass_secret.foobar", "name", "terraform-provider-lastpass resource generated test"),
				),
			},
		},
	})
}

func testAccResourceSecretDestroy(s *terraform.State) error {
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

func testAccResourceSecretExists(n string, secret *lastpass.Secret) resource.TestCheckFunc {
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

func testAccResourceSecretGeneratedLength(n string, length int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		provider := testAccProvider.Meta().(*lastpass.Client)
		secrets, err := provider.Read(rs.Primary.ID)
		if err != nil {
			return err
		}
		if len(secrets) != 1 && len(secrets[0].Password) != length {
			return fmt.Errorf("password has wrong length %d", len(secrets[0].Password))
		}
		return nil
	}
}

const testAccResourceSecretConfig_basic = `
resource "lastpass_secret" "foobar" {
    name = "terraform-provider-lastpass resource basic test"
    username = "gopher"
    password = "hunter2"
    note = "secret note"
}`

const testAccResourceSecretConfig_generated = `
resource "lastpass_secret" "foobar" {
	name = "terraform-provider-lastpass resource generated test"
	generate {
        length = 24
        use_symbols = false
    }
}`
