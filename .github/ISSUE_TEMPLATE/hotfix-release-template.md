---
name: Hotfix Release Checklist
about: Start a checklist for the next hot fix release (e.g. security patches, code/config hot fixes)
title: Hotfix Release checklist for vX.Y.Z
labels: ''
assignees: ''

---

# Description

This issue tracks the requirements for the next hot fix release candidate

<Please provide a brief description of the hot fix release candidate>

# Github Issues
<Please link related Github Issues to be released below>

* [Github Issues](link)
* ....


# Pre-release checklist

These steps must be completed before we consider deploying code to production (ie. move tickets to 'ready for release' state)

**System Documentation**
- [ ] Sequence diagrams
- [ ] Data mapping sheet (for data governance)
- [ ] Data flows
- [ ] API specs

**Testing**
- [ ] Unit tests
- [ ] Code coverage ≥90% (via sonarqube)
- [ ] System tests
- [ ] System integration tests

**Code and security scans**
- [ ] Checkmarx: High/medium issues addressed and/or seek guidance from security
- [ ] Blackduck: issues addressed and/or seek guidance from security (focus on security vulnerability issues)

**Release candidate ready**
- [ ] Deployed to Pre-Prod
- [ ] Tested headless or against client app
- [ ] Informal PenSec Testing (new features only)


# Ready for release

These steps must be completed during release readiness work

**Release documentation**
- TSR
  - [ ] Drafted
  - [ ] Approved
- Implementation Runsheet
  - [ ] Implementation tasks
  - [ ] Verification tasks (TVT, BVT)
  - [ ] Rollback tasks