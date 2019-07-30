// +build integration

// run with go test -tags=integration
package lastpass

import (
	"os"
	"os/exec"
	"testing"
)

var client Client

func init() {
	/* load test data */
	client.Username = os.Getenv("LASTPASS_USER")
	client.Password = os.Getenv("LASTPASS_PASSWORD")
}

// Run full integration test of all CRUD methods
func TestClient(t *testing.T) {
	cmd := exec.Command("lpass", "logout", "-f")
	err := cmd.Run()
	if err != nil {
		t.Error(err)
		return
	}
	err = client.login()
	if err != nil {
		t.Error(err)
	}
	s := Secret{
		Name:     "myintegrationtest",
		URL:      "https://example.com",
		Username: "user",
		Password: "pw",
		Note:     "ABC\nDEF\nGHJ",
	}
	template := s.getTemplate()
	expect := `Name: myintegrationtest
URL: https://example.com
Username: user 
Password: pw
Notes:    # Add notes below this line.
ABC
DEF
GHJ
`
	if template != expect {
		t.Error("Template not as expected")
	}
	s, err = client.Create(s)
	if err != nil {
		t.Error(err)
	}
	s.Username = "user2"
	s.Password = "pw2"
	s.URL = "https://example2.com"
	s.Note = "123\n456\n789"
	s.Name = "myintegrationtest2"
	err = client.Update(s)
	if err != nil {
		t.Error(err)
	}
	secrets, err := client.Read(s.ID)
	if err != nil {
		t.Error(err)
	}
	if secrets[0].Name != "myintegrationtest2" && secrets[0].URL != "https://example2.com" && secrets[0].Username != "user2" && secrets[0].Password != "pw2" && secrets[0].Note != "123\n456\n789" {
		t.Error("Read() did not receive expected values.")
	}
	err = client.Delete(s.ID)
	if err != nil {
		t.Error(err)
	}
}
