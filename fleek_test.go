package fleek

import (
	"fmt"
	"os"
	"testing"
)

func TestGetSitesByTeamId(t *testing.T) {
	f, err := New(os.Getenv("FLEEK_API_TOKEN"))
	if err != nil {
		t.Fatalf("%s", err)
	}

	sites, err := f.GetSitesByTeamId("mrusme-team")
	if err != nil {
		t.Fatalf("%s", err)
	}

	for _, site := range sites {
		fmt.Printf(
			"Site ID: %v\nName: %s\nPlatform: %s\nUpdated at: %s\n\n",
			site.Id,
			site.Name,
			site.Platform,
			site.UpdatedAt,
		)
	}

}

func TestGetSiteBySlug(t *testing.T) {
	f, err := New(os.Getenv("FLEEK_API_TOKEN"))
	if err != nil {
		t.Fatalf("%s", err)
	}

	site, err := f.GetSiteBySlug("xn-gckvb8fzb")
	if err != nil {
		t.Fatalf("%s", err)
	}

	fmt.Printf(
		"Site ID: %v\nName: %s\nPlatform: %s\nTeam ID: %v\n\n",
		site.Id,
		site.Name,
		site.Platform,
		site.Team.Id,
	)

}
