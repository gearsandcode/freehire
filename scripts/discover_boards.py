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

import sys
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
