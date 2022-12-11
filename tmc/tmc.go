package tmc

import (
	"io"
	"net/http"
	"sync"
	"time"
)

type Func func(key string) (interface{}, error)

type TMCache struct {
	done  chan struct{}
	mu    sync.Mutex
	items map[string]*item
	f     Func
}

type item struct {
	done     chan struct{}
	deadline int64
	res      result
}

type result struct {
	value interface{}
	err   error
}

func NewTMCache(fun Func, cleanupTimeout time.Duration) *TMCache {
	tmc := &TMCache{
		done:  make(chan struct{}),
		items: make(map[string]*item),
		f:     fun,
	}

	go cleanup(tmc, cleanupTimeout)

	return tmc
}

func (tmc *TMCache) routineCleanup() {
	now := time.Now().UnixNano()

	tmc.mu.Lock()

	for k, i := range tmc.items {
		if i.deadline <= now {
			delete(tmc.items, k)
		}
	}

	tmc.mu.Unlock()
}

func cleanup(tmc *TMCache, cleanupTimeout time.Duration) {
	ticker := time.NewTicker(cleanupTimeout)

	for {
		select {
		case <-ticker.C:
			tmc.routineCleanup()
		case <-tmc.done:
			ticker.Stop()
			return
		}
	}
}

func (tmc *TMCache) Get(key string, ttl time.Duration) (interface{}, bool, error) {
	chit := false
	tmc.mu.Lock()
	i := tmc.items[key]
	if i == nil {
		i = &item{
			deadline: time.Now().UnixNano() + int64(ttl),
			done:     make(chan struct{}),
		}
		tmc.items[key] = i
		tmc.mu.Unlock()

		i.res.value, i.res.err = tmc.f(key)

		close(i.done)
	} else {
		chit = true
		tmc.mu.Unlock()
		<-i.done
	}
	return i.res.value, chit, i.res.err
}

func (tmc *TMCache) Del(key string) {
	tmc.mu.Lock()
	delete(tmc.items, key)
	tmc.mu.Unlock()
}

func (tmc *TMCache) EraseAll() {
	for k := range tmc.items {
		delete(tmc.items, k)
	}
}

func (tmc *TMCache) Close() {
	tmc.mu.Lock()

	close(tmc.done)
	tmc.EraseAll()
	tmc.items = nil

	tmc.mu.Unlock()
}

