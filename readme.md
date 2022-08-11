# adr
A tool for managing Architectural Decision Records (ADRs). This tool will enable the management of an ADR repository (a log of decisions recorded in Markdown files) via an intuitive command line interface.

There are 3 broad use cases:
1. Create, update, and supersede records (the functionality of [adr-tools](https://github.com/npryce/adr-tools)).
2. "Lint" your records in ci to ensure they fit the declared pattern. Not everybody will always choose to use this tool, but you can ensure there is no functional difference by using this tool to test in continuous integration
3. Generate reports

There are 3 location-based scenarios where you may want to manage ADRs:
1. In some software project. This is the default that ADRs were intended for. The default location is the `docs/decisions` subdirectories under the root of your project.
2. In some centralized place where the scope is greater than a specific project. For example, in traditional organizations it is typical to have many projects/applications falling under the purview of one (or a small group of) architect. In this case it can be desirable to have a central decision repository that does not contain code or other artifacts, so the default `docs/decisions` subdirectory structure makes less sense and should be configurable.
3. In some software project but the `docs/decisions` default will not work or does not make sense for some reason. Another scenario indicating that the storage location should be configurable.

## Features

- [ ] Create a new ADR
- [ ] Link an ADR to another ADR within the same repository
- [ ] Supersede an existing ADR with another ADR
- [ ] Update the status of an existing ADR
- [ ] Configure custom directory where ADRs are maintained
- [ ] Configure custom statuses for your ADRs
- [ ] Store configuration per ADR repository (directory)
- [ ] Lint ADRs

