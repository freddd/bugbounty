package fdns

import "regexp"


type Site struct {
	Name      string
	SiteRegex *regexp.Regexp
	Verify Verify
}

type Verify struct {
	Kind      string // HTTP/NXDOMAIN
	Condition int
	Regexp    *regexp.Regexp
}

var GithubIO = Site {
	Name: "github.io",
	SiteRegex: regexp.MustCompile("github\\.io"),
	Verify: Verify {
		Kind: "HTTP",
		Condition: 404,
		Regexp: regexp.MustCompile(`(?m)There isn't a GitHub Pages site here`),
	},
}

var BitbucketIO = Site {
	Name: "bitbucket.io",
	SiteRegex: regexp.MustCompile("bitbucket\\.io"),
	Verify: Verify {
		Kind: "HTTP",
		Condition: 404,
		Regexp: regexp.MustCompile(`(?m)Repository not found`),
	},
}

var BitbucketOrg = Site {
	Name: "bitbucket.org",
	SiteRegex: regexp.MustCompile("bitbucket\\.org"),
	Verify: Verify {
		Kind: "HTTP",
		Condition: 404,
		Regexp: regexp.MustCompile(`(?m)Repository not found`),
	},
}

var CampaignMonitor = Site {
	Name: "campaignmonitor.com",
	SiteRegex: regexp.MustCompile("campaignmonitor\\.com"),
	Verify: Verify {
		Kind: "HTTP",
		Condition: 404,
		Regexp: regexp.MustCompile(`(?m)Trying to access your account`),
	},
}

var FeedPress = Site {
	Name: "feedpress.me",
	SiteRegex: regexp.MustCompile("feedpress\\.me"),
	Verify: Verify{
		Kind: "HTTP",
		Condition: 404,
		Regexp: regexp.MustCompile(`(?m)The feed has not been found."`),
	},
}

var AWSDualStack = Site {
	Name: "aws-dualstack",
	SiteRegex: regexp.MustCompile(`(?m)^[a-z0-9.\-]{0,63}\.?s3.dualstack\.(eu|ap|us|ca|sa)-\w{2,14}-\d{1,2}\.amazonaws.com$`),
	Verify: Verify{
		Kind: "HTTP",
		Condition: 404,
		Regexp: regexp.MustCompile(`(?m)The specified bucket does not exist`),
	},
}

var AWSWebsite = Site {
	Name: "aws-website",
	SiteRegex: regexp.MustCompile(`(?m)^[a-z0-9.\-]{0,63}\.?s3-website[.-](eu|ap|us|ca|sa|cn)-\w{2,14}-\d{1,2}\.amazonaws.com(\.cn)?$`),
	Verify: Verify{
		Kind: "HTTP",
		Condition: 404,
		Regexp: regexp.MustCompile(`(?m)The specified bucket does not exist`),
	},
}

var AWSGeneric = Site {
	Name: "aws-generic",
	SiteRegex: regexp.MustCompile(`(?m)^[a-z0-9.\-]{0,63}\.?s3.amazonaws\.com$`),
	Verify: Verify{
		Kind: "HTTP",
		Condition: 404,
		Regexp: regexp.MustCompile(`(?m)The specified bucket does not exist`),
	},
}

var AWSGenericWithCountry = Site {
	Name: "aws-generic-country",
	SiteRegex: regexp.MustCompile(`(?m)^[a-z0-9.\-]{0,63}\.?s3[.-](eu|ap|us|ca|sa)-\w{2,14}-\d{1,2}\.amazonaws.com$`),
	Verify: Verify{
		Kind: "HTTP",
		Condition: 404,
		Regexp: regexp.MustCompile(`(?m)The specified bucket does not exist`),
	},
}

var WordPress = Site {
	Name: "wordpress.com",
	SiteRegex: regexp.MustCompile("wordpress\\.com"),
	Verify: Verify{
		Kind: "HTTP",
		Condition: 404,
		Regexp: regexp.MustCompile(`(?m)Do you want to register`),
	},
}

