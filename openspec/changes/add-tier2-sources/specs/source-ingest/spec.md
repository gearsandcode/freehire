## ADDED Requirements

### Requirement: Personio, Breezy, Pinpoint, Rippling, BambooHR, and Join.com are registered providers

The system SHALL register `personio`, `breezy`, `pinpoint`, `rippling`, `bamboohr`, and
`join.com` adapters so boards on these platforms can be listed in `sources.yml`. Each
adapter SHALL yield the normalized job shape (at least title, url, location, remote flag,
description, and the platform's native posting id) with the `description` as sanitized
HTML assembled from the platform's authoritative HTML field(s), consistent with the
existing adapters. An adapter whose list endpoint omits the description SHALL fetch each
posting's detail with bounded concurrency rather than yield an empty body.

#### Scenario: Personio XML feed is crawled in one request

- **WHEN** `sources.yml` lists a board with provider `personio`
- **THEN** the adapter fetches the board's `…jobs.personio.com/xml` feed in one request and
  yields each `<position>` with a sanitized HTML description assembled from its inline
  `jobDescriptions`

#### Scenario: Breezy board carries the body inline

- **WHEN** a `breezy` board is crawled
- **THEN** the adapter fetches the board's JSON in one request and yields each posting with
  a sanitized HTML description from the inline body field

#### Scenario: Pinpoint board carries the body inline

- **WHEN** a `pinpoint` board is crawled
- **THEN** the adapter fetches the board's `…/postings.json` in one request and yields each
  posting with a sanitized HTML description from the inline body field

#### Scenario: Rippling posting gains its description from detail

- **WHEN** a `rippling` board is crawled
- **THEN** the adapter fetches the board's job list and, per posting, fetches its detail with
  bounded concurrency to obtain the description, still yielding the normalized job shape

#### Scenario: BambooHR posting gains its description from detail

- **WHEN** a `bamboohr` board is crawled
- **THEN** the adapter fetches `…/careers/list` and, per posting, fetches `…/careers/{id}/detail`
  with bounded concurrency to obtain the description, still yielding the normalized job shape

#### Scenario: Join.com postings are fetched from its captured feed

- **WHEN** a `join.com` board is crawled
- **THEN** the adapter fetches the company's postings from join.com's public feed and yields
  each with a sanitized HTML description, still yielding the normalized job shape

#### Scenario: A board with no open postings yields no jobs without error

- **WHEN** any of these providers' feeds returns an empty posting list for a configured board
- **THEN** the adapter yields zero jobs and returns no error, so the board is simply skipped
