victims-bot
===========

GitHub bot for victims. Will take care of handling submissions
through GitHub merging.

Building
========

For all targets run ``make help``

Dynamic
-------
1. Install dependencies via ``govendor`` by running: ``make deps``
2. Run ``make victims-bot``
3. Execute ``./victims-bot -secret=$VICTIMS_GITHUB_SECRET``

Static
------
1. Install dependencies via ``govendor`` by running: ``make deps``
2. Run ``make static-victims-bot``
3. Execute ``./victims-bot -secret=$VICTIMS_GITHUB_SECRET``


Configuration
=============

| CLI             | Environment Variable        | Description                |
|-----------------|-----------------------------|----------------------------|
| bind            | VICTIMS_BOT_BIND            | Host:Port to listen on     |
| git-repo        | VICTIMS_BOT_GIT_REPO        | Repo to push/modify/pull   |
| github-username | VICTIMS_BOT_GITHUB_USERNAME | GitHub username of the bot |
| github-password | VICTIMS_BOT_GITHUB_PASSWORD | GitHub password to push    |
| secret          | VICTIMS_BOT_SECRET          | GitHub secret for webhooks |


Deployment
==========

* [OpenShift](/deployment/openshift/)
