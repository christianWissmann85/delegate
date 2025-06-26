# Delegate Documentation

Welcome to the Delegate documentation! Delegate is an MCP server that allows Claude Code to delegate heavy tasks to other LLMs, saving precious context tokens.

## 📚 Documentation Structure

### 🏗️ Architecture
- [**Architecture Specification**](architecture/architecture-spec.md) - Complete technical specification
- [**Module Architecture**](architecture/module-architecture.md) - Module boundaries and interfaces
- [**Day 0 Decisions**](architecture/day-0-decisions.md) - Key implementation decisions

### 🔧 Development
- [**Implementation Roadmap**](development/implementation-roadmap.md) - 3-week development plan
- [**Project Charter**](development/PROJECT_CHARTER.md) - Mission and core principles
- [**NO SCOPE CREEP**](development/NO_SCOPE_CREEP.md) - The sacred document
- [**Testing Strategy**](development/delegate-testing.md) - Test plans and approaches
- [**Implementation Questions**](development/implementation-questions.md) - Technical decisions

### 📖 Guides
- [**Getting Started**](guides/getting-started-guide.md) - Quick start guide
- [**Claude Code Guide**](guides/claude-code-guide.md) - How to use Delegate with Claude Code
- [**Use Cases**](guides/claude-code-use-cases.md) - Creative ways to use Delegate
- [**Deployment Guide**](guides/deployment-guide.md) - Installation and configuration

### 📋 Reference
- [**API Reference**](reference/api-reference.md) - Complete tool specifications
- [**Model Reference**](reference/model-reference-card.md) - Which model to use when
- [**MCP Tool Descriptions**](reference/mcp-tool-descriptions.md) - Tool registration details

## 🚀 Quick Links

- **For Users**: Start with the [Getting Started Guide](guides/getting-started-guide.md)
- **For Claude Code**: Read the [Claude Code Guide](guides/claude-code-guide.md)
- **For Developers**: Check the [Implementation Roadmap](development/implementation-roadmap.md)
- **For Contributors**: Read [NO SCOPE CREEP](development/NO_SCOPE_CREEP.md) first!

## 🎯 Core Philosophy

Delegate does three things:
1. **invoke** (STEP 1) - Delegate tasks to LLMs (files <1MB each, total can exceed)
2. **check** (STEP 2) - Get metadata about outputs before retrieving
3. **read** (STEP 3) - Retrieve results or write directly to disk (saves tokens!)

That's it. No more. Read [NO SCOPE CREEP](development/NO_SCOPE_CREEP.md) to understand why.