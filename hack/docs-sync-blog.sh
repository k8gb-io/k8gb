#!/usr/bin/env bash

set -euo pipefail

die() {
  echo "ERROR: $*" >&2
  exit 1
}

require_cmd() {
  command -v "$1" >/dev/null 2>&1 || die "$1 is required"
}

ghpages_dir="${1:-${GHPAGES_DIR:-}}"
site_url_base="${SITE_URL_BASE:-https://k8gb.io}"
site_url_base="${site_url_base%/}"

[[ -n "$ghpages_dir" ]] || die "GHPAGES_DIR is required"
[[ -d "$ghpages_dir" ]] || die "GHPAGES_DIR does not exist: $ghpages_dir"
[[ -d docs/blog ]] || die "docs/blog not found; run from the source checkout"

require_cmd find
require_cmd mkdocs
require_cmd yq

versions=()
while IFS= read -r version; do
  versions+=("$version")
done < <(find "$ghpages_dir" -mindepth 1 -maxdepth 1 -type d -name 'v*' \
  -exec basename {} \; | sort -V)

if [[ "${#versions[@]}" -eq 0 ]]; then
  echo "WARN: No versioned directories found to sync"
  exit 0
fi

tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT

echo "Syncing blog posts across all deployed versions..."
for version in "${versions[@]}"; do
  site_dir="$tmp_dir/$version"
  echo "Building blog for $version"

  yq eval ".site_url = \"$site_url_base/$version\" | .extra.version.provider = \"mike\"" mkdocs.yml \
    | mkdocs build --config-file - --site-dir "$site_dir"

  [[ -d "$site_dir/blog" ]] || die "mkdocs did not produce blog output for $version"

  rm -rf "$ghpages_dir/$version/blog"
  mkdir -p "$ghpages_dir/$version"
  cp -R "$site_dir/blog" "$ghpages_dir/$version/blog"
done

echo "Successfully synced blog to ${#versions[@]} version(s)"
