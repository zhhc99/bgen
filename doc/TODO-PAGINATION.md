# TODO: Pagination

首页和标签页在文章数量增多后需要分页.

---

## 配置

`blog.yaml` 新增可选项:

```yaml
pagination:
  page_size: 10  # 每页文章数, 默认 10
```

---

## URL 结构

静态站点分页使用目录形式, 不用查询参数:

```
/           -> 第 1 页
/page/2/    -> 第 2 页
/page/3/    -> 第 3 页

/tags/tech/         -> 标签页第 1 页
/tags/tech/page/2/  -> 标签页第 2 页
```

---

## 模板变量

分页页面向模板传入 `Pager` 对象:

```
Pager.Current    当前页码 (从 1 开始)
Pager.Total      总页数
Pager.HasPrev    是否有上一页
Pager.HasNext    是否有下一页
Pager.PrevURL    上一页 URL
Pager.NextURL    下一页 URL
```

---

## SEO

分页页面在 `<head>` 中加入:

```html
<link rel="prev" href="...">  <!-- 非第一页 -->
<link rel="next" href="...">  <!-- 非最后一页 -->
```

告知搜索引擎页面间的关系.
