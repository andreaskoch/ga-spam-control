# Google Analytics Spam Control

Command-line utility for automating the fight against Google Analytics referral spam

Google Analytics [referrer spam](https://en.wikipedia.org/wiki/Referer_spam) is huge and never ending pain.
There are hundreds of known referrer spam domains and every other day a new one pops up. And the only way to keep the spammers from skewing your web analytics reports is to their stupid domain names – one by one.

![ga-spam-control logo](files/assets/ga-spam-control-logo-300x200.png)

**ga-spam-control** is a small command-line utility that keeps your Google Analytics spam filters up-to-date, automatically.

## How does ga-spam-control work?

**ga-spam-control** creates filters for your Google Analytics accounts that block known referrer spam domains from your analytics reports and keeps these filter up-to-date.

To always protect your analytics reports from annoying false entries ga-spam-control **combines multiple community-maintained lists** of known spam domains:

- [ddofborgs' Analytics Ghost Spam List](https://github.com/ddofborg/analytics-ghost-spam-list)
- [Stevie Rays'  apache-nginx-referral-spam-blacklist](https://github.com/Stevie-Ray/apache-nginx-referral-spam-blacklist)
- [Piwik Referrer spam blacklist](https://github.com/piwik/referrer-spam-blacklist)

This gives you the ability to completely automate your spam protection process. Just let ga-spam-control check your Google Analytics accounts daily for new spam. And when it detects new spam; update your filters.

## Available Commands

The command line utility provides the following actions.

**Spam Control Filter Actions**

In order to protect your Google Analytics account from spam **ga-spam-control** creates filters which blocks known referrer spam domains from your analytics reports. These are the commands that help you to review and update your spam filters:

- **filters status** displays the spam-control status of all your accounts or for a specific account
- **filters update** creates or updates the spam-control filters for a specific account
- **filters remove** removes all previously created spam-control filters from an account

**Referrer Spam Domains Actions**

The basis for the spam filters is an up-to-date list of known referrer spam domains. And with these commands you can review and update the spam-domain lists:

- **domains list** prints a list of all currently known referrer spam domains
- **domains update** downloads the latest referrer spam domain name lists and updates your local list of known referrer spam domains
- **domains find** allows you to manually review the last `n` days of analytics data and mark domain names as spam

Which domains are currently considered spam is stored in the `~/.ga-spam-control/spam-domains/community.txt` and `~/.ga-spam-control/spam-domains/personal.txt`.

## Using ga-spam-control

```bash
ga-spam-control <command> [<args> ...]
```

### Help

Print information about the available actions:

```bash
ga-spam-control help
```

Print detailed help information about the different arguments and flags of a specific action:

```bash
ga-spam-control help <actionname>
```

### Authorizing ga-spam-control to access your Google Analytics accounts

The first time you perform an action, you will be displayed an oAuth authorization dialog.
If you permit the requested rights the authentication token will be stored in your home directory (`~/.ga-spam-control/credentials.json`).

To sign out you can either delete the file or de-authorize the "Google Analytics Spam Control" app in your Google App Permissions at https://security.google.com/settings/security/permissions.

### Get your spam-Control status

Display the current spam-control **status** for all accounts that you have access to:

```bash
ga-spam-control filters status
```

Display the spam-control status in a parseable format:

```bash
ga-spam-control filters status --quiet
```

Display the current spam-control **status** for a specific Google Analytics account:

```bash
ga-spam-control filters status <accountID>
```

### Install or update filters

Create or update the spam-control filters of a given Google Analytics account:

```bash
ga-spam-control filters update <accountID>
```

### Uninstall filters

Remove the spam-control filters of a given Google Analytics account:

```bash
ga-spam-control filters remove <accountID>
```

### List all known spam domains

Print a list of your known referrer spam domains names (community & personal):

```bash
ga-spam-control domains list
```

### Update your list of known spam domains

Update your local community list of known referrer spam domain names:

```bash
ga-spam-control domains update
```

### Find new spam domains

Find referrer spam domain names in your Google Analtics data. Review the hostnames of the last `n` days of one of your Google Analytics accounts and mark those which you consider spam. All marked domain names will be added to your personal referrer spam list:

```bash
ga-spam-control domains find <accountID> <numberOfDaysToLookBack>
```

By default ga-spam-control will use the last 90 days of analytics data. But if you want to review less or more days you can specify the number of days yourself.

## Installation

The command-line package is [github.com/andreaskoch/ga-spam-control/cli](cli/main.go). You can clone the repository or install it with `go get github.com/andreaskoch/ga-spam-control` and then run the [make.go](make.go) script:

```bash
go run make.go -test
go run make.go -install
go run make.go -crosscompile
```

Or with **make**:

```
make test
make install
make crosscompile
```

## Licensing

ga-spam-control is licensed under the Apache License, Version 2.0.
See [LICENSE](LICENSE) for the full license text.

## Roadmap

Ideally Google would just include a spam-protection into Google Analytics but until then here are some ideas for additional features and possible improvements:

- Make remote spam domain providers configurable
- Publish spam domain names that you found in your Google Analytics accounts back to the community lists.
- Populate my own list of known referrer spam domains with the results from the `find-spam-domains` action.
  - Automatic daily upload from the ga-spam-control clients
  - Review of the additions by trusted community members or by a tool which checks the listed website
- Create and update a "No Referrer Spam" segment and update it during the normal update process.
Unfortunately I will need Google to add create and update support to the Google Analytics API for this to work (see: [analytics-issues - Issue 174: Create Advanced Segment and Customized Report Through API](https://code.google.com/p/analytics-issues/issues/detail?id=174)).
- Until Google supports segment creation via the API I ga-spam-control can at least print the necessary segment content to support manual editing of spam segments.
- Use machine learning to automatically identify new referrer spam.
Earlier versions of ga-spam-control already used a machine learning model. But unfortunately I could only train the model to detect new referrer spam for a single website - the model did not work well enough when I applied it to websites with different usage patterns.
- Other options for detecting referrer spam automatically
  - Correlate analytics data with web server logs to identify referrer spam
  - Do a word analysis of the referrer site and use regular e-mail techniques to identify spam sites

Let me know if you have other ideas, or if want one of the features implemented next.

## Related Resources

### Referrer Spam

- [What is referrer spam?](https://en.wikipedia.org/wiki/Referer_spam)
- [Google Analytics Help Forum - Referral Spam Traffic](https://www.en.advertisercommunity.com/t5/Referral-Spam-Traffic/bd-p/Referral_Spam_Traffic)

### Lists of Referrer Spam Domains

There are multiple curated lists of referrer spam domains out there that you can use to create filters for your analytics accounts.

- [Analytics Ghost Spam List](https://github.com/ddofborg/analytics-ghost-spam-list)
- [Stevie Ray: apache-nginx-referral-spam-blacklist](https://github.com/Stevie-Ray/apache-nginx-referral-spam-blacklist)
- [Piwik Referrer spam blacklist](https://github.com/piwik/referrer-spam-blacklist)
- [Referrer Spam Blocker Blacklist](https://referrerspamblocker.com/blacklist)
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
