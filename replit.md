# Overview

openNAMU is a Python-based wiki engine that provides a complete wiki platform with support for multiple markup languages (NamuMark and Markdown). It's built using Flask as the web framework and supports both SQLite and MySQL databases. The project includes features like user management, ACL (Access Control Lists), discussion boards (BBS), file uploads, version control, and comprehensive admin tools.

# User Preferences

Preferred communication style: Simple, everyday language.

# System Architecture

## Core Framework
- **Web Framework**: Flask with async support (flask[async])
- **Server**: Waitress WSGI server for production deployment
- **Language**: Python 3.8+
- **Deployment Options**: Native Python, Docker containerization

## Database Layer
- **Primary Database**: SQLite with WAL (Write-Ahead Logging) mode enabled
- **Optional Database**: MySQL/MariaDB support via mysqlclient and pymysql
- **Database Configuration**: Controlled via `data/set.json` with db_type selection
- **Connection Pattern**: Context manager pattern (`get_db_connect()`) for automatic connection lifecycle

## Architecture Patterns
- **Modular Route System**: Route handlers organized in separate modules under `route/` directory
- **Hybrid Python-Go Architecture**: Performance-critical operations delegated to Go binaries via subprocess calls (`python_to_golang()` and `python_to_golang_sync()`)
- **Template Rendering**: Jinja2 templates with skin system for customizable UI
- **Async/Await Support**: Asynchronous routes for improved concurrency on I/O operations

## Key Features

### Document Management
- **Markup Support**: NamuMark (primary), Markdown (beta)
- **Version Control**: Full history tracking with revision IDs
- **Document Operations**: Create, edit, move, delete, revert with ACL enforcement
- **Backlink System**: Automatic link tracking and relationship management
- **Edit Requests**: Workflow for proposing changes to restricted documents

### Access Control
- **ACL System**: Granular permissions at document and function level
- **User Roles**: Owner, admin groups, regular users, IP users
- **Permission Groups**: Hierarchical authority system with custom groups
- **Ban System**: Flexible blocking with time limits and various restriction levels

### BBS (Bulletin Board System)
- **Two Modes**: Comment-based and thread-based discussions
- **Post Management**: Create, edit, delete, pin, hide posts
- **Nested Comments**: Thread-style comment replies with depth tracking
- **ACL Integration**: Per-board access control

### File Management
- **Upload System**: Image and file uploads with extension filtering
- **Storage**: File-based storage with SHA-224 hashing for names
- **Image Handling**: Pillow library for image processing
- **Size Limits**: Configurable upload size restrictions

### User System
- **Authentication**: Password-based with bcrypt/SHA-256 hashing options
- **2FA Support**: Two-factor authentication capability
- **User Levels**: Experience-based level system
- **Watchlists**: Document monitoring and notifications

## Data Flow
1. HTTP Request → Flask route handler
2. Authentication/ACL check via `acl_check()` and `ip_check()`
3. Database operations via context managers
4. Business logic processing (Python or delegated to Go)
5. Template rendering with `flask.render_template()` and `easy_minify()`
6. Response with HTML or JSON

## Rendering Pipeline
- **Parser**: `render_set()` function handles markup parsing
- **Caching**: Rendered content can be cached for performance
- **Timeout Protection**: `edit_timeout()` prevents long-running renders
- **Include System**: Transclusion support for embedding documents

## Security Features
- **Captcha Integration**: reCAPTCHA support for spam prevention
- **XSS Protection**: Input filtering via `html_filter` table
- **Edit Filters**: Regex-based content filtering
- **CSRF Protection**: Form validation with tokens
- **Rate Limiting**: Slow edit checks and IP-based throttling

## Internationalization
- **Language Files**: JSON-based translations in `lang/` directory
- **Dynamic Language Loading**: `get_lang()` function for runtime translation
- **Supported Languages**: English (en-US), Korean (ko-KR)

# External Dependencies

## Core Python Libraries
- **flask**: Web framework and routing
- **waitress**: Production WSGI server
- **aiohttp**: Async HTTP client for external requests
- **asyncio**: Asynchronous I/O support
- **requests**: Synchronous HTTP client

## Data Processing
- **orjson**: Fast JSON serialization/deserialization (optional but recommended)
- **diff-match-patch**: Text diffing for version comparison
- **pillow**: Image processing and manipulation

## Database Drivers
- **pymysql**: Pure Python MySQL client
- **mysqlclient**: C-based MySQL connector (optional, higher performance)

## Frontend Assets (External Projects)
- **highlight.js**: Syntax highlighting for code blocks
- **KaTeX**: Mathematical formula rendering
- **Feather Icons**: Icon set for UI elements
- **Go modules**: Performance-critical operations (see route_go/go.mod)

## System Integration
- **Go Binaries**: Platform-specific binaries (AMD64/ARM64) for Linux/Windows
  - Downloaded from: `https://github.com/openNAMU/GopenNAMU/releases/`
  - Located in: `route_go/bin/`
- **SMTP Support**: Email functionality for notifications and password resets

## Development Tools
- **Emergency Tool**: `emergency_tool.py` for administrative recovery operations
- **Language Helper**: `lang/help_tool.py` for translation management