# 🎯 Internship Tracker

> A modern full-stack web application for tracking and managing internship applications, interviews, lifecycle stages, and AI resume analysis. Built with **Vue 3, TypeScript, Go (Gin), and SQLite/MySQL**.  
> 一款现代化的全栈实习投递看板与求职全流程管理系统。采用 **Vue 3 + TypeScript + Go (Gin) + SQLite/MySQL** 构建。

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

### ✨ Features
- **Interactive Kanban Board**: Visual columns (`Wishlist`, `Applied`, `OA`, `Interview`, `Offer`, `Rejected`) with status transitions.
- **Application Lifecycle Tracking**: Log timelines, status transitions, interview dates, and notes for every job application.
- **AI Resume Analysis**: Match resume text against job requirements and generate structured feedback scores and improvement suggestions.
- **Calendar & Timeline Views**: Calendar overview for application deadlines, interviews, and follow-ups.
- **Admin Management Console**: System overview statistics, user management, and application inspection.
- **Multi-Database Support**: **Zero-configuration SQLite** out-of-the-box (no local database installation required) with optional MySQL support.
- **Unified Single-Port Deployment**: Production build serves both the Vue SPA and Go REST API from a single lightweight container.

### 🏗️ Architecture
```
├── server/             # Go REST API backend (Gin, SQLite/MySQL, JWT Auth)
│   ├── main.go         # Entry point & static file routing
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

### 🚀 Quick Start (Local)

#### 1. Run the Backend
By default, the server uses built-in SQLite with zero configuration:
```bash
cd server
go run .
```
*Starts on `http://127.0.0.1:8080`, automatically creating `./data/internship_tracker.db` and seeding test accounts.*

#### 2. Run the Frontend
```bash
cd web
npm install
npm run dev
```
*Visit `http://127.0.0.1:5173/` in your browser.*

### 🔑 Demo Accounts

| Role | Email | Password | Access |
|------|-------|----------|--------|
| **Demo User** | `demo@example.com` | `Demo@123456` | Kanban Board, Timeline, Calendar, AI Resume Review |
| **Admin** | `admin@example.com` | `Admin@123456` | Full Admin Dashboard & User Management |

---

## 🇨🇳 中文项目说明

### ✨ 核心功能
- **全流程投递看板**：可视化卡片列（`想投/Wishlist`、`已投/Applied`、`笔试/OA`、`面试/Interview`、`录用/Offer`、`归档/Rejected`），支持快速调整状态与详情查看。
- **求职生命周期时间线**：详细记录每一次投递流转、面试通知、备忘录与跟进记录。
- **AI 简历匹配分析**：根据岗位 JD 深度分析简历匹配度，输出契合度评分、竞争优势与改进建议。
- **面试日程日历**：按日期查看即将到来的笔试、面试与投递截止日。
- **轻量级管理后台**：查看平台数据统计（投递总数、用户数、今日新增）、用户状态管理与投递巡查。
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

### 🔑 预置体验账号

| 角色 | 邮箱 | 密码 | 权限与数据 |
|------|------|------|------------|
| **普通演示用户** | `demo@example.com` | `Demo@123456` | 预置名企真实投递、时间线与日历数据 |
| **系统管理员** | `admin@example.com` | `Admin@123456` | 访问管理控制台（`/admin`）、数据概览、用户管理 |

---

## 🐳 Docker 部署 / Run with Docker

```bash
docker build -t internship-tracker .
docker run -p 8080:8080 internship-tracker
```
访问 `http://localhost:8080` 即可使用。

---

## ☁️ 云端一键上线指南 (Render / Railway)

1. 推送代码至 GitHub。
2. 登录 [Render](https://render.com)。
3. 点击 **New +** → **Web Service**，选择你的 GitHub 仓库。
4. 环境选择 **Docker**（Render 将自动识别 Dockerfile 和 render.yaml）。
5. 点击 **Deploy Web Service**，免费获取在线公网演示链接！

---

## 📄 License

MIT License.
