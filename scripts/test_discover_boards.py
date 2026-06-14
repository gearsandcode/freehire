#!/usr/bin/env python3
"""Plain-assert tests for discover_boards (no pytest). Run: python3 scripts/test_discover_boards.py"""

import sys
import traceback
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent))
import ats_boards  # noqa: E402
import discover_boards as d  # noqa: E402


def test_provider_hosts_subset_of_validators():
    assert set(d.PROVIDER_HOSTS) <= set(ats_boards.VALIDATORS), \
        "every discoverable provider must have a validator"
    assert d.PROVIDER_HOSTS["ashby"] == "jobs.ashbyhq.com"


def _run():
    fns = [v for k, v in sorted(globals().items()) if k.startswith("test_") and callable(v)]
    failed = 0
    for fn in fns:
        try:
            fn()
            print(f"ok   {fn.__name__}")
        except Exception:
            failed += 1
            print(f"FAIL {fn.__name__}")
            traceback.print_exc()
    print(f"\n{len(fns) - failed}/{len(fns)} passed")
    return 1 if failed else 0


if __name__ == "__main__":
    raise SystemExit(_run())
