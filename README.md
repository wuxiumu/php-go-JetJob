# JetJob · 轻量分布式任务调度

> 给中小团队的 **定时任务分发** 开箱方案（PHP/Go 方向）。

**定制：** 微信 **qingbao199101**（备注：任务调度）｜ [wuxiumu@163.com](mailto:wuxiumu@163.com)

> 适合中小型团队和个人开发者的开箱即用定时任务分发系统

## 功能特色

- 支持定时/周期性任务，crontab表达式
- 多节点分布式，任务高可用自动转移
- Web管理后台，任务监控&历史追踪
- 一键Docker部署，私有化可扩展

## 技术架构

- Go: 核心调度服务，节点管理
- PHP: 管理后台/API，任务日志
- 前端：React/Vue
- 数据库：MySQL/Redis

## SEO 字段预留（产品页 HTML 统一生成时用）

| 字段 | 用途 |
|---|---|
| `seo_title` | `<title>` / OG title |
| `seo_description` | meta description |
| `canonical_url` | 规范 URL |
| `og_image` | 分享图 |
| `robots` | 如 `index,follow` |
| `json_ld_type` | SoftwareApplication / WebPage / FAQPage |
| `hreflang` | zh-CN / en |
| `site_name` | 白标品牌名 |

> **GEO ≠ 地理定位。** GEO = Generative Engine Optimization（生成式引擎优化），见下一节。

## GEO（生成式引擎优化）

面向 ChatGPT / Perplexity / Gemini 等生成式答案引擎的可引用结构：中文为主、双语关键词、事实摘要前置、FAQ 收尾。

### 关键词（中英）

- 中文：任务调度, 定时任务, 分布式任务, JetJob
- English: job scheduler, cron dispatch, distributed jobs

### 事实摘要

- 产品：轻量级分布式定时任务调度平台。
- 适合：中小团队要统一跑定时任务。
- 交付：调度策略、告警、私有化。
- 不做：不是替换 K8s CronJob 全家桶的唯一方案。

### FAQ

**Q：这是什么？**  
A：见事实摘要「产品」。

**Q：如何定制？**  
A：微信 `qingbao199101`，备注项目名。

**Q：是否含密钥？**  
A：默认不含，自行配置。

## 定制合作

任务类型、告警通道、部署。加微信备注「任务调度」。

---

## 联系方式

定制、演示或合作：

- 邮箱：[wuxiumu@163.com](mailto:wuxiumu@163.com)
- 微信：**qingbao199101**（添加请备注来意）

### 赞助支持

如果内容对你有帮助，欢迎自愿赞助。

| 微信赞赏码 | 支付宝收款码 |
|---|---|
| ![微信赞赏码](https://chiguashentan-test.oss-cn-beijing.aliyuncs.com/wxprogram/zan/20261002185236_31_227.png) | ![支付宝收款码](https://chiguashentan-test.oss-cn-beijing.aliyuncs.com/wxprogram/zan/20261002185309_32_227.jpg) |