func HttpGetBody(url string) (interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

var TestUrls []string = []string{"www.fonts.googleapis.com", "www.facebook.com", "www.twitter.com", "www.google.com", "www.youtube.com", "www.s.w.org", "www.instagram.com", "www.googletagmanager.com", "www.linkedin.com", "www.ajax.googleapis.com", "www.plus.google.com", "www.gmpg.org", "www.pinterest.com", "www.fonts.gstatic.com", "www.wordpress.org", "www.en.wikipedia.org", "www.youtu.be", "www.maps.google.com", "www.itunes.apple.com", "www.github.com", "www.bit.ly", "www.play.google.com", "www.goo.gl", "www.docs.google.com", "www.cdnjs.cloudflare.com", "www.vimeo.com", "www.support.google.com", "www.google-analytics.com", "www.maps.googleapis.com", "www.flickr.com", "www.vk.com", "www.t.co", "www.reddit.com", "www.amazon.com", "www.medium.com", "www.sites.google.com", "www.drive.google.com", "www.creativecommons.org", "www.microsoft.com", "www.developers.google.com", "www.adobe.com", "www.soundcloud.com", "www.theguardian.com", "www.apis.google.com", "www.ec.europa.eu", "www.lh3.googleusercontent.com", "www.chrome.google.com", "www.cloudflare.com", "www.nytimes.com", "www.maxcdn.bootstrapcdn.com", "www.support.microsoft.com", "www.blogger.com", "www.forbes.com", "www.s3.amazonaws.com", "www.code.jquery.com", "www.dropbox.com", "www.translate.google.com", "www.paypal.com", "www.apps.apple.com", "www.tinyurl.com", "www.etsy.com", "www.theatlantic.com", "www.m.facebook.com", "www.archive.org", "www.amzn.to", "www.cnn.com", "www.policies.google.com", "www.commons.wikimedia.org", "www.issuu.com", "www.i.imgur.com", "www.wordpress.com", "www.wp.me", "www.businessinsider.com", "www.yelp.com", "www.mail.google.com", "www.support.apple.com", "www.t.me", "www.apple.com", "www.washingtonpost.com", "www.bbc.com", "www.gstatic.com", "www.imgur.com", "www.amazon.de", "www.bbc.co.uk", "www.googleads.g.doubleclick.net", "www.mozilla.org", "www.eventbrite.com", "www.slideshare.net", "www.w3.org", "www.forms.gle", "www.platform.twitter.com", "www.accounts.google.com", "www.telegraph.co.uk", "www.messenger.com", "www.web.archive.org", "www.secure.gravatar.com", "www.usatoday.com", "www.huffingtonpost.com", "www.stackoverflow.com", "www.fb.com", "www.npr.org", "www.techcrunch.com", "www.wired.com", "www.eepurl.com", "www.reuters.com", "www.arxiv.org", "www.i0.wp.com", "www.finance.yahoo.com", "www.player.vimeo.com", "www.amazon.fr", "www.developer.mozilla.org", "www.tumblr.com", "www.cnbc.com", "www.booking.com", "www.zoom.us", "www.tools.google.com", "www.wa.me", "www.paypal.me", "www.slack.com", "www.adwords.google.com", "www.independent.co.uk", "www.discord.gg", "www.open.spotify.com", "www.tripadvisor.com", "www.gsuite.google.com", "www.behance.net", "www.thinkwithgoogle.com", "www.who.int", "www.cloud.google.com", "www.calendar.google.com", "www.ibm.com", "www.analytics.google.com", "www.lh5.googleusercontent.com", "www.docs.microsoft.com", "www.groups.google.com", "www.mixcloud.com", "www.google.co.in", "www.latimes.com", "www.gist.github.com", "www.download.macromedia.com", "www.dailymail.co.uk", "www.myaccount.google.com", "www.line.me", "www.unsplash.com", "www.addons.mozilla.org", "www.msdn.microsoft.com", "www.storage.googleapis.com", "www.pagead2.googlesyndication.com", "www.cdn.jsdelivr.net", "www.g.co", "www.meetup.com", "www.bloomberg.com", "www.books.google.com", "www.vox.com", "www.imdb.com", "www.arstechnica.com", "www.support.mozilla.org", "www.amazon.co.uk", "www.statista.com", "www.google.ca", "www.nbcnews.com", "www.steemit.com", "www.upload.wikimedia.org", "www.aws.amazon.com", "www.amazon.co.jp", "www.w3schools.com", "www.weibo.com", "www.twitch.tv", "www.moz.com", "www.drupal.org", "www.nypost.com", "www.fr.wikipedia.org", "www.ebay.com", "www.sourceforge.net", "www.godaddy.com", "www.ncbi.nlm.nih.gov", "www.wsj.com", "www.thenextweb.com", "www.cdc.gov", "www.photos.google.com", "www.trends.google.com", "www.opera.com", "www.ftc.gov", "www.j.mp", "www.about.me", "www.gov.uk", "www.i.ytimg.com", "www.squareup.com", "www.google.co.uk", "www.buzzfeed.com", "www.developer.apple.com", "www.pixabay.com", "www.quora.com", "www.cnet.com", "www.google.com.au", "www.myspace.com", "www.theglobeandmail.com", "www.shopify.com", "www.whitehouse.gov", "www.marketingplatform.google.com", "www.kickstarter.com", "www.freelancer.com", "www.evernote.com", "www.gmail.com", "www.static.wixstatic.com", "www.fortune.com", "www.wix.com", "www.dailymotion.com", "www.goodreads.com", "www.podcasts.google.com", "www.google.co.jp", "www.support.cloudflare.com", "www.surveymonkey.com", "www.cbsnews.com", "www.entrepreneur.com", "www.bing.com", "www.engadget.com", "www.lh4.googleusercontent.com", "www.tandfonline.com", "www.cargocollective.com", "www.join.slack.com", "www.api.whatsapp.com", "www.shutterstock.com", "www.pbs.twimg.com", "www.blog.google", "www.ads.google.com", "www.time.com", "www.lifehacker.com", "www.patreon.com", "www.scribd.com", "www.lh6.googleusercontent.com", "www.feeds.feedburner.com", "www.trello.com", "www.ted.com", "www.helpx.adobe.com", "www.researchgate.net", "www.fda.gov", "www.instructables.com", "www.google.de", "www.digitalocean.com", "www.dribbble.com", "www.ssl.gstatic.com", "www.variety.com", "www.indiegogo.com", "www.salesforce.com", "www.uk.linkedin.com", "www.nasa.gov", "www.codepen.io", "www.newyorker.com", "www.abcnews.go.com", "www.snapchat.com", "www.googleadservices.com", "www.sciencedirect.com", "www.ft.com", "www.nature.com", "www.gofundme.com", "www.500px.com", "www.qz.com", "www.dl.dropboxusercontent.com", "www.store.steampowered.com", "www.webmasters.googleblog.com", "www.podcasts.apple.com", "www.oracle.com", "www.m.youtube.com", "www.gartner.com", "www.web.facebook.com", "www.blogs.msdn.com", "www.mailchi.mp", "www.digg.com", "www.skype.com", "www.periscope.tv", "www.google.com.br", "www.googleblog.blogspot.com", "www.de.wikipedia.org", "www.search.google.com", "www.npmjs.com", "www.motherboard.vice.com", "www.aboutads.info", "www.pexels.com", "www.sxsw.com", "www.fiverr.com", "www.code.google.com", "www.cbc.ca", "www.venturebeat.com", "www.hangouts.google.com", "www.mashable.com", "www.upwork.com", "www.bitly.com", "www.music.apple.com", "www.loc.gov", "www.news.google.com", "www.spotify.com", "www.netflix.com", "www.hbr.org", "www.pbs.org", "www.stumbleupon.com", "www.walmart.com", "www.foursquare.com", "www.amazon.ca", "www.tiny.cc", "www.services.google.com", "www.xkcd.com", "www.target.com", "www.developers.facebook.com", "www.nginx.com", "www.canva.com", "www.s7.addthis.com", "www.elmundo.es", "www.i2.wp.com", "www.gizmodo.com", "www.expedia.com", "www.picasaweb.google.com", "www.money.cnn.com", "www.timesofindia.indiatimes.com", "www.yandex.ru", "www.vice.com", "www.de-de.facebook.com", "www.pastebin.com", "www.samsung.com", "www.rt.com", "www.ifttt.com", "www.dev.mysql.com", "www.giphy.com", "www.use.fontawesome.com", "www.houzz.com", "www.developer.android.com", "www.wikipedia.org", "www.speakerdeck.com", "www.productforums.google.com", "www.digitaltrends.com", "www.steamcommunity.com", "www.canada.ca", "www.networkadvertising.org", "www.linktr.ee", "www.calendly.com", "www.lg.com", "www.yahoo.com", "www.last.fm", "www.amazon.es", "www.businesswire.com", "www.deezer.com", "www.dx.doi.org", "www.udemy.com", "www.edition.cnn.com", "www.1.bp.blogspot.com", "www.l.facebook.com", "www.dl.dropbox.com", "www.store.google.com", "www.gumroad.com", "www.azure.microsoft.com", "www.i1.wp.com", "www.ubuntu.com", "www.google.ru", "www.makeuseof.com", "www.raw.githubusercontent.com", "www.airbnb.com", "www.bluehost.com", "www.deviantart.com", "www.pewresearch.org", "www.google.cn", "www.profiles.google.com", "www.login.microsoftonline.com", "www.google.nl", "www.google.es", "www.fastcompany.com", "www.is.gd", "www.telegram.me", "www.mediafire.com", "www.themeforest.net", "www.openstreetmap.org", "www.pewinternet.org", "www.photos.app.goo.gl", "www.buff.ly", "www.events.google.com", "www.link.springer.com", "www.bizjournals.com", "www.residentadvisor.net", "www.golang.org", "www.vogue.com", "www.s3-eu-west-1.amazonaws.com", "www.boardgamegeek.com", "www.un.org", "www.i.redd.it", "www.ad.doubleclick.net", "www.ow.ly", "www.get.google.com", "www.pt.slideshare.net", "www.buffer.com", "www.mozilla.com", "www.chromium.org", "www.patheos.com", "www.ancestry.com", "www.youtube-nocookie.com", "www.barnesandnoble.com", "www.amzn.com", "www.hulu.com", "www.greenpeace.org", "www.spiegel.de", "www.apps.facebook.com", "www.census.gov", "www.app.box.com", "www.windows.microsoft.com", "www.yandex.com", "www.sketchfab.com", "www.avvo.com", "www.searchengineland.com", "www.shareasale.com", "www.abc.net.au", "www.trustpilot.com", "www.plaza.rakuten.co.jp", "www.marriott.com", "www.amazon.it", "www.ca.linkedin.com", "www.uber.com", "www.coursera.org", "www.adage.com", "www.foxnews.com", "www.flic.kr", "www.developer.chrome.com", "www.pl.wikipedia.org", "www.change.org", "www.spreaker.com", "www.filezilla-project.org", "www.flipboard.com", "www.pcworld.com", "www.news.yahoo.com", "www.pitchfork.com", "www.weforum.org", "www.istockphoto.com", "www.9to5mac.com", "www.1drv.ms", "www.uspto.gov", "www.irs.gov", "www.copyblogger.com", "www.example.com", "www.maps.google.co.jp", "www.fonts.google.com", "www.siteground.com", "www.sony.net", "www.ja.wikipedia.org", "www.creativemarket.com", "www.smithsonianmag.com", "www.mobile.twitter.com", "www.acm.org", "www.accenture.com", "www.git-scm.com", "www.bbb.org", "www.cdn.shopify.com", "www.schema.org", "www.psychologytoday.com", "www.redcross.org", "www.nodejs.org", "www.europa.eu", "www.ameblo.jp", "www.gitlab.com", "www.sciencemag.org", "www.hollywoodreporter.com", "www.freshbooks.com", "www.heise.de", "www.chicagotribune.com", "www.epa.gov", "www.anchor.fm", "www.telegram.org", "www.ikea.com", "www.bloglovin.com", "www.eventbrite.co.uk", "www.googlewebmastercentral.blogspot.com", "www.huffpost.com", "www.weebly.com", "www.producthunt.com", "www.google.fr", "www.click.linksynergy.com", "www.vanityfair.com", "www.diigo.com", "www.about.fb.com", "www.ouest-france.fr", "www.namecheap.com", "www.codecanyon.net", "www.appstore.com", "www.tools.ietf.org", "www.rockpapershotgun.com", "www.columbia.edu", "www.constantcontact.com", "www.webmd.com", "www.eur-lex.europa.eu", "www.aliexpress.com", "www.business.facebook.com", "www.material.io", "www.neilpatel.com", "www.discordapp.com", "www.m.me", "www.design.google", "www.francetvinfo.fr", "www.slashgear.com", "www.4shared.com", "www.python.org", "www.amazon.com.au", "www.storify.com", "www.elegantthemes.com", "www.europe1.fr", "www.pwc.com", "www.activecampaign.com", "www.freepik.com", "www.ravelry.com", "www.apnews.com", "www.propublica.org", "www.smugmug.com", "www.inc.com", "www.google.it", "www.zazzle.com", "www.dol.gov", "www.scoop.it", "www.bit.do", "www.yadi.sk", "www.eclipse.org", "www.economist.com", "www.xing.com", "www.hbo.com", "www.nationalgeographic.com", "www.jetbrains.com", "www.iconfinder.com", "www.cafepress.com", "www.strava.com", "www.es.wikipedia.org", "www.kiva.org", "www.doi.org", "www.boredpanda.com", "www.justgiving.com", "www.pinterest.ca", "www.stats.wp.com", "www.zdnet.com", "www.sproutsocial.com", "www.rebrand.ly", "www.ranker.com", "www.a.co", "www.waze.com", "www.nydailynews.com", "www.marketwatch.com", "www.mailchimp.com", "www.livescience.com", "www.flaticon.com", "www.ilpost.it", "www.connect.facebook.net", "www.4.bp.blogspot.com", "www.maps.google.co.nz", "www.s0.wp.com", "www.amazon.com.br", "www.ustream.tv", "www.lulu.com", "www.socialmediatoday.com", "www.market.android.com", "www.rollingstone.com", "www.journals.sagepub.com", "www.zillow.com", "www.xinhuanet.com", "www.affiliate-program.amazon.com", "www.politico.com", "www.pnas.org", "www.laughingsquid.com", "www.access.redhat.com", "www.s.ytimg.com", "www.it.linkedin.com", "www.dashlane.com", "www.copyright.gov", "www.news.bbc.co.uk", "www.mentalfloss.com", "www.zeit.de", "www.google.cz", "www.artsandculture.google.com", "www.startnext.com", "www.launchpad.net", "www.lemonde.fr", "www.sciencedaily.com", "www.edx.org", "www.smashwords.com", "www.philips.co.uk", "www.smile.amazon.com", "www.intel.com", "www.livestream.com", "www.brookings.edu", "www.thehill.com", "www.popsci.com", "www.blog.us.playstation.com", "www.fb.me", "www.academic.oup.com", "www.php.net", "www.wetransfer.com", "www.blog.livedoor.jp", "www.amazon.in", "www.families.google.com", "www.obsproject.com", "www.t.qq.com", "www.android.com", "www.blogs.windows.com", "www.hostgator.com", "www.google.gr", "www.theverge.com", "www.aub.edu.lb", "www.digiday.com", "www.gimp.org", "www.treasury.gov", "www.geocities.com", "www.healthline.com", "www.3.bp.blogspot.com", "www.coinbase.com", "www.gnu.org", "www.abc.com", "www.aljazeera.com", "www.getresponse.com", "www.dev.to", "www.thesun.co.uk", "www.lynda.com", "www.sellfy.com", "www.prnt.sc", "www.help.ubuntu.com", "www.whatsapp.com", "www.www8.hp.com", "www.notion.so", "www.vmware.com", "www.kraken.com", "www.vine.co", "www.marthastewart.com", "www.discord.me", "www.dailycaller.com", "www.mega.nz", "www.media.giphy.com", "www.flavors.me", "www.blip.tv", "www.prntscr.com", "www.journals.plos.org", "www.gravatar.com", "www.discogs.com", "www.europarl.europa.eu", "www.redhat.com", "www.bugs.chromium.org", "www.vanmiubeauty.com", "www.metro.co.uk", "www.google.ie", "www.relapse.com", "www.earth.google.com", "www.bandcamp.com", "www.br.linkedin.com", "www.gopro.com", "www.t.ly", "www.get.adobe.com", "www.news.mit.edu", "www.billboard.com", "www.britannica.com", "www.onlinelibrary.wiley.com", "www.techsmith.com", "www.sfgate.com", "www.redbubble.com", "www.kstatic.googleusercontent.com", "www.repubblica.it", "www.disqus.com", "www.khanacademy.org", "www.skfb.ly", "www.androidauthority.com", "www.accessify.com", "www.tensorflow.org", "www.chrisjdavis.org", "www.box.net", "www.poetryfoundation.org", "www.law.cornell.edu", "www.penguinrandomhouse.com", "www.refinery29.com", "www.starwars.com", "www.fr.linkedin.com", "www.video.google.com", "www.prnewswire.com", "www.springer.com", "www.zalo.me", "www.reverbnation.com", "www.getpocket.com", "www.webdesign.about.com", "www.goo.gle", "www.note.mu", "www.thetimes.co.uk", "www.fao.org", "www.bol.com", "www.yoursite.com", "www.sophos.com", "www.automattic.com", "www.podbean.com", "www.sublimetext.com", "www.archives.gov", "www.vizio.com", "www.thelancet.com", "www.guardian.co.uk", "www.news.nationalgeographic.com", "www.tunein.com", "www.stitcher.com", "www.coinmarketcap.com", "www.amnestyusa.org", "www.docs.wixstatic.com", "www.cambridge.org", "www.validator.w3.org", "www.asus.com", "www.problogger.net", "www.worldbank.org", "www.ticketportal.cz", "www.dw.com", "www.static.googleusercontent.com", "www.www-01.ibm.com", "www.business.google.com", "www.es.linkedin.com", "www.songkick.com", "www.adweek.com", "www.lifehack.org", "www.apple.co", "www.sutterhealth.org", "www.globalnews.ca", "www.breitbart.com", "www.thumbtack.com", "www.eff.org", "www.snopes.com", "www.crowdrise.com", "www.hkrsa.asia", "www.metmuseum.org", "www.pixlr.com", "www.institutvajrayogini.fr", "www.rhtlbn.mosgorcredit.ru", "www.itchyfeetonthecheap.com", "www.reacts.ru", "www.db.tt", "www.udacity.com", "www.we.tl", "www.blogtalkradio.com", "www.audacityteam.org", "www.youcaring.com", "www.hp.com", "www.0.gravatar.com", "www.health.harvard.edu", "www.meta.wikimedia.org", "www.google.ch", "www.nhs.uk", "www.mixi.jp", "www.symfony.com", "www.business.linkedin.com", "www.mtv.com", "www.1.gravatar.com", "www.cyber.law.harvard.edu", "www.bitcointalk.org", "www.esa.int", "www.slate.com", "www.oecd.org", "www.raspberrypi.org", "www.sfexaminer.com", "www.warriorforum.com", "www.g.page", "www.bhphotovideo.com", "www.axios.com", "www.pinterest.co.uk", "www.unesco.org", "www.msn.com", "www.code.visualstudio.com", "www.netbeans.org", "www.blockchain.info", "www.untappd.com", "www.espn.com", "www.cbs.com", "www.michigan.gov", "www.oprah.com", "www.ok.ru", "www.mhlw.go.jp", "www.newegg.com", "www.2.bp.blogspot.com", "www.pipes.yahoo.com", "www.themarthablog.com", "www.yummly.com", "www.bitpay.com", "www.buymeacoffee.com", "www.deepl.com", "www.blogs.scientificamerican.com", "www.audible.com", "www.bitbucket.org", "www.it.wikipedia.org", "www.in.linkedin.com", "www.cse.google.com", "www.use.typekit.net", "www.ru.wikipedia.org", "www.faz.net", "www.adf.ly", "www.cancerresearchuk.org", "www.moma.org", "www.bigthink.com", "www.bandsintown.com", "www.addthis.com", "www.support.office.com", "www.upi.com", "www.zen.yandex.ru", "www.collegehumor.com", "www.lefigaro.fr", "www.gplus.to", "www.themify.me", "www.rtve.es", "www.nba.com", "www.myfitnesspal.com", "www.chris.pirillo.com", "www.eonline.com", "www.science.sciencemag.org", "www.ietf.org", "www.ko-fi.com", "www.otto.de", "www.tkqlhce.com", "www.unity3d.com", "www.blog.naver.com", "www.scratch.mit.edu", "www.freewebs.com", "www.drift.com", "www.postmates.com", "www.geek.com", "www.geocities.jp", "www.beian.gov.cn", "www.meet.google.com", "www.flattr.com", "www.webroot.com", "www.thoughtcatalog.com", "www.google.pl", "www.seroundtable.com", "www.checkpoint.com", "www.form.jotform.com", "www.xbox.com", "www.infusionsoft.com", "www.g1.globo.com", "www.capterra.com", "www.ea.com", "www.hawaii.edu", "www.stock.adobe.com", "www.osha.gov", "www.buzzfeednews.com", "www.jstor.org", "www.adssettings.google.com", "www.iheart.com", "www.en-gb.facebook.com", "www.foodnetwork.com", "www.eventim.de", "www.economictimes.indiatimes.com", "www.presseportal.de", "www.google.se", "www.skillshare.com", "www.puu.sh", "www.examiner.com", "www.stripe.com", "www.firstdata.com", "www.en.advertisercommunity.com", "www.download.microsoft.com", "www.s-media-cache-ak0.pinimg.com", "www.ctt.ec", "www.cia.gov", "www.docker.com", "www.homedepot.com", "www.humblebundle.com", "www.chiark.greenend.org.uk", "www.pandora.com", "www.google.co.za", "www.payhip.com", "www.thingiverse.com", "www.inkscape.org", "www.adobe.ly", "www.css-tricks.com", "www.ja-jp.facebook.com", "www.jamanetwork.com", "www.maps.gstatic.com", "www.wattpad.com", "www.chase.com", "www.cisco.com", "www.bmj.com", "www.louvre.fr", "www.aclu.org", "www.squarespace.com", "www.gum.co", "www.idealo.de", "www.imore.com", "www.rakuten.com", "www.fbi.gov", "www.united.com", "www.google.pt", "www.event.on24.com", "www.nejm.org", "www.colorado.edu", "www.franchising.com", "www.logitech.com", "www.redbull.com", "www.tf1.fr", "www.keep.google.com", "www.hostinger.com", "www.nbc.com", "www.stats.g.doubleclick.net", "www.ign.com", "www.blog.hubspot.com", "www.1.envato.market", "www.accessdata.fda.gov", "www.google.be", "www.help.apple.com", "www.technorati.com", "www.pond5.com", "www.ucl.ac.uk", "www.a2hosting.com", "www.amzn.asia", "www.nicovideo.jp", "www.indiewire.com", "www.mp.weixin.qq.com", "www.codeproject.com", "www.1.usa.gov", "www.gleam.io", "www.space.com", "www.gitter.im", "www.lh5.ggpht.com", "www.leparisien.fr", "www.kotaku.com", "www.img.youtube.com", "www.huffingtonpost.co.uk", "www.agoda.com", "www.flow.microsoft.com", "www.chronicle.com", "www.createspace.com", "www.lh3.ggpht.com", "www.verizon.com", "www.teespring.com", "www.data.worldbank.org", "www.blog.goo.ne.jp", "www.kobo.com", "www.globalsign.com", "www.money.yandex.ru", "www.animoto.com", "www.google.co.nz", "www.privacy.microsoft.com", "www.online.wsj.com", "www.patft.uspto.gov", "www.rottentomatoes.com", "www.cell.com", "www.de.linkedin.com", "www.health.com", "www.geni.us", "www.news.harvard.edu", "www.lesechos.fr", "www.fas.org", "www.smashingmagazine.com", "www.7-zip.org", "www.realvnc.com", "www.news.discovery.com", "www.mlb.com", "www.autodesk.com", "www.popularmechanics.com", "www.forms.office.com", "www.rdio.com", "www.sendspace.com", "www.synology.com", "www.angelfire.com", "www.france24.com", "www.us.battle.net", "www.stjude.org", "www.denverpost.com", "www.japantimes.co.jp", "www.stocktwits.com", "www.blogs.adobe.com", "www.envato.com", "www.airtable.com", "www.ticketmaster.com", "www.ocw.mit.edu", "www.whc.unesco.org", "www.neh.gov", "www.purl.org", "www.nfl.com", "www.monster.com", "www.codex.wordpress.org", "www.dafont.com", "www.tesla.com", "www.lmgtfy.com", "www.olympic.org", "www.technet.microsoft.com", "www.baidu.com", "www.desktop.github.com", "www.overcast.fm", "www.allmusic.com", "www.cdbaby.com", "www.vsco.co", "www.nvidia.com", "www.stadt-bremerhaven.de", "www.faa.gov", "www.blogs.yahoo.co.jp", "www.bild.de", "www.lenovo.com", "www.funnyordie.com", "www.ssl.google-analytics.com", "www.abebooks.com", "www.spectrum.ieee.org", "www.thinkgeek.com", "www.snip.ly", "www.pixiv.net", "www.buzzsprout.com", "www.patents.google.com", "www.orcid.org", "www.cdn.ampproject.org", "www.blog.feedspot.com", "www.vr.google.com", "www.opinionator.blogs.nytimes.com", "www.ericsson.com"}
