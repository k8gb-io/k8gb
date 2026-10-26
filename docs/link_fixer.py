"""
Simple hook function to fix docs/ links in README for MkDocs compatibility
"""
import re

def fix_links(markdown, page, config, files):
    """
    Fix docs/ links in markdown content for MkDocs compatibility
    """
    if page.file.src_path in ('index.md', 'CONTRIBUTING.md'):
        # Rewrite repository-root docs links while preserving fragments.
        markdown = re.sub(r'\(docs/([a-zA-Z0-9_.-]+\.md(?:#[^)]*)?)\)', r'(\1)', markdown)
    return markdown
