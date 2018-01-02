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
