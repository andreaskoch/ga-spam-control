# Google Analytics Spam Control

Command-line utility for blocking referrer spam from your Google Analytics accounts

Google Analytics [referrer spam](https://en.wikipedia.org/wiki/Referer_spam) is pain.
**ga-spam-control** is a small command-line utility that helps you to keep your Google Analtics spam filters up-to-date.

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

- [Analytics ghost spam list](https://github.com/ddofborg/analytics-ghost-spam-list)
- [Piwik Referrer spam blacklist](https://github.com/piwik/referrer-spam-blacklist)
- [Referrer Spam Blocker Blacklist](https://referrerspamblocker.com/blacklist)
- [Stevie Ray: apache-nginx-referral-spam-blacklist](https://github.com/Stevie-Ray/apache-nginx-referral-spam-blacklist)
- [My own list of referral spam domains](spam-domains/referrer-spam-domains.txt)

### Tools

- [Online Tool: Analytics Referrer/Ghost Spam Blocker](https://www.adwordsrobot.com/en/tools/ga-referrer-spam-killer)
- [Spam Filter Installer](http://www.simoahava.com/spamfilter/)
- [Referrer Spam Blocker](https://referrerspamblocker.com/)

### Google Analytics Segments

- [Analytics Spam Blocker ](https://www.google.com/analytics/gallery/#posts/search/%3F_.tab%3DMy%26_.sort%3DDATE%26_.start%3D0%26_.viewId%3DgyNgK6N3R6iK-UphdU8M6w/)

### Google Analytics API

- [Google Analytics Account API](https://developers.google.com/analytics/devguides/config/mgmt/v3/mgmtReference/management/accounts/list)
- [Google Analytics Filter API](https://developers.google.com/analytics/devguides/config/mgmt/v3/mgmtReference/management/filters)
- [Google Analytics Filter Expressions](https://developers.google.com/analytics/devguides/reporting/core/v3/reference#filters)
- [Google Analytics Data Management](https://developers.google.com/analytics/devguides/config/mgmt/v3/data-management)
- [Google Analytics Profile Filter Links](https://developers.google.com/analytics/devguides/config/mgmt/v3/mgmtReference/management/profileFilterLinks)
