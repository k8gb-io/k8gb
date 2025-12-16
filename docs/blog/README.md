# K8GB Blog

This directory contains the K8GB blog powered by the Material for MkDocs blog plugin.

## Structure

```
blog/
├── index.md          # Blog homepage
├── .authors.yml      # Author definitions
├── posts/            # Blog posts directory
│   ├── *.md         # Individual blog posts
└── README.md        # This file
```

## Writing Blog Posts

### Creating a New Post

1. Create a new Markdown file in `docs/blog/posts/`
2. Use a descriptive filename (e.g., `my-awesome-feature.md`)
3. Add the required frontmatter (see template below)

### Post Template

```markdown
---
date: 
  created: YYYY-MM-DD
  updated: YYYY-MM-DD  # Optional
authors:
  - author-id          # Must match .authors.yml
categories:
  - Category Name      # Must match allowed categories
tags:
  - tag1
  - tag2
readtime: 5            # Optional: override calculated reading time
pin: true              # Optional: pin to top of blog
draft: true            # Optional: mark as draft
---

# Your Post Title

Brief description that appears in the blog index.

<!-- more -->

Full post content goes here...
```

### Frontmatter Fields

- **date.created** (required): Publication date in YYYY-MM-DD format
- **date.updated** (optional): Last update date
- **authors** (required): List of author IDs from `.authors.yml`
- **categories** (optional): Must be from allowed list in `mkdocs.yml`
- **tags** (optional): Free-form tags for content organization
- **readtime** (optional): Override calculated reading time in minutes
- **pin** (optional): Pin post to top of blog index
- **draft** (optional): Mark as draft (excluded from builds)

### Allowed Categories

Current allowed categories (defined in `mkdocs.yml`):
- Releases
- Tutorials  
- Community
- Technical Deep Dive
- Integrations

### Authors

Authors are defined in `.authors.yml`. To add a new author:

```yaml
authors:
  your-github-username:
    name: Your Full Name
    description: Your role/description
    avatar: https://github.com/your-github-username.png
```

### Content Guidelines

1. **Use descriptive titles** - Help readers understand what they'll learn
2. **Add a clear excerpt** - The content before `<!-- more -->` appears in the index
3. **Include code examples** - Use proper syntax highlighting
4. **Link to relevant docs** - Help readers find additional information
5. **Add tags** - Help with content discovery
6. **Use appropriate categories** - Keep content organized

### Images and Assets

- Place images in `docs/images/blog/` directory
- Use relative paths: `![Alt text](../../images/blog/my-image.png)`
- Optimize images for web (compress, appropriate formats)

### Code Blocks

Use fenced code blocks with language specification:

```yaml
apiVersion: k8gb.absa.oss/v1beta1
kind: Gslb
metadata:
  name: example
spec:
  strategy:
    type: roundRobin
```

### Internal Links

Link to other documentation:
- `[Installation Guide](../../tutorials.md)`
- `[Components](../../components.md)`

### External Links

Always open in new tabs for external links:
- `[K8GB GitHub](https://github.com/k8gb-io/k8gb){:target="_blank"}`

## Local Development

To preview your blog posts locally:

```bash
# Install dependencies
pip install mkdocs mkdocs-material mkdocs-git-revision-date-localized-plugin mkdocs-simple-hooks

# Serve locally
mkdocs serve

# View at http://localhost:8000/blog/
```

## Publishing

Blog posts are automatically published when:
1. Merged to the `master` branch
2. The `draft: true` flag is removed
3. GitHub Actions builds and deploys to GitHub Pages

## Best Practices

1. **Review before publishing** - Use draft mode for work-in-progress
2. **Keep posts focused** - One main topic per post
3. **Update existing posts** - Use `date.updated` when making significant changes
4. **Cross-reference content** - Link to related documentation and posts
5. **Engage with community** - Encourage comments and feedback via Slack

## Examples

See existing posts in the `posts/` directory for examples of:
- Release announcements
- Technical deep dives
- Tutorial content
- Community stories
- Integration guides

## Questions?

- Join [#k8gb on CNCF Slack](https://cloud-native.slack.com/archives/C021P656HGB)
- Open an issue on [GitHub](https://github.com/k8gb-io/k8gb/issues)
- Check the [Material for MkDocs blog documentation](https://squidfunk.github.io/mkdocs-material/plugins/blog/)