# Google Analytics Spam Control

Command-line utility for blocking referer spam from your Google Analytics accounts

Google Analytics [referer spam](https://en.wikipedia.org/wiki/Referer_spam) is pain.
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

- [Referer spam](https://en.wikipedia.org/wiki/Referer_spam)
- [List of spam domains](https://github.com/ddofborg/analytics-ghost-spam-list)
- [Online Tool: Analytics Referrer/Ghost Spam Blocker](https://www.adwordsrobot.com/en/tools/ga-referrer-spam-killer)
- [Google Analytics Account API](https://developers.google.com/analytics/devguides/config/mgmt/v3/mgmtReference/management/accounts/list)
- [Google Analytics Filter API](https://developers.google.com/analytics/devguides/config/mgmt/v3/mgmtReference/management/filters)
- [Google Analytics Filter Expressions](https://developers.google.com/analytics/devguides/reporting/core/v3/reference#filters)
- [Google Analytics Data Management](https://developers.google.com/analytics/devguides/config/mgmt/v3/data-management)
- [Google Analytics Profile Filter Links](https://developers.google.com/analytics/devguides/config/mgmt/v3/mgmtReference/management/profileFilterLinks)
