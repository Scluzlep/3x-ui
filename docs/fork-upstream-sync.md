# Stable upstream synchronization for the ECH fork

The fork follows official MHSanaei/3x-ui **stable** releases. Every six hours
the sync workflow checks the latest official tag, merges it into `main`, and
runs `go test ./internal/sub`. If the merge or tests fail, it opens an Issue
without overwriting `main`.

A successful merge triggers `Release 3X-UI` with its version tag, builds
seven Linux packages and Windows amd64, and publishes a **normal, non-prerelease**
GitHub Release (`v3.11.0`, `v3.12.0`, ...). The published tag points to
the fork commit containing both upstream changes and the ECH patch. A version
is published once and is never silently replaced. GitHub Issues provide
success/error notifications.

To install, use the fork's **stable** release page:
https://github.com/Scluzlep/3x-ui/releases/latest

The web-panel updater is already pointed at Scluzlep/3x-ui.
Leave the **Dev Channel OFF**. The normal update button and the bundled
`x-ui.sh` menu use this fork's latest stable Release, not `dev-latest`.

If you installed a prior dev-tagged fork binary for the same base version,
the existing stable-channel UI may incorrectly report "up-to-date" because it
compares only the base version number. Install the first fork stable release
manually once; subsequent official version bumps will be detected normally.

No workflows install updates to the VPS automatically.
