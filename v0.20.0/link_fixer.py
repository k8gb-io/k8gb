"""
Simple hook function to fix docs/ links in README for MkDocs compatibility
"""
import re

def fix_links(markdown, page, config, files):
    """
    Fix docs/ links in markdown content for MkDocs compatibility
    """
    if page.file.src_path == 'index.md':
        # Replace (docs/filename.md) with (filename.md) for the index page
        markdown = re.sub(r'\(docs/([a-zA-Z0-9_.-]+\.md)\)', r'(\1)', markdown)
    return markdown
