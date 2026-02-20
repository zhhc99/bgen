# bgen 主题编写指南

bgen 的主题由两部分构成: **样式** (`style.css`) 和**布局** (`layouts/`).
两者均可独立覆盖, 也可组合使用. 找不到用户文件时, bgen 回退到内置默认值.

> 这篇文档由 AI 生成, 仅供辅助阅读, 项目开发者不对此负责. 制作主题时, 最快的上手方式是直接从 `internal/site/templates` 和 `internal/site/static` 开始修改. 辛苦了~

---

## 文件结构

```
your-blog/
├── static/
│   └── style.css        # 覆盖全局样式
└── layouts/
    ├── base.html        # 页面骨架 (导航 / head / 脚本)
    ├── index.html       # 首页文章列表
    ├── single.html      # 文章详情页
    ├── page.html        # 独立页面 (about 等)
    ├── tags.html        # 所有标签列表
    ├── tag.html         # 单个标签下的文章
    ├── search.html      # 搜索页
    └── 404.html         # 404 页
```

只需放置你想覆盖的文件, 其余继续使用内置模板. 你的文件会**完全替换**内置模板.

---

## 样式覆盖 (`static/style.css`)

内置样式通过 CSS 自定义属性 (变量) 管理视觉 token, 覆盖它们是最快的主题方式.

### 内置变量

```css
:root {
  /* 颜色 */
  --bg:       #faf9f5;   /* 页面背景 */
  --fg:       #141413;   /* 主文字 */
  --text-2:   #5c5a55;   /* 次要文字 */
  --muted:    #b0aea5;   /* 弱化文字 (日期、作者等) */
  --border:   #e8e6dc;   /* 边框、分隔线 */
  --card-bg:  #ffffff;   /* 文章卡片背景 */
  --code-bg:  #f3f1ea;   /* 行内代码、TOC 背景 */
  --accent:   #d97757;   /* 强调色 (hover、引用线、标签) */
  --link:     #6a9bcc;   /* 链接默认色 */

  /* 布局 */
  --max-w:    44rem;     /* 正文最大宽度 */
  --radius:   0.75rem;   /* 圆角 */

  /* 字体 */
  --font:  system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
  --mono:  ui-monospace, "Cascadia Code", "Fira Code", monospace;
}

/* Dark Mode */
html[data-theme="dark"] {
  --bg:      hsl(60, 2.7%, 14.5%);
  --card-bg: hsl(60, 2.5%, 18%);
  --code-bg: hsl(60, 2.0%, 11%);
  --fg:      hsl(60, 8.0%, 88%);
  --text-2:  hsl(60, 3.0%, 62%);
  --muted:   hsl(60, 2.0%, 42%);
  --border:  hsl(60, 2.0%, 22%);
}
```

### 代码块高亮主题

代码块高亮由 highlight.js 提供. 在 `base.html` 中更换 CDN 链接里的主题名即可:

```html
<!-- layouts/base.html -->
<link rel="stylesheet"
  href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/tokyo-night-dark.min.css">
```