var Tumblr = Site {
	Name: "tumblr.com",
	SiteRegex: regexp.MustCompile("tumblr\\.com"),
	Verify: Verify{
		Kind: "HTTP",
		Condition: 404,
		Regexp: regexp.MustCompile(`(?m)Whatever you were looking for doesn't currently exist at this address`),
	},
}

var Strikingly = Site {
	Name: "strikinglydns.com",
	SiteRegex: regexp.MustCompile("strikinglydns\\.com"),
	Verify: Verify{
		Kind: "HTTP",
		Condition: 404,
		Regexp: regexp.MustCompile(`(?m)page not found`),
	},
}

var FlyIO = Site {
	Name: "fly.io",
	SiteRegex: regexp.MustCompile("fly\\.io"),
	Verify: Verify{
		Kind: "HTTP",
		Condition: 404,
		Regexp: regexp.MustCompile(`(?m)404 Not Found`),
	},
}

var FlyDNS = Site {
	Name: "flydns.net",
	SiteRegex: regexp.MustCompile("flydns\\.net"),
	Verify: Verify{
		Kind: "HTTP",
		Condition: 404,
		Regexp: regexp.MustCompile(`.*`),
	},
}

var HelpJuice = Site {
	Name: "helpjuice.com",
	SiteRegex: regexp.MustCompile("helpjuice\\.com"),
	Verify: Verify{
		Kind: "HTTP",
		Condition: 404,
		Regexp: regexp.MustCompile(`(?m)We could not find what you're looking for.`),
	},
}

var HelpScoutDocs = Site {
	Name: "helpscoutdocs.com",
	SiteRegex: regexp.MustCompile("helpscoutdocs.com"),
	Verify: Verify{
		Kind: "HTTP",
		Condition: 404,
		Regexp: regexp.MustCompile(`(?m)No settings were found for this company`),
	},
}

var MyJetBrains = Site {
	Name: "myjetbrains.com",
	SiteRegex: regexp.MustCompile("myjetbrains\\.com"),
	Verify: Verify{
		Kind: "HTTP",
		Condition: 404,
		Regexp: regexp.MustCompile(`(?m)is not a registered InCloud YouTrack`),
	},
}

var ReadmeIO = Site {
	Name: "readme.io",
	SiteRegex: regexp.MustCompile("readme\\.io"),
	Verify: Verify{
		Kind: "HTTP",
		Condition: 404,
		Regexp: regexp.MustCompile(`(?m)"Project doesnt exist\.\.\. yet!`),
	},
}

var SurgeSH = Site {
	Name: "surge.sh",
	SiteRegex: regexp.MustCompile("surge\\.sh"),
	Verify: Verify{
		Kind: "HTTP",
		Condition: 404,
		Regexp: regexp.MustCompile(`(?m)project not found`),
	},
}

var UptimeRobot = Site {
	Name: "uptimerobot.com",
	SiteRegex: regexp.MustCompile("uptimerobot\\.com"),
	Verify: Verify{
		Kind: "HTTP",
		Condition: 404,
		Regexp: regexp.MustCompile(`(?m)page not found`),
	},
}

var CargoCollective = Site {
	Name: "cargocollective.com",
	SiteRegex: regexp.MustCompile("cargocollective\\.com"),
	Verify: Verify{
		Kind: "HTTP",
		Condition: 404,
		Regexp: regexp.MustCompile(`(?m)404 Not Found`),
	},
}

var HatenaBlog = Site {
	Name: "hatenablog.com",
	SiteRegex: regexp.MustCompile("hatenablog\\.com"),
	Verify: Verify{
		Kind: "HTTP",
		Condition: 404,
		Regexp: regexp.MustCompile(`(?m)404 Blog is not found`),
	},
}

var IntercomHelp = Site {
	Name: "intercom.help",
	SiteRegex: regexp.MustCompile("intercom\\.help"),
	Verify: Verify{
		Kind: "HTTP",
		Condition: 404,
		Regexp: regexp.MustCompile("Uh oh\\. That page doesn't exist\\."),
	},
}

