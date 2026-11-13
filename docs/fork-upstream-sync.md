# Automated upstream stable release synchronization

The fork automatically polls the [official 3x-ui release](https://github.com/MHSanaei/3x-ui/releases)
every six hours. It can also be run from **Actions → Sync upstream stable
releases → Run workflow**.

**Important for GitHub forks:** GitHub may disable scheduled workflows by
default in a newly forked repository. Visit **Actions → Sync upstream stable
releases**. If the page shows **Enable workflow**, click it once. Then use
**Run workflow** for an initial manual check; the schedule is not active until
the workflow is enabled.

New stable tags are merged into `main` (not force-reset), preserving the
fork's Mihomo ECH converter. If there is a conflict or `go test ./internal/sub`
fails, the workflow stops, leaves `main` unchanged, and opens a GitHub Issue.

On success the workflow dispatches the Release and CI Actions explicitly.
That is required because GitHub workflow pushes made with `GITHUB_TOKEN`
do **not** trigger ordinary `on: push` workflows. The existing Release
workflow publishes the rolling `dev-latest` assets, and another workflow
opens an Issue with the build result and download URL.

The panel's version checker, update script, and bundled `x-ui.sh` menu point
to **Scluzlep/3x-ui** instead of upstream. In the web panel, enable the
**Dev update channel** to compare commit IDs against the fork's
`dev-latest` release. Updates are never deployed to any VPS automatically;
you must approve each installation.

To update from a shell, after backing up the DB, use:

```bash
curl -fL https://raw.githubusercontent.com/Scluzlep/3x-ui/main/update.sh -o /tmp/xui-fork-update.sh
sudo env XUI_UPDATE_TAG=dev-latest bash /tmp/xui-fork-update.sh
```

**Enable Issues in Settings → General → Features → Issues first** (this fork currently has Issues disabled). Otherwise a build result is visible only in Actions, because issue notifications cannot be posted.

Configure GitHub notifications for Issues and Releases in this fork to receive
the `[upstream-sync]` notifications. Scheduled checks are best-effort and
may be delayed by GitHub.