可用主题列表见 [highlight.js 主题预览](https://highlightjs.org/demo). 常用选择:

| 主题名 | 风格 |
|---|---|
| `github-dark` | GitHub 深色 (默认) |
| `tokyo-night-dark` | Tokyo Night |
| `catppuccin-mocha` | Catppuccin |
| `nord` | Nord |
| `atom-one-light` | 浅色 |

---

## 布局覆盖 (`layouts/`)

布局文件使用 Go 标准库 `html/template` 语法.

### 模板继承机制

所有页面模板都继承 `base.html`. 固定模式如下:

```html
<!-- layouts/single.html -->
{{template "base.html" .}}
{{define "title"}}{{.Post.Title}} - {{.Site.Config.Title}}{{end}}
{{define "content"}}
  <!-- 你的页面内容 -->
{{end}}
```

`base.html` 提供两个 block:
- `title` — `<title>` 标签内容, 有默认值
- `content` — `<main>` 内的主体内容, **必须定义**

### 可用数据

所有模板均可访问 `.Site`:

```
.Site.Config.Title          → 博客标题 (blog.yaml: title)
.Site.Config.BaseURL        → 站点根 URL (blog.yaml: base_url)
.Site.Config.BasePath       → URL 路径前缀, 通常为 "" 或 "/~john"
.Site.Config.Hero.Header    → 首页 hero 标题
.Site.Config.Hero.Content   → 首页 hero 副文本
.Site.Config.Nav            → map[string]string, 键: "search" / "tags"
.Site.Config.L10n           → map[string]string, 键: "toc"
.Site.Posts                 → []Post, 所有文章 (按时间倒序)
.Site.Tags                  → map[string][]Post
.Site.Pages                 → map[string]Page, 独立页面
```

#### Post 字段

```
.Title          string
.Date           time.Time
.Tags           []string
.Slug           string
.URL            string          → 如 /posts/hello/
.Summary        string
.Author         string
.Cover          string          → 封面相对路径, 如 /posts/hello/cover.jpg; 无封面时为空
.Content        template.HTML   → pandoc 生成的正文 HTML
.TOC            template.HTML   → pandoc 生成的目录 HTML; 无标题时为空
```

#### Page 字段

```
.Title          string
.Slug           string
.URL            string
.Content        template.HTML
```

### 链接 & 路径

所有内部链接必须加 `BasePath` 前缀, 以兼容部署在子路径下的站点:

```html
<a href="{{.Site.Config.BasePath}}{{.Post.URL}}">{{.Post.Title}}</a>
<img src="{{.Site.Config.BasePath}}{{.Post.Cover}}">
```

静态资源 (CSS / JS) 同理:

```html
<link rel="stylesheet" href="{{.Site.Config.BasePath}}/style.css">
```

### 各模板的上下文

| 模板 | 额外可用字段 |
|---|---|
| `index.html` | 仅 `.Site` |
| `single.html` | `.Post` (Post) |
| `page.html` | `.Page` (Page) |
| `tags.html` | 仅 `.Site` |
| `tag.html` | `.Tag` (string), `.Posts` ([]Post) |
| `search.html` | 仅 `.Site` |
| `404.html` | 仅 `.Site` |

### 内置行为: 不要改掉它们

`base.html` 包含两段逻辑, 覆盖时请保留:

**Live reload** (dev server 用, 生产环境自动跳过):
```html
<script>
  if (location.hostname === 'localhost' || location.hostname === '127.0.0.1') {
    var socket = new WebSocket('ws://' + location.host + '/__reload');
    socket.onmessage = function (e) { if (e.data === 'reload') location.reload(); };
  }
</script>
```

**BasePath meta** (search 页的 JS 依赖它定位 `search.json`):
```html
<meta name="base-path" content="{{.Site.Config.BasePath}}">
```

**copy.js** (代码块复制按钮):
```html
<script src="{{.Site.Config.BasePath}}/copy.js" defer></script>
```

---

## 完整示例: 极简双栏布局

以下示例将首页改为左侧固定导航 + 右侧文章流的双栏布局.

**`layouts/base.html`** — 只改结构, 保留必要脚本:

```html
<!DOCTYPE html>
<html lang="zh">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{block "title" .}}{{.Site.Config.Title}}{{end}}</title>
  <meta name="base-path" content="{{.Site.Config.BasePath}}">
  <link rel="stylesheet" href="{{.Site.Config.BasePath}}/style.css">
  <link rel="alternate" type="application/rss+xml"
        title="{{.Site.Config.Title}}"
        href="{{.Site.Config.BasePath}}/feed.xml">
  <link rel="stylesheet"
        href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/github-dark.min.css">
  <script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/highlight.min.js" defer></script>
  <script defer>document.addEventListener('DOMContentLoaded', function(){ hljs.highlightAll(); });</script>
  <script>window.MathJax = { tex: { inlineMath: [['\\(','\\)']], displayMath: [['\\[','\\]']] } };</script>
  <script id="MathJax-script" async src="https://cdn.jsdelivr.net/npm/mathjax@3/es5/tex-chtml.js"></script>
</head>
<body>
  <aside class="sidebar">
    <a href="{{.Site.Config.BasePath}}/" class="site-title">{{.Site.Config.Title}}</a>
    <nav>
      {{if index .Site.Config.Nav "search"}}<a href="{{.Site.Config.BasePath}}/search/">{{index .Site.Config.Nav "search"}}</a>{{end}}
      {{if index .Site.Config.Nav "tags"}}<a href="{{.Site.Config.BasePath}}/tags/">{{index .Site.Config.Nav "tags"}}</a>{{end}}
      {{range $slug, $p := .Site.Pages}}<a href="{{$.Site.Config.BasePath}}{{$p.URL}}">{{$p.Title}}</a>{{end}}
    </nav>
  </aside>
  <main>{{block "content" .}}{{end}}</main>
  <script src="{{.Site.Config.BasePath}}/copy.js" defer></script>
  <script>
    if (location.hostname === 'localhost' || location.hostname === '127.0.0.1') {
      var s = new WebSocket('ws://' + location.host + '/__reload');
      s.onmessage = function(e){ if(e.data==='reload') location.reload(); };
    }
  </script>
</body>
</html>
```

**`static/style.css`** — 追加双栏布局, 其余继承内置:

```css
body {
  display: grid;
  grid-template-columns: 14rem 1fr;
  max-width: 72rem;
  gap: 0 3rem;
  padding: 2rem 2rem 5rem;
}

.sidebar {
  position: sticky;
  top: 2rem;
  height: fit-content;
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.sidebar nav {
  display: flex;
  flex-direction: column;
  gap: 0.6rem;
  font-size: 0.9rem;
}

@media (max-width: 700px) {
  body { grid-template-columns: 1fr; }
  .sidebar { position: static; }
}
```

---

## 注意事项

- 模板使用 `html/template`, 不是 `text/template`. 输出 HTML 内容时用 `template.HTML` 类型的字段 (`.Content`, `.TOC`) 即可, bgen 已处理好转义.
- 覆盖 `base.html` 时, 确保保留 `{{block "content" .}}` 占位, 否则所有页面内容都会消失.
- 样式文件由 bgen 先写入内置 `style.css`, 再用用户的 `static/style.css` **覆盖整个文件** (不是追加合并). 若只想微调少数变量, 在用户 CSS 文件里重新声明 `:root` 变量即可, 其他规则从内置继承.
