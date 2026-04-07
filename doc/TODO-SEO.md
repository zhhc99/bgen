# TODO: SEO

bgen 的设计哲学是约定优于配置, 因此 SEO 相关信息应尽量从已有数据自动推断,
不引入新的必填配置项.

---

## 对搜索引擎和爬虫

### sitemap.xml

构建时自动生成, 列出所有文章和页面的 URL 及最后修改时间.
无需用户配置.

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

防止重复内容问题, 每个页面都应有.

---

## 分享预览 (Open Graph / Twitter Card)

分享到社交平台时决定展示效果, 统一处理.

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

数据来源与 description 一致, 不引入额外配置.

---

## 暂不做

- JSON-LD / schema.org: 收益有限, 实现较繁琐, 与简洁哲学不符.
- `<meta name="author">`: 意义不大, 暂跳过.
