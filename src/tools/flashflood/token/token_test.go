package token_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"tools/flashflood/token"
)

var _ = Describe("GetTokenFromUAA", func() {
	It("makes correct request to UAA", func() {
		server := ghttp.NewServer()

		server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST",
					"/oauth/token",
				),
				ghttp.VerifyContentType("application/x-www-form-urlencoded"),
				ghttp.VerifyBasicAuth("cf", ""),
				func(w http.ResponseWriter, req* http.Request) { // Verify data in body
					req.ParseForm()
					Expect(req.Form).To(HaveKeyWithValue("grant_type", []string{"password"}))
					Expect(req.Form).To(HaveKeyWithValue("username", []string{"u"}))
					Expect(req.Form).To(HaveKeyWithValue("password", []string{"p"}))
					Expect(req.Form).To(HaveKeyWithValue("client_id", []string{"cf"}))
				},
				ghttp.RespondWithJSONEncoded(http.StatusOK, map[string]string{
					"token_type":   "bearer",
					"access_token": "token",
				}),
			))

		token, err := token.GetTokenFromUAA("u", "p", server.URL(), false)

		Expect(server.ReceivedRequests()).To(HaveLen(1))

		Expect(err).NotTo(HaveOccurred())
		Expect(token).To(Equal("bearer token"))
	})
})
