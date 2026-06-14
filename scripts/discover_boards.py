#!/usr/bin/env python3
"""Query-driven open-web ATS board discovery.

Run a free-text query across search channels (DuckDuckGo, Google CSE, GitHub code
search, Common Crawl), extract ATS board URLs, dedup against sources/*.yml,
validate each board live, and print (or --write) ready-to-paste YAML.

Usage:
    python3 scripts/discover_boards.py --query "fintech berlin" \
            --provider ashby,lever --channel ddg,github,google,cc [--write] [--limit N]

Stdlib only; the github channel shells out to `gh`; google needs GOOGLE_CSE_KEY/_CX.
"""

from __future__ import annotations

import re
import sys
import urllib.parse
import urllib.request
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent))
from ats_boards import VALIDATORS  # noqa: E402

# Provider -> the ATS host to put in a `site:` search / Common Crawl prefix.
PROVIDER_HOSTS = {
    "greenhouse": "job-boards.greenhouse.io",
    "lever": "jobs.lever.co",
    "ashby": "jobs.ashbyhq.com",
    "smartrecruiters": "jobs.smartrecruiters.com",
    "workable": "apply.workable.com",
    "recruitee": "recruitee.com",
    "bamboohr": "bamboohr.com",
    "breezy": "breezy.hr",
    "personio": "jobs.personio.com",
    "teamtailor": "teamtailor.com",
}

# DuckDuckGo's HTML endpoint rejects unusual UAs; use a browser-like one.
BROWSER_UA = (
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 "
    "(KHTML, like Gecko) Chrome/124.0 Safari/537.36"
)


def get_text(url: str, timeout: int = 25) -> str:
    """GET a URL with a browser UA, returning decoded text ('' on failure)."""
    req = urllib.request.Request(url, headers={"User-Agent": BROWSER_UA})
    try:
        with urllib.request.urlopen(req, timeout=timeout) as r:
            return r.read().decode("utf-8", "replace")
    except Exception as e:
        print(f"  ! GET failed ({url[:60]}...): {e}", file=sys.stderr)
        return ""


def parse_ddg_html(html: str) -> set[str]:
    """Extract target URLs from DuckDuckGo HTML results (unwrap ?uddg= redirects)."""
    return {urllib.parse.unquote(m) for m in re.findall(r"uddg=([^\"&]+)", html)}


def channel_ddg(host: str, query: str, limit: int) -> set[str]:
    """site:<host> <query> via DuckDuckGo HTML -> raw target URLs."""
    q = urllib.parse.quote(f"site:{host} {query}")
    html = get_text(f"https://html.duckduckgo.com/html/?q={q}")
    urls = parse_ddg_html(html)
    return set(list(urls)[:limit]) if limit else urls
