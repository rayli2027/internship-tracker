# 🎯 Internship Tracker

> A modern full-stack web application for tracking and managing internship applications, interviews, lifecycle stages, and AI resume analysis. Built with **Vue 3, TypeScript, Go (Gin), and SQLite/MySQL**.  
> 一款现代化的全栈实习投递看板与求职全流程管理系统。采用 **Vue 3 + TypeScript + Go (Gin) + SQLite/MySQL** 构建。

[![Live Demo](https://img.shields.io/badge/Live-Demo%20Online-brightgreen?style=for-the-badge&logo=render)](https://internship-tracker-zfj8.onrender.com)
[![Go](https://img.shields.io/badge/Backend-Go%201.22+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Vue](https://img.shields.io/badge/Frontend-Vue%203%20%2B%20TypeScript-4FC08D?style=flat&logo=vuedotjs)](https://vuejs.org)
[![Database](https://img.shields.io/badge/Database-SQLite%20%2F%20MySQL-003B57?style=flat&logo=sqlite)](https://sqlite.org)
[![Docker](https://img.shields.io/badge/Deploy-Docker%20%2F%20Render-2496ED?style=flat&logo=docker)](https://docker.com)

---

## 🌐 Language Navigation / 语言导航
- [English Documentation](#-english-documentation)
- [中文项目说明](#-中文项目说明)

---

## 🇺🇸 English Documentation

### 🔗 Live Demo
Experience the deployed web application live:  
👉 **[https://internship-tracker-zfj8.onrender.com](https://internship-tracker-zfj8.onrender.com)**

> *Note: Hosted on Render's free tier. If inactive, the server may take ~30–45 seconds to spin up on first request.*

### ✨ Key Features
- **Interactive Kanban Board**: Visual columns (`Wishlist`, `Applied`, `OA`, `Interview`, `Offer`, `Rejected`) with dynamic status transitions.
- **Application Lifecycle Tracking**: Log timelines, status transitions, interview dates, and notes for every job application.
- **AI Resume Compatibility Analysis**: Match resume text against job requirements and generate structured feedback scores and improvement suggestions.
- **Calendar & Timeline Views**: Calendar overview for application deadlines, interviews, and follow-ups.
- **Role-Based Access Control**: Multi-tenant user isolation with separate administrative overview and user auditing.
- **Multi-Database Support**: **Zero-configuration pure-Go SQLite** out-of-the-box (no local database installation required) with optional MySQL support.
- **Unified Single-Port Deployment**: Production build serves both the Vue SPA and Go REST API from a single lightweight container.

### 🏗️ Architecture
```
├── server/             # Go REST API backend (Gin, SQLite/MySQL, JWT Auth)
│   ├── main.go         # Entry point & static SPA file routing
│   ├── config.go       # Environment & database configuration
│   ├── database.go     # Schema migrations & demo data seeding
│   ├── handlers.go     # REST API route handlers
│   ├── lifecycle.go    # Application lifecycle events
│   └── resume.go       # Resume parsing & AI analysis
├── web/                # Vue 3 Frontend (Vite, TypeScript, Tailwind/CSS)
│   ├── src/views/      # Board, Calendar, Timeline, Resume, Admin views
│   └── src/components/ # Modular UI components & drawers
├── Dockerfile          # Multi-stage production container build
└── render.yaml         # One-click Render deployment configuration
```

### 🚀 Quick Start (Local Development)

#### 1. Run the Backend
By default, the server uses built-in SQLite with zero configuration:
```bash
cd server
go run .
```
*Starts on `http://127.0.0.1:8080`, automatically creating `./data/internship_tracker.db` and seeding sample data.*

#### 2. Run the Frontend
```bash
cd web
npm install
npm run dev
```
*Visit `http://127.0.0.1:5173/` in your browser.*

### 🔑 Demo Account

Visitors can log in immediately using the pre-seeded demo user account or register their own:

| Account | Email | Password | Access |
|---------|-------|----------|--------|
| **Demo User** | `demo@example.com` | `Demo@123456` | Full Kanban Board, Timeline, Calendar, AI Resume Review |

---

## 🇨🇳 中文项目说明

### 🔗 在线体验
在线访问已部署的完整全栈服务：  
👉 **[https://internship-tracker-zfj8.onrender.com](https://internship-tracker-zfj8.onrender.com)**

> *提示：项目部署于 Render 免费实例，若处于休眠状态首次打开需等待约 30~45 秒冷启动。*

### ✨ 核心功能
- **全流程投递看板**：可视化卡片列（`想投/Wishlist`、`已投/Applied`、`笔试/OA`、`面试/Interview`、`录用/Offer`、`归档/Rejected`），支持快速调整状态与详情抽屉。
- **求职生命周期时间线**：详细记录每一次投递流转、面试通知、备忘录与跟进记录。
- **AI 简历匹配分析**：根据岗位 JD 深度分析简历匹配度，输出契合度评分、竞争优势与改进建议。
- **面试日程日历**：按日期查看即将到来的笔试、面试与投递截止日。
- **用户隔离与权限控制**：基于 JWT 认证与数据隔离，支持普通用户与管理巡查权限。
- **零配置开箱即用数据库**：内置 **纯 Go SQLite 驱动**，本地无需安装任何 MySQL 或 Docker，启动即自动建表并灌入演示数据；同时兼容 MySQL。
- **一体化轻量部署**：支持单端口同构部署，后端直接托管前端编译产物。

### 🚀 快速启动

#### 1. 启动后端
```bash
cd server
go run .
```
*后端将在 `http://127.0.0.1:8080` 启动，自动生成数据库文件并写入真实演示岗位数据（Google、Stripe、字节跳动、Jane Street 等）。*

#### 2. 启动前端
```bash
cd web
npm install
npm run dev
```
*在浏览器打开 `http://127.0.0.1:5173/` 即可体验。*

### 🔑 演示账号

访客可直接使用预置演示账号登录体验，亦支持自由注册新账号：

| 账号 | 邮箱 | 密码 | 权限与数据 |
|------|------|------|------------|
| **演示用户** | `demo@example.com` | `Demo@123456` | 预置名企真实投递、时间线与日历数据 |

---

## 🐳 Docker 部署 / Run with Docker

```bash
docker build -t internship-tracker .
docker run -p 8080:8080 internship-tracker
```
访问 `http://localhost:8080` 即可使用。

---

## 📄 License

MIT License.
