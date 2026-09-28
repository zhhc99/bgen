# TODO: SEO

bgen 从已有数据生成 SEO 信息, 新增配置项均为可选.

---

## 对搜索引擎和爬虫

### sitemap.xml

构建时自动生成, 列出所有文章和页面的 URL 及最后修改时间.

### robots.txt

构建时生成, 默认内容:

```
User-agent: *
Allow: /
Sitemap: <base_url>/sitemap.xml
```

---

## 页面 meta 标签

以下标签加入 `base.html` 的 `<head>`, 通过 block 机制允许各页面覆盖.

### `<html lang>`

目前硬编码为 `zh`. 应从 `blog.yaml` 读取:

```yaml
lang: zh  # 默认 zh
```

### description

```html
<meta name="description" content="...">
```

- 单篇文章: 优先取 front matter 的 `summary`, 否则截取正文前 120 字.
- 首页/标签页/搜索页: 取 `blog.yaml` 中的 `description` 字段 (新增可选项).

### canonical

```html
<link rel="canonical" href="<base_url><page_url>">
```

每个页面声明规范 URL, 帮助搜索引擎识别重复内容.

---

## 分享预览 (Open Graph / Twitter Card)

为社交平台提供分享预览信息:

```html
<meta property="og:title"       content="...">
<meta property="og:description" content="...">
<meta property="og:url"         content="...">
<meta property="og:type"        content="article">  <!-- 首页用 website -->
<meta property="og:image"       content="...">      <!-- 封面图, 无则省略 -->

<meta name="twitter:card"        content="summary_large_image">
<meta name="twitter:title"       content="...">
<meta name="twitter:description" content="...">
<meta name="twitter:image"       content="...">     <!-- 封面图, 无则省略 -->
```

复用页面标题, description, URL 和封面数据.

---

## 暂不做

- JSON-LD / schema.org
- `<meta name="author">`
