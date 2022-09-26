# adr
A tool for managing Architectural/Any Decision Records (ADRs). This tool will enable the management of an ADR repository (a log of decisions recorded in Markdown files) via an intuitive, cross-platform, command line interface.

There are 2 broad use cases:
1. An individual needing to manage a decision repository
2. CI/CD

Things you may want to do with this tool:
1. Create, update, and supersede records (the functionality of [adr-tools](https://github.com/npryce/adr-tools)).
2. "Lint" your records in ci to ensure they fit the declared pattern. Not everybody will always choose to use this tool, but you should be able to ensure there is no functional difference by using this tool to test in continuous integration. These are not currently planned.
3. Generate reports. Not sure if we're going to support this, some statistics could be interesting but maybe not that useful in the end.

There are 3 location-based scenarios where you may want to manage ADRs:
1. In some software project. This is the default that ADRs were intended for. The default location is the `docs/decisions` subdirectories under the root of your project. This scenario is our default.
2. In some centralized place where the scope is greater than a specific project. For example, in traditional organizations it is typical to have many projects/applications falling under the purview of one architect (or a small group). In this case it can be desirable to have a central decision repository that does not contain code or other artifacts, so the default `docs/decisions` subdirectory structure makes less sense and should be configurable.
3. In some software project but the `docs/decisions` default will not work or does not make sense for some reason. Another scenario indicating that the storage location should be configurable. 

## Planned Features

- [x] Create a new ADR based on the Nygard format
- [x] Store configuration per ADR repository
- [x] Link an ADR to another ADR within the same repository (and back)
- [ ] Supersede an existing ADR with another ADR
- [x] Update the status of an existing ADR
- [x] Configure custom directory where ADRs are maintained

## Potential Future Features
- [ ] Customizable ADR format with Nygard by default (partial due to config file)
- [ ] Configure custom statuses for your ADRs
- [ ] Lint ADRs
- [ ] Default global init values (used when initializing the repo)
