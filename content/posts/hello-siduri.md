---
title: Hello, Siduri
slug: hello-siduri
date: 2026-08-25
summary: Why I am building a small, honest site about shipping software with agents.
plain_summary: I am documenting what it takes to ship useful software with an agent while keeping a person responsible for the result.
tags:
  - method
  - build-log
draft: true
---

I am building Siduri to document the work behind agent-assisted software, not to pretend the agent is the work.

The name comes from the tavern keeper who gives Gilgamesh useful advice and then points him toward the next step. That is the model here. I want the tools to move quickly, and I want a person to stay at the gate where a decision becomes public.

## What this site is for

I will write about things I actually build: the prompts that helped, the approaches that failed, and the small numbers that tell me whether a shortcut was worth taking. A build log should include the dead end. A tool release should include the maintenance cost. A postmortem should say what changed afterward.

The first version of this site is itself an example. It is a Go program that turns Markdown into static HTML. It has no database, no client-side framework, and no analytics script. The initial build has one content file, five allowed tags, and a 32768-byte guard on the agent contract. Those are not impressive numbers. They are concrete numbers, which is more useful.

## The line I am drawing

An agent can read a requirement, write a patch, run a check, and explain a failure. It cannot decide that an unpublished post should become public. That last action stays human. The repository should make this boundary visible instead of hiding it behind an automation button.

I am starting with an introduction because the project needs a thesis before it needs a backlog. Later posts will test the thesis against real tools and real failures. If the site says a system is deterministic, the build should prove it. If it says a draft is not published, the draft should be absent from the output. If either claim stops being true, the check should turn red.

That is the work. I will record the next result when I have one.
