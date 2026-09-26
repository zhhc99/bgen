# bgen 设计文档

bgen 是一个用 Go 写的极简静态博客生成器. 目标用户是追求简洁, 有一定技术能力, 想要完全理解自己工具的写作者.

## 用户维护的文件

见 `README.md` "📝 配置项和约定".

## 编码格式原则

- 代码自述, 尽量避免注释.
- 当代码无法自述时, 用中文写简洁注释.
- 任何时候都不允许使用中文标点. 只能使用英文标点.
- 在规范的前提下, 尽量用最简短的代码实现功能.

## bgen 做的事

1. 读取 `content/` 下所有 markdown 文件
2. 每篇文章用 Pandoc 处理: markdown -> HTML, 处理 TeX, 图注, 代码块
3. Pandoc 同时生成 TOC
4. 将 HTML 内容注入 Go html/template 模板
5. 在临时目录生成静态文件, 成功后替换 output/, 清除旧文件
6. 独立页面照常生成, 导航按 nav 列表展示
7. nav 包含 /search/ 时生成搜索页和 search.json (标题 + URL + 日期)
8. dev 模式: 本地 HTTP server + 文件监听自动重建

## 生成的页面

| 页面     | URL                  | 说明                   |
| :------- | :------------------- | :--------------------- |
| 文章列表 | `/`                  | 所有文章按时间倒序     |
| 文章页   | `/posts/slug/`       | 正文 + TOC + tags      |
| tag 页   | `/tags/math/`        | 该 tag 下的文章        |
| 搜索页   | `/search/`           | 前端 Fuse.js, 只搜标题 |
| 特殊页面 | `/about/`, `/links/` | 纯内容, 无列表逻辑     |
| 404      | `/404.html`          | 错误反馈页             |

导航由 nav 中的 title 和 url 按顺序组成, 支持站内和外部链接. nav 包含 /tags/ 时生成标签索引和标签列表页.

文章和独立页面均支持 ignore: true, 跳过生成. 输出目录由 bgen 管理, 构建失败保留上次结果. 输出路径不能覆盖项目或源文件.

## 技术栈

- 语言: Go
- Markdown 处理: 调用本地 Pandoc (`exec.Command`)
- 模板: 标准库 `html/template`
- YAML 解析: `gopkg.in/yaml.v3`
- 文件监听: `github.com/fsnotify/fsnotify`
- Dev server: `net/http`, `github.com/coder/websocket`
- 前端搜索: Fuse.js (CDN), 消费 search.json
- 依赖极少, 编译为单一二进制

## 命令行

```
bgen init         # 初始化目录结构
bgen build        # 构建到 output/
bgen serve        # dev server, watch + reload
bgen version
bgen help
```

## 主题/自定义

默认模板和样式内置在二进制里 (embed). 用户在项目根目录放同名文件即可覆盖:

- `layouts/single.html` 覆盖文章页模板
- `static/style.css` 覆盖样式

默认样式应是真正好看的, 干净的, 不是 placeholder. 用户不覆盖也能直接用.

## 计划做的事

- pandoc parse 应该用 errgroup
- `extractImageRefs` 也支持内嵌 html, 如 `<img src="...">`
- 用 bgen 生成的 example & usage site

## 不做的事

- 文章的 l10n (责任不在写作者)
- 复杂 taxonomy (只有 tags, 没有 categories)
- 任何需要文档才能理解的配置项
