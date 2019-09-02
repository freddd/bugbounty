package third_party

import (
	"crypto/tls"
	"fdns/logger"
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/parnurzeal/gorequest"
	"os"
	"regexp"
	"strings"
	"time"
)

type Target struct {
	Name            string
	Vulnerabilities []Vulnerability
}

type Vulnerability struct {
	Condition   int
	Description string
	Check       func(domain string)
}

var Nexus = Target{
	Name: "Nexus Repository",
	Vulnerabilities: []Vulnerability{
		{
			Condition: 204,
			Description: "Default credentials usernames/password",
			Check: func(domain string) {
				//logger.DefaultLogger.Debug("%s: Credentials (nexus) - starting", domain)

				anonymous := checkDefaultUsernameAndPasswordInNexus(domain, "YW5vbnltb3Vz", "YW5vbnltb3Vz")
				if anonymous {
					logger.DefaultLogger.Info("Domain: %s is vulnerable login with anonymous credentials", domain)
				}

				admin := checkDefaultUsernameAndPasswordInNexus(domain, "YWRtaW4=", "YWRtaW4xMjM=")
				if admin {
					logger.DefaultLogger.Info("Domain: %s is vulnerable login with admin credentials", domain)
				}
			},
		},
		{
			Condition:   200,
			Description: "CVE-2019-7238, RCE via API, version is lower than 3.15.0",
			Check: func(domain string) {
				//logger.DefaultLogger.Debug("%s: CVE-2019-7238, RCE via API - starting", domain)

				vulnerableVersion, _ := semver.NewVersion("3.15.0")
				res, body, httpErrors := execute(domain)

				if httpErrors != nil {
					return
				}

				if res.StatusCode != 200 {
					return
				}

				re := regexp.MustCompile(`_v=(\d.*)"`)
				unsanitized := re.FindString(body)
				if unsanitized == "" {
					return
				}

				sanitized := unsanitized[strings.LastIndex(unsanitized, "_v=") + 3:strings.Index(unsanitized, "-")]


				v, err := semver.NewVersion(sanitized)
				if err != nil {
					logger.DefaultLogger.Error("domain: %s, %+v (%s)", domain, err, unsanitized)
				}


				if v.LessThan(vulnerableVersion) {
					logger.DefaultLogger.Info("Domain: %s is vulnerable potentially vulnerable to RCE (CVE-2019-7238)", domain)
				}
			},
		},
	},
}

var Artifactory = Target{
	Name: "JFrog Artifactory",
	Vulnerabilities: []Vulnerability{
		{
			Condition: 200,
			Description: "Default admin credentials usernames/password (CVE-2019-9733)",
			Check: func(domain string) {
				logger.DefaultLogger.Debug("%s: CVE-2019-9733 (jfrog) - starting", domain)

				accessAdmin := checkDefaultUsernameAndPasswordInArtifactory(domain, "access-admin", "password")
				if accessAdmin {
					logger.DefaultLogger.Info("Domain: %s is vulnerable login with access-admin credentials", domain)
				}

				admin := checkDefaultUsernameAndPasswordInArtifactory(domain, "admin", "password")
				if admin {
					logger.DefaultLogger.Info("Domain: %s is vulnerable login with admin credentials", domain)
				}
			},
		},
	},
}
/*
var Jira = Target{
	Name: "Atlassian Jira",
	Vulnerabilities: Vulnerabilities{
		Condition: 200,
		Check: func(domain string) bool {
			return false
		},
	},
}
*/
var Jenkins = Target{
	Name: "Jenkins CI",
	Vulnerabilities: []Vulnerability{
		{
			Condition: 200,
			Description: "CVE-2018-1000861 (Jenkins version < 2.138)",
			Check: func(domain string) {
				vulnerableVersion, err := semver.NewVersion("2.138")
				res, _, errs := execute(domain)

				if errs != nil {
					return
				}

				if res.StatusCode != 200 {
					return
				}

				logger.DefaultLogger.Info("%s: returned: %d (%s)", domain, res.StatusCode, res.Header.Get("X-Jenkins"))

				v, err := semver.NewVersion(res.Header.Get("X-Jenkins"))
				if err != nil {
					return
				}

				if v.LessThan(vulnerableVersion) {
					logger.DefaultLogger.Info("Domain: %s is potentially vulnerable to CVE-2018-1000861", domain)
				}
			},
		},
	},
}