var KinstaCloud = Site {
	Name: "kinsta.cloud",
	SiteRegex: regexp.MustCompile("kinsta\\.cloud"),
	Verify: Verify{
		Kind: "HTTP",
		Condition: 404,
		Regexp: regexp.MustCompile("No Site For Domain"),
	},
}

var LaunchRock = Site {
	Name: "launchrock.com",
	SiteRegex: regexp.MustCompile("launchrock\\.com"),
	Verify: Verify{
		Kind: "HTTP",
		Condition: 404,
		Regexp: regexp.MustCompile(`(?m)"It looks like you may have taken a wrong turn somewhere. Don't worry\.\.\.it happens to all of us\.`),
	},
}

var PantheonSite = Site {
	Name: "pantheonsite.io",
	SiteRegex: regexp.MustCompile("pantheonsite\\.io"),
	Verify: Verify{
		Kind: "HTTP",
		Condition: 404,
		Regexp: regexp.MustCompile(`(?m)404 error unknown site!`),
	},
}

var UserVoice = Site {
	Name: "uservoice.com",
	SiteRegex: regexp.MustCompile("uservoice\\.com"),
	Verify: Verify{
		Kind: "HTTP",
		Condition: 404,
		Regexp: regexp.MustCompile(`(?m)This UserVoice subdomain is currently available!`),
	},
}

var AzureCloudApp = Site {
	Name: "cloudapp.net",
	SiteRegex: regexp.MustCompile("cloudapp\\.net"),
	Verify: Verify{
		Kind: "HOST",
		Condition: -1,
		Regexp: regexp.MustCompile("NXDOMAIN"),
	},
}

var AzureCloudAppCom = Site {
	Name: "cloudapp.azure.com",
	SiteRegex: regexp.MustCompile("cloudapp\\.azure\\.com"),
	Verify: Verify{
		Kind: "HOST",
		Condition: -1,
		Regexp: regexp.MustCompile("NXDOMAIN"),
	},
}

var AzureWebsite = Site {
	Name: "azurewebsites.net",
	SiteRegex: regexp.MustCompile("azurewebsites\\.net"),
	Verify: Verify{
		Kind: "HOST",
		Condition: -1,
		Regexp: regexp.MustCompile("NXDOMAIN"),
	},
}

var BlobCoreWindowsNet = Site {
	Name: "blob.core.windows.net",
	SiteRegex: regexp.MustCompile("blob\\.core\\.windows\\.net"),
	Verify: Verify{
		Kind: "HOST",
		Condition: -1,
		Regexp: regexp.MustCompile("NXDOMAIN"),
	},
}

var AzureAPI = Site {
	Name: "azure-api.net",
	SiteRegex: regexp.MustCompile("azure-api\\.net"),
	Verify: Verify{
		Kind: "HOST",
		Condition: -1,
		Regexp: regexp.MustCompile("NXDOMAIN"),
	},
}

var AzureHdInsight = Site {
	Name: "azurehdinsight.net",
	SiteRegex: regexp.MustCompile("azurehdinsight\\.net"),
	Verify: Verify{
		Kind: "HOST",
		Condition: -1,
		Regexp: regexp.MustCompile("NXDOMAIN"),
	},
}

var AzureEdge = Site {
	Name: "azureedge.net",
	SiteRegex: regexp.MustCompile("azureedge\\.net"),
	Verify: Verify{
		Kind: "HOST",
		Condition: -1,
		Regexp: regexp.MustCompile("NXDOMAIN"),
	},
}

var AzureContainer = Site {
	Name: "azurecontainer.io",
	SiteRegex: regexp.MustCompile("azurecontainer\\.io"),
	Verify: Verify{
		Kind: "HOST",
		Condition: -1,
		Regexp: regexp.MustCompile("NXDOMAIN"),
	},
}

var AzureDatabaseWindows = Site {
	Name: "database.windows.net",
	SiteRegex: regexp.MustCompile("database\\.windows\\.net"),
	Verify: Verify{
		Kind: "HOST",
		Condition: -1,
		Regexp: regexp.MustCompile("NXDOMAIN"),
	},
}

