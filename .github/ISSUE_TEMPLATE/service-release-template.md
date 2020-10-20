---
name: New Feature/Enhancements Release Checklist
about: Start a checklist for the next production release of new features and enhancements
title: Release checklist for vX.Y.Z
labels: ''
assignees: ''

---

# Description

This issue tracks the requirements for the next release candidate

<Please provide a brief description of the release candidate>

# Jira Epics & Github Milestones
<Please link related JIRA epics and Github Milestones to be released below>

* [Jira ticket](links)
* [Github milestone](links)
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

**PenSec testing**
- [ ] Review and testing by security partner
- [ ] PenSec outcomes documented and/or issues addressed

**Release Architecture Diagram** (handled by architecture)
- [ ] Infrastructure diagrams
- [ ] Network diagrams

**Release documentation**
- TSR
  - [ ] Drafted
  - [ ] Approved
- Implementation Runsheet
  - [ ] Implementation tasks
  - [ ] Verification tasks (TVT, BVT)
  - [ ] Rollback tasks

**Pre-deployment change requirements**

add/remove as necessary

- [ ] Firewall changes
- [ ] Certificate changes
- [ ] ...

**CODEX review** (handled by release lead)
- [ ] completed

**Data governance review** (handled by release lead)
- [ ] Data impact questionnaire
- [ ] Approval

**Fin Crime Assessment** (handled by release lead)
- [ ] Drafted
- [ ] Approved