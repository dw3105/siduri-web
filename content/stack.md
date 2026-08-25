# Stack

The tools I actually use to make and publish this site.

## Build

- **Go 1.26** owns loading, routing, rendering, and the build binary.
- **templ** generates the type-safe HTML components. Generated Go files are committed so a clean checkout has a visible build boundary.
- **Markdown** is the content format. A small repository renderer turns it into HTML; posts use a closed vocabulary of tags.

## Runtime

- **Static files** are the product: HTML, CSS, feed data, and images can be served by an asset host.
- **CSS** is plain and local. There is no UI framework, analytics beacon, or third-party script in the first build.
- **Git** is the review surface. A draft remains a draft until a person changes that decision.

## What I watch

Build time, output diffs, dependency count, and the places where an agent made a plausible but wrong assumption. The stack is allowed to change; the human gate is not.
