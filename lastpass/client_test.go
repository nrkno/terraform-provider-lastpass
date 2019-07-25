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
	r := Record{
		Name:     "myintegrationtest",
		URL:      "https://example.com",
		Username: "user",
		Password: "pw",
		Note:     "ABC\nDEF\nGHJ",
	}
	template := r.getTemplate()
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
	r, err = client.Create(r)
	if err != nil {
		t.Error(err)
	}
	r.Username = "user2"
	r.Password = "pw2"
	r.URL = "https://example2.com"
	r.Note = "123\n456\n789"
	r.Name = "myintegrationtest2"
	err = client.Update(r)
	if err != nil {
		t.Error(err)
	}
	r, err = client.Read(r.ID)
	if err != nil {
		t.Error(err)
	}
	if r.Name != "myintegrationtest2" && r.URL != "https://example2.com" && r.Username != "user2" && r.Password != "pw2" && r.Note != "123\n456\n789" {
		t.Error("Read() did not receive expected values.")
	}
	err = client.Delete(r.ID)
	if err != nil {
		t.Error(err)
	}
}
