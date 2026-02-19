# bgen

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/zhhc99/bgen)](https://github.com/zhhc99/bgen)
[![Latest Release](https://img.shields.io/github/v/release/zhhc99/bgen)](https://github.com/zhhc99/bgen/releases)

**静态博客站点生成器.**

核心想法: 把时间花在用 Markdown 书写博文, 而不是做各种配置上.

## 🛠 功能

- [x] Pandoc 的一切: Tex, 图注
- [x] 目录, 标签, 搜索
- [x] Dev server + 文件监听自动重建
- [x] 自定义主题

> 搜索功能基于 Fuse.js.

## 📦 快速安装

**依赖:** [Pandoc](https://pandoc.org/installing.html) (用于 Markdown 解析)

建议使用 `go install`:

```bash
# 可能需要手动将 '$(go env GOPATH)/bin' 添加到 PATH
go install github.com/zhhc99/bgen@latest
```

## 📖 快速开始

**确保博客项目有如下结构:**

```
content/
├── posts/
│   ├── 2024-01-hello.md       # 简单文章
│   ├── 2024-01-hello.jpg      # 文章封面图
│   └── complex-post/          # 复杂文章用 bundle
│       ├── index.md
│       ├── cover.png
│       └── figure1.png
└── about.md                   # 特殊页面
layouts/                       # (可选) 覆盖模板
static/                        # (可选) 静态文件, 原样打包
blog.yaml                      # 站点配置
```

然后就可以愉快地写文了.

**本地预览:**

```bash
bgen serve    # 启动 dev server, 监听文件变化
```

**构建生产版本:**

```bash
bgen build    # 输出到 output/
```

## 📝 配置项和约定

**blog.yaml 示例:**

```yaml
title: My Blog
base_url: https://example.com # 站点根域名
hero:                         # 首页顶部展示区
  header: John
  content: This is my blog!
nav:                          # 导航栏项目的名称. 填 "" 删去对应项
  search: search
  tags: tags
l10n:
  toc: Table of Contents
front-matter-defaults:        # markdown 元数据的默认值
  author: John
```

**Markdown 前置元数据:**

```
---
title: Hello World
date: 2024-01-01
tags: [tech, life]
slug: slug-to-this-post           # 默认为文件名
summary: this post has nothing... # 默认从文章截取
author: Alice                     # 若不填写, 由 blog.yaml 覆盖
---

这是文章正文, 使用 Pandoc's Markdown. 支持 TeX: $E = mc^2$

![这是图注](./image.png)
```

## 🤔 常见问题

**Q: 如何自定义主题?**

A: 主题由 style (外观) 和 layout (布局) 构成.

- style: 位于 `static/style.css`, 找不到时回退到内置模板.
- layout: 位于 `layouts/`, 找不到时回退到内置模板.

**Q: layout 有哪些?**

A: 见仓库 `internal/site/templates`, 默认内容非常简单. 例如, `layouts/single.html` 覆盖文章页模板. TeX 解析默认使用 MathJax, 位于 `internal/site/templates/base.html`.

**Q: 如何添加导航页面?**

A: 在 `content/` 下直接创建的 markdown 会被 bgen 理解成可导航的单独页面.

**Q: 封面图怎么添加?**

A: 文章同目录下放同名图片文件 (如 `post.md` 对应 `post.jpg`), 或在 bundle 中使用 `cover.jpg`. 支持各种常见图片格式 (但不包括 svg).

## 🔨 编译源代码

```bash
git clone https://github.com/zhhc99/bgen.git
cd bgen
go build -o bgen .
```

带版本号编译:

```bash
go build -ldflags "-X 'main.Version=v1.0.0'"
```

## 🚀 发布

推送 tag 触发 goreleaser:

```bash
git tag vX.Y.Z
git push origin vX.Y.Z
```

## 🎯 想法

**写博客应该是在写作, 不是在配置工具.**

现有工具 (Hugo, Jekyll 等) 配置项太重, 基本功能 (如 TeX) 做不好, 每个主题还有细微区别. 这给我带来了困扰.

bgen 是我对此交出的答卷.

我弄了一个工具 [hugo2bgen](https://github.com/zhhc99/hugo2bgen/) 用来从 hugo 快速迁移. 工具可以处理 frontmatter 的字段并迁移封面图, 不修改原始文件.
