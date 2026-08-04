"""
Simple hook functions to fix docs links for MkDocs compatibility.
"""
import re

def fix_links(markdown, page, config, files):
    """
    Fix docs/ links in markdown content for MkDocs compatibility
    """
    if page.file.src_path == 'index.md':
        # Replace (docs/filename.md) with (filename.md) for the index page
        markdown = re.sub(r'\(docs/([a-zA-Z0-9_.-]+\.md)\)', r'(\1)', markdown)
    if page.file.src_path == 'crossplane_globalapp.md':
        # Symlink to examples/crossplane/globalapp/README.md; rewrite sibling asset paths
        # so images resolve from the nav page path at docs/crossplane_globalapp.md.
        markdown = markdown.replace(
            '](assets/',
            '](examples/crossplane/globalapp/assets/',
        )
    return markdown
