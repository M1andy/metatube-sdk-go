package avbase

import (
	"testing"

	"github.com/metatube-community/metatube-sdk-go/provider/internal/testkit"
)

func TestAVBase_GetMovieInfoByID(t *testing.T) {
	testkit.Test(t, New, []string{
		"prestige:ABP-588",
		"tameike:MEYD-856",
		"SSIS-354",
		"glory:GVH-186",
		"idp:IPX-534",
		"idp:IPX-867",
		"SSNI-369",
		"SSNI-392",
		"SSNI-497",
		"SSNI-520",
		"SSNI-436",
		"PPPE-052",
	})
}

func TestAVBase_SearchMovie(t *testing.T) {
	testkit.Test(t, New, []string{
		"ABP-588",
		"MEYD-856",
		"SSIS-354",
		"HMN",
	})
}
