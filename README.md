# Google Analytics Spam Control

Command-line utility for blocking referrer spam from your Google Analytics accounts

Google Analytics [referrer spam](https://en.wikipedia.org/wiki/Referer_spam) is pain.
**ga-spam-control** is a small command-line utility that helps you to keep your Google Analytics spam filters up-to-date.

## Usage

```bash
ga-spam-control <action> [..options]
```

Display the current **status** of your spam filters:

```bash
ga-spam-control status
```

**Update** your spam filters:

```bash
ga-spam-control update
```

**Remove** all spam filters:

```bash
ga-spam-control remove
```

## Related Resources

### Referrer Spam

- [What is referrer spam?](https://en.wikipedia.org/wiki/Referer_spam)

### Lists of Referrer Spam Domains

There are multiple curated lists of referrer spam domains out there that you can use to create filters for your analytics accounts.

- [Analytics ghost spam list](https://github.com/ddofborg/analytics-ghost-spam-list)
- [Piwik Referrer spam blacklist](https://github.com/piwik/referrer-spam-blacklist)
- [Referrer Spam Blocker Blacklist](https://referrerspamblocker.com/blacklist)
- [Stevie Ray: apache-nginx-referral-spam-blacklist](https://github.com/Stevie-Ray/apache-nginx-referral-spam-blacklist)
- [My own list of referral spam domains](spam-domains/referrer-spam-domains.txt)

### Other Spam Blocker Tools

ga-spam-control is not the first and not the only tool that helps you to block referrer spam from your Google Analytics accounts.

- [Online Tool: Analytics Referrer/Ghost Spam Blocker](https://www.adwordsrobot.com/en/tools/ga-referrer-spam-killer)
- [Spam Filter Installer](http://www.simoahava.com/spamfilter/)
- [Referrer Spam Blocker](https://referrerspamblocker.com/)

### Google Analytics: Segments

Filters **prevent** referrer spam from getting into your Google Analytics accounts.
But filters don't help you with referrer spam that already reached your reports. In order to filter this spam out you can use segments that filter out the spammy traffic:

- [Analytics Spam Blocker ](https://www.google.com/analytics/gallery/#posts/search/%3F_.tab%3DMy%26_.sort%3DDATE%26_.start%3D0%26_.viewId%3DgyNgK6N3R6iK-UphdU8M6w/)

### Google Analytics: Bot and Spider Filtering

Google Analytics has a setting to block bots and spiders from your Google Analytics reports.

1. Goto `Google Analytics > Admin > Account > Property > View > View Settings`
2. Goto `Bot Filtering`
3. Check `Exclude all hits from known bots and spiders`

This feature is not advertised much by Google. The only time it officially got mentioned by is in a Google Plus post: [Google Analytics - Introducing Bot and Spider Filtering](https://plus.google.com/+GoogleAnalytics/posts/2tJ79CkfnZk).

I am not yet sure if this flag does the trick. One would assume that is would be easy for Google to exclude all referrer spam and block the stupid spammers once and for all.

### Google Analytics: API

- [Google Analytics Account API](https://developers.google.com/analytics/devguides/config/mgmt/v3/mgmtReference/management/accounts/list)
- [Google Analytics Filter API](https://developers.google.com/analytics/devguides/config/mgmt/v3/mgmtReference/management/filters)
- [Google Analytics Filter Expressions](https://developers.google.com/analytics/devguides/reporting/core/v3/reference#filters)
- [Google Analytics Data Management](https://developers.google.com/analytics/devguides/config/mgmt/v3/data-management)
- [Google Analytics Profile Filter Links](https://developers.google.com/analytics/devguides/config/mgmt/v3/mgmtReference/management/profileFilterLinks)