func execute(domain string) (gorequest.Response, string, []error) {
	resp, body, httpErrors := gorequest.
		New().
		TLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		Get(fmt.Sprintf("https://%s", domain)).
		Timeout(3*time.Second).
		End()

	for _, err := range httpErrors {
		if os.IsTimeout(err) {
			resp, body, httpErrors = gorequest.
				New().
				TLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
				Get(fmt.Sprintf("http://%s", domain)).
				Timeout(3*time.Second).
				End()
		}
	}

	return resp, body, httpErrors
}

func nexusLoginRequest(protocol string, domain string, username string, password string) *gorequest.SuperAgent {
	return gorequest.
		New().
		Post(fmt.Sprintf("%s://%s/service/rapture/session", protocol, domain)).
		TLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		Set("sec-fetch-mode", "cors").
		Set("x-requested-with", "XMLHttpRequest").
		Set("authority", domain).
		Set("sec-fetch-site", "same-origin").
		Set("origin", fmt.Sprintf("%s://%s", protocol, domain)).
		Set("x-nexus-ui", "true").
		Set("accept", "*/*").
		Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36").
		Set("referer", fmt.Sprintf("%s://%s/", protocol, domain)).
		Type("form").
		Timeout(3*time.Second).
		SendMap(map[string]string{"username": username, "password": password})
}

func artifactoryLoginRequest(protocol string, domain string, username string, password string) *gorequest.SuperAgent {
	return gorequest.
		New().
		Post(fmt.Sprintf("%s://%s/artifactory/ui/auth/login?_spring_security_remember_me=false", protocol, domain)).
		TLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		Type("json").
		Set("X-Forwarded-For", "127.0.0.1").
		Set("Sec-Fetch-Mode", "cors").
		Set("Sec-Fetch-Site","same-origin").
		Set("X-Requested-With", "artUI").
		Set("Request-Agent", "artifactoryUI").
		Set("referer", fmt.Sprintf("%s://%s/artifactory/webapp/", protocol, domain)).
		Set("origin", fmt.Sprintf("%s://%s", protocol, domain)).
		Set("accept", "*/*").
		Set("serial",  "61").
		Set("authority", domain).
		Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36").
		Timeout(3*time.Second).
		Send(fmt.Sprintf(`{"user": "%s", "password": "%s", "type":"login"}`, username, password))
}

func checkDefaultUsernameAndPasswordInNexus(domain string, username string, password string) bool {
	resp, _, httpErrors := nexusLoginRequest("https", domain, username, password).End()

	for _, err := range httpErrors {
		if os.IsTimeout(err) {
			resp, _, httpErrors = nexusLoginRequest("http", domain, username, password).End()
		}
	}

	if len(httpErrors) > 0 {
		return false
	}

	if resp.StatusCode == 204 {
		return true
	}

	return false
}

func checkDefaultUsernameAndPasswordInArtifactory(domain string, username string, password string) bool {
	resp, _, httpErrors := artifactoryLoginRequest("https", domain, username, password).End()

	for _, err := range httpErrors {
		if os.IsTimeout(err) {
			resp, _, httpErrors = artifactoryLoginRequest("http", domain, username, password).End()
		}
	}

	if len(httpErrors) > 0 {
		return false
	}

	if resp.StatusCode == 200 {
		return true
	}

	return false
}

var Targets = map[string]Target{
	"Nexus": Nexus,
	"Artifactory": Artifactory,
	"Jenkins": Jenkins,
}