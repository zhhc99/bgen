# bgen 主题编写指南

bgen 的主题包含**样式** (`style.css`) 和**布局** (`layouts/`), 均可独立替换. 缺少用户文件时使用内置版本. `custom.css` 用于增量调整样式.

> 本文由 AI 生成, 仅供参考, 项目开发者不对内容负责. 内置布局和样式见 `internal/site/templates` 和 `internal/site/static`.

---

## 文件结构

```
your-blog/
├── static/
│   ├── style.css        # (可选) 完全替换内置样式
│   └── custom.css       # (可选) 增量补充和覆盖当前主题样式
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

按需创建以上文件. `static/style.css` 和 `layouts/` 下的同名文件会完全替换内置文件; `static/custom.css` 只需包含样式改动.

---

## 样式覆盖 (`static/custom.css` / `static/style.css`)

日常调整可直接写入 `static/custom.css`. 例如:

```css
/* static/custom.css */
:root {
  --accent: #387b67;
  --max-w: 50rem;
}
```

默认布局仅在 `static/custom.css` 存在时引用它, 加载顺序为 `style.css`, highlight.js 高亮样式, `custom.css`. 两个本地 CSS 文件分别输出到站点根目录, 引用自动带上 `BasePath` 前缀.

覆盖遵循 CSS 层叠规则: 同优先级时, 后加载的声明生效; 选择器权重和 `!important` 仍然有效. 修改深色模式变量时, 可沿用内置的 `html[data-theme="dark"]` 选择器.

完整替换样式时, 使用 `static/style.css`, 也可继续用 `static/custom.css` 调整. 自定义 `layouts/base.html` 时, 需保留 `custom.css` 的条件引用, 见下方示例.

内置 CSS 变量和深色模式定义见 [style.css](../internal/site/static/style.css).

代码块高亮主题在 `layouts/base.html` 的 highlight.js 样式链接中选择. 可用主题见 [highlight.js 主题预览](https://highlightjs.org/demo).

---

## 布局覆盖 (`layouts/`)

布局文件使用 Go 标准库 `html/template` 语法.

### 模板继承机制

页面模板通过 `base.html` 组织布局, 例如:

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
.Site.Config.Nav            → []NavItem, 按配置顺序排列, 字段: Title / URL
.Site.Config.NavTitle "/search/" → 对应导航名称, 未配置时为空
.Site.Config.NavURL "/about/"    → 添加部署路径前缀, 外部链接保持原样
.Site.Config.Text "toc"     → 界面文案, 优先使用 text-override, 未配置时使用默认值
.Site.Posts                 → []Post, 所有文章 (按时间倒序)
.Site.Tags                  → map[string][]Post
.Site.Pages                 → map[string]Page, 独立页面
.Site.HasCustomCSS          → bool, 项目是否提供 static/custom.css
```

文案键: `toc`, `search-placeholder`, `not-found`, `go-home`, `copy`. 在模板中通过 `{{.Site.Config.Text "键名"}}` 读取.

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
.CoverCaption   string          → front matter 的 cover_caption, 普通文本, 可为空
.Content        template.HTML   → pandoc 生成的正文 HTML
.TOC            template.HTML   → pandoc 生成的目录 HTML; 无标题时为空
```

默认文章模板在封面和图注都存在时生成 `.cover-figure`, 内含 `.featured-image` 和 `figcaption`, 位于目录之前. 自定义模板通过 `.Post.CoverCaption` 输出图注, 文本由 `html/template` 转义.

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
<!-- 如果加载了代码高亮等样式, 将它们放在 custom.css 之前 -->
{{if .Site.HasCustomCSS}}<link rel="stylesheet" href="{{.Site.Config.BasePath}}/custom.css">{{end}}
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

### 布局中的功能片段

自定义 `base.html` 时, 保留以下片段以支持对应功能:

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

## 模板中的 HTML 内容

`html/template` 自动转义普通文本. `.Content` 和 `.TOC` 的类型为 `template.HTML`, 可直接输出转换后的 HTML. `base.html` 通过 `{{block "content" .}}{{end}}` 插入页面主体, 自定义布局时需保留此占位.
