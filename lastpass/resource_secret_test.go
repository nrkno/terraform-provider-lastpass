package lastpass

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nrkno/terraform-provider-lastpass/api"
)

func TestAccResourceSecret_Basic(t *testing.T) {
	var secret api.Secret
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
						"lastpass_secret.foobar", "notes", "FOO\nBAR\n"),
				),
			},
		},
	})
}

func testAccResourceSecretDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "lastpass_secret" {
			continue
		}
		s := api.Secret{
			ID:       rs.Primary.ID,
			Name:     rs.Primary.Attributes["name"],
			Username: rs.Primary.Attributes["username"],
			Password: rs.Primary.Attributes["password"],
			URL:      rs.Primary.Attributes["url"],
			Group:    rs.Primary.Attributes["group"],
			Share:    rs.Primary.Attributes["share"],
			Notes:    rs.Primary.Attributes["notes"],
		}

		err := c.Delete(&s)
		if err != nil {
			return err
		}
		_, err = c.Read(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Secret still exists")
		}
	}
	return nil
}

func testAccResourceSecretExists(n string, secret *api.Secret) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Secret ID is set")
		}
		c := testAccProvider.Meta().(*api.Client)
		tmp, err := c.Read(rs.Primary.ID)
		if err != nil {
			return err
		}
		if tmp.ID != rs.Primary.ID {
			return fmt.Errorf("Secret not found")
		}
		*secret = *tmp
		return nil
	}
}

const testAccResourceSecretConfig_basic = `
resource "lastpass_secret" "foobar" {
    name = "terraform-provider-lastpass resource basic test"
    username = "gopher"
    password = "hunter2"
	notes = <<EOF
FOO
BAR
EOF
}`
