# adr
A tool for managing Architectural/Any Decision Records (ADRs). This tool will enable the management of an ADR repository (a log of decisions recorded in Markdown files) via an intuitive, cross-platform, command line interface.

## Current Features
1. Default Nygard template for simple usage. No need to initialize the repo or maintain `.adr.yaml`
2. Simple configuration via `.adr.yaml`
   1. Repository path (i.e. where the ADR documents for this repo are stored)
   2. Title template (i.e. how your documents are named)
   3. Body template (i.e. how your documents will look by default when created)
3. Link ADRs together
   1. Freeform linking w/ individual messages for link and backlink
   2. Superseding, a special case of linking

## Features under consideration
- Initialize w/ first decision to record decisions
- More pre-defined ADR formats from which to choose
- Status enforcement (e.g. choose from predefined list)
- Lint ADRs (e.g. ensure title fits the format, required sections present, etc)
- Default global init values for folks managing multiple repos (used when initializing the repo)
- Reports (e.g. status ratios, lead time from proposal to acceptance, etc)
