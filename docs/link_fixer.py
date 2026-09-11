"""
Rewrite hrefs in files that MkDocs publishes via docs/ symlinks.

README.md, CONTRIBUTING.md, and ADOPTERS.md live at the repository root. GitHub
needs `docs/<page>.md` from those files. MkDocs' docs_dir is already `docs/`, so
the same href would look for docs/docs/<page>.md and warn as missing.

docs/crossplane_globalapp.md is a symlink to the example README. Image hrefs
must stay `assets/...` for that README; the nav path needs them prefixed.

Linking strategy:
- Repo-root files: keep `docs/<page>.md` (GitHub). This hook strips the prefix.
- Files under docs/: use same-tree relative links (`local.md`, `resource_ref.md`).
- Targets outside docs/ (adr/, LICENSE, ...): use a GitHub blob or tree URL.
- Symlinked nav pages: keep paths correct for the real file; this hook rewrites
  them for the MkDocs src_path.
"""
import re

# docs/index.md -> ../README.md, docs/CONTRIBUTING.md -> ../CONTRIBUTING.md,
# docs/ADOPTERS.md -> ../ADOPTERS.md
_REPO_ROOT_PAGES = frozenset({"index.md", "CONTRIBUTING.md", "ADOPTERS.md"})
_DOCS_MD_LINK = re.compile(r"\(docs/([a-zA-Z0-9_.-]+\.md)(#[^)]*)?\)")
_CROSSPLANE_NAV = "crossplane_globalapp.md"
_CROSSPLANE_ASSETS = "examples/crossplane/globalapp/assets/"
_CROSSPLANE_ASSET_LINK = re.compile(r"\(assets/([^)]+)\)")


def rewrite_repo_root_docs_links(markdown):
    """Turn (docs/page.md) into (page.md) for MkDocs."""
    return _DOCS_MD_LINK.sub(lambda m: f"({m.group(1)}{m.group(2) or ''})", markdown)


def rewrite_crossplane_nav_assets(markdown):
    """Prefix example-local asset hrefs for the Crossplane nav symlink."""
    return _CROSSPLANE_ASSET_LINK.sub(rf"({_CROSSPLANE_ASSETS}\1)", markdown)


def fix_links(markdown, page, config, files):
    src = page.file.src_path
    if src in _REPO_ROOT_PAGES:
        return rewrite_repo_root_docs_links(markdown)
    if src == _CROSSPLANE_NAV:
        return rewrite_crossplane_nav_assets(markdown)
    return markdown
