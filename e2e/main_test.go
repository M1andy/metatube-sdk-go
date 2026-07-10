package e2e

import (
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/metatube-community/metatube-sdk-go/provider/10musume"
	"github.com/metatube-community/metatube-sdk-go/provider/1pondo"
	"github.com/metatube-community/metatube-sdk-go/provider/avbase"
	"github.com/metatube-community/metatube-sdk-go/provider/caribbeancom"
	"github.com/metatube-community/metatube-sdk-go/provider/dahlia"
	"github.com/metatube-community/metatube-sdk-go/provider/fanza"
	"github.com/metatube-community/metatube-sdk-go/provider/fc2"
	"github.com/metatube-community/metatube-sdk-go/provider/heyzo"
	"github.com/metatube-community/metatube-sdk-go/provider/javbus"
	"github.com/metatube-community/metatube-sdk-go/provider/mgstage"
	"github.com/metatube-community/metatube-sdk-go/provider/tokyo-hot"
	mt "github.com/metatube-community/metatube-sdk-go/provider"
)

func proxyCheck(t *testing.T) {
	t.Helper()
	if os.Getenv("E2E_SKIP_IF_NO_PROXY") == "" {
		return
	}
	proxy := os.Getenv("HTTP_PROXY")
	if proxy == "" {
		proxy = os.Getenv("HTTPS_PROXY")
	}
	if proxy == "" {
		t.Skip("no proxy configured")
	}
	host, port, err := net.SplitHostPort(proxy)
	if err != nil {
		host = proxy
		port = "80"
	}
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), 5*time.Second)
	if err != nil {
		t.Skipf("proxy %s not reachable: %v", proxy, err)
	}
	conn.Close()
	t.Logf("proxy %s reachable", proxy)
}

// === Layer 1 (P0): AVBase — must pass ===

func TestE2E_AVBase_GetBuildID(t *testing.T) {
	proxyCheck(t)
	ab := avbase.New()
	buildID, err := ab.GetBuildID()
	require.NoError(t, err, "AVBase getBuildID failed — uTLS fingerprint issue?")
	require.NotEmpty(t, buildID)
	t.Logf("buildID: %s", buildID)
}

func TestE2E_AVBase_GetMovieInfoByID(t *testing.T) {
	proxyCheck(t)
	ab := avbase.New()
	for _, id := range []string{"prestige:ABP-588", "tameike:MEYD-856", "SSIS-354", "glory:GVH-186", "idp:IPX-534"} {
		t.Run(id, func(t *testing.T) {
			info, err := ab.GetMovieInfoByID(id)
			require.NoError(t, err)
			require.NotNil(t, info)
			require.True(t, info.IsValid())
			require.NotEmpty(t, info.Title)
			t.Logf("title=%q provider=%s", info.Title, info.Provider)
		})
		time.Sleep(500 * time.Millisecond)
	}
}

func TestE2E_AVBase_SearchMovie(t *testing.T) {
	proxyCheck(t)
	ab := avbase.New()
	for _, kw := range []string{"ABP-588", "MEYD-856", "SSIS-354"} {
		t.Run(kw, func(t *testing.T) {
			results, err := ab.SearchMovie(kw)
			require.NoError(t, err)
			require.NotEmpty(t, results)
			t.Logf("%d results", len(results))
		})
		time.Sleep(500 * time.Millisecond)
	}
}

// === Layer 2 (P1): Main providers ===

func TestE2E_JavBus_GetMovieInfoByID(t *testing.T) {
	proxyCheck(t)
	testGetMovieInfoByID(t, javbus.New(), []string{"SMBD-77", "SSNI-776", "ABP-331"})
}

func TestE2E_FANZA_GetMovieInfoByID(t *testing.T) {
	proxyCheck(t)
	testGetMovieInfoByID(t, fanza.New(), []string{"prst00022", "196glod00359", "midv00047"})
}

func TestE2E_Caribbeancom_GetMovieInfoByID(t *testing.T) {
	proxyCheck(t)
	testGetMovieInfoByID(t, caribbeancom.New(), []string{"050422-001", "031222-001"})
}

func TestE2E_HEYZO_GetMovieInfoByID(t *testing.T) {
	proxyCheck(t)
	testGetMovieInfoByID(t, heyzo.New(), []string{"0841", "0805"})
}

func TestE2E_1Pondo_GetMovieInfoByID(t *testing.T) {
	proxyCheck(t)
	testGetMovieInfoByID(t, onepondo.New(), []string{"071319_870", "042922_001"})
}

func TestE2E_MGS_GetMovieInfoByID(t *testing.T) {
	proxyCheck(t)
	testGetMovieInfoByID(t, mgstage.New(), []string{"300MAAN-778"})
}

// === Layer 3 (P2): Additional providers ===

func TestE2E_DAHLIA_GetMovieInfoByID(t *testing.T) {
	proxyCheck(t)
	testGetMovieInfoByID(t, dahlia.New(), []string{"dldss339"})
}

func TestE2E_TOKYOHOT_GetMovieInfoByID(t *testing.T) {
	proxyCheck(t)
	testGetMovieInfoByID(t, tokyohot.New(), []string{"n1633"})
}

func TestE2E_FC2_GetMovieInfoByID(t *testing.T) {
	proxyCheck(t)
	testGetMovieInfoByID(t, fc2.New(), []string{"2812904"})
}

func TestE2Etenmusume_GetMovieInfoByID(t *testing.T) {
	proxyCheck(t)
	testGetMovieInfoByID(t, tenmusume.New(), []string{"042922_01"})
}

func testGetMovieInfoByID(t *testing.T, p mt.MovieProvider, ids []string) {
	t.Helper()
	for i, id := range ids {
		t.Run(id, func(t *testing.T) {
			info, err := p.GetMovieInfoByID(id)
			require.NoError(t, err, "provider=%s", p.Name())
			require.NotNil(t, info)
			require.True(t, info.IsValid())
			require.NotEmpty(t, info.Title)
			t.Logf("%s: title=%q", p.Name(), info.Title)
		})
		if i < len(ids)-1 {
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func TestMain(m *testing.M) {
	fmt.Println("=== metatube-sdk-go E2E Tests ===")
	fmt.Printf("Fingerprint mode: %s\n", os.Getenv("MT_FINGERPRINT_MODE"))
	fmt.Printf("HTTP_PROXY: %s\n", os.Getenv("HTTP_PROXY"))
	os.Exit(m.Run())
}
