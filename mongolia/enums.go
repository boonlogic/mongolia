package mongolia

import (
	"encoding/base64"
	"math/rand"
	"strings"
	"time"
)

func dropped() string {
	if tmp, err := base64.StdEncoding.DecodeString("YW5uaWhpbGF0ZWQuCmJsb3R0ZWQgb3V0LgpkZXN0cm95ZWQuCmRlbW9saXNoZWQuCmVsaW1pbmF0ZWQuCmV4cHVuZ2VkLgpleHRlcm1pbmF0ZWQuCmV4dGlycGF0ZWQuCmxpcXVpZGF0ZWQuCm9ibGl0ZXJhdGVkLgpqdXN0IGdvdCBjYW5jZWxsZWQuCndhcyBzdW1tYXJpbHkgZXhlY3V0ZWQuCmhhcyBiZWVuIGNvbnNpZ25lZCB0byBvYmxpdmlvbi4K"); err == nil {
		syns := strings.Split(string(tmp), "\n")
		r := rand.New(rand.NewSource(time.Now().Unix()))
		if i := r.Intn(20); i == 0 {
			return syns[r.Intn(len(syns))]
		}
	}
	return "dropped"
}