var AzureDataLakeStore = Site {
	Name: "azuredatalakestore.net",
	SiteRegex: regexp.MustCompile("azuredatalakestore\\.net"),
	Verify: Verify{
		Kind: "HOST",
		Condition: -1,
		Regexp: regexp.MustCompile("NXDOMAIN"),
	},
}

var AzureSearchWindows = Site {
	Name: "search.windows.net",
	SiteRegex: regexp.MustCompile("search\\.windows\\.net"),
	Verify: Verify{
		Kind: "HOST",
		Condition: -1,
		Regexp: regexp.MustCompile("NXDOMAIN"),
	},
}

var AzureCR = Site {
	Name: "azurecr.io",
	SiteRegex: regexp.MustCompile("azurecr\\.io"),
	Verify: Verify{
		Kind: "HOST",
		Condition: -1,
		Regexp: regexp.MustCompile("NXDOMAIN"),
	},
}

var AzureRedisCache = Site {
	Name: "redis.cache.windows.net",
	SiteRegex: regexp.MustCompile("redis\\.cache\\.windows\\.net"),
	Verify: Verify{
		Kind: "HOST",
		Condition: -1,
		Regexp: regexp.MustCompile("NXDOMAIN"),
	},
}

var AzureServicebusWindows = Site {
	Name: "servicebus.windows.net",
	SiteRegex: regexp.MustCompile("servicebus\\.windows\\.net"),
	Verify: Verify{
		Kind: "HOST",
		Condition: -1,
		Regexp: regexp.MustCompile("NXDOMAIN"),
	},
}

var VisualStudio = Site {
	Name: "visualstudio.com",
	SiteRegex: regexp.MustCompile("visualstudio\\.com"),
	Verify: Verify{
		Kind: "HOST",
		Condition: -1,
		Regexp: regexp.MustCompile("NXDOMAIN"),
	},
}

var TrafficManager = Site {
	Name: "trafficmanager.net",
	SiteRegex: regexp.MustCompile("trafficmanager\\.net"),
	Verify: Verify{
		Kind: "HOST",
		Condition: -1,
		Regexp: regexp.MustCompile("NXDOMAIN"),
	},
}

var Heroku = Site {
	Name: "herokudns.com",
	SiteRegex: regexp.MustCompile("herokudns\\.com"),
	Verify: Verify{
		Kind: "HTTP",
		Condition: 404,
		Regexp: regexp.MustCompile(`(?m)No such app`),
	},
}

var Fastly = Site {
	Name: ".fastly.net",
	SiteRegex: regexp.MustCompile("\\.fastly\\.net"),
	Verify: Verify{
		Kind: "HTTP",
		Condition: 500,
		Regexp: regexp.MustCompile(`(?m)unknown domain`),
	},
}

var TeamTailor = Site {
	Name: "teamtailor.com",
	SiteRegex: regexp.MustCompile("ext\\.teamtailor\\.com"),
	Verify: Verify{
		Kind: "HTTP",
		Condition: 409,
		Regexp: regexp.MustCompile(`.*`),
	},
}


var VulnerableSites = []Site {
	GithubIO,
	FeedPress,
	AWSDualStack,
	AWSGeneric,
	AWSGenericWithCountry,
	AWSWebsite, WordPress,
	Tumblr,
	Strikingly,
	FlyIO,
	FlyDNS,
	HelpJuice,
	HelpScoutDocs,
	MyJetBrains,
	ReadmeIO,
	SurgeSH,
	UptimeRobot,
	CargoCollective,
	HatenaBlog,
	IntercomHelp,
	KinstaCloud,
	PantheonSite,
	UserVoice,
	AzureCloudApp,
	AzureCloudAppCom,
	AzureWebsite,
	BlobCoreWindowsNet,
	AzureAPI,
	AzureHdInsight,
	AzureEdge,
	AzureContainer,
	AzureDatabaseWindows,
	AzureDataLakeStore,
	AzureSearchWindows,
	AzureCR,
	AzureRedisCache,
	AzureServicebusWindows,
	VisualStudio,
	TrafficManager,
	Heroku,
	CampaignMonitor,
	BitbucketIO,
	BitbucketOrg,
	//Fastly,
	//TeamTailor,
}