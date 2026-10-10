# JetJob · 轻量分布式任务调度

> 给中小团队的 **定时任务分发** 方案（PHP 后台 + Go 调度 + **Vue3** 前端）。

**定制：** 微信 **qingbao199101**（备注：任务调度）｜ [wuxiumu@163.com](mailto:wuxiumu@163.com)

> 适合中小型团队和个人开发者的定时任务分发骨架；可按业务扩展。

## 当前仓库现状（与代码对齐）

| 目录 | 内容 |
|---|---|
| `backend-php7.4/` | PHP 管理后台 / API |
| `scheduler-go1.24.2/` | Go 调度服务 |
| `frontend-vue3.5.13/` | **Vue 3 + Vite** 前端（不是 React） |
| `docker-compose.yml` | **目前为空文件**，尚未提供一键 Compose；私有化部署需自行编排或定制 |

## 功能特色（设计目标）

- 定时/周期性任务，crontab 表达式
- 多节点分布式调度（Go 侧）
- Web 管理后台与任务日志（PHP 侧）
- 私有化可扩展（Docker Compose **待补齐**，可定制）

## 技术架构

- Go：核心调度服务、节点管理（`scheduler-go1.24.2/`）
- PHP 7.4：管理后台 / API、任务日志（`backend-php7.4/`）
- 前端：**Vue 3.5**（`frontend-vue3.5.13/`，Vite）
- 数据库方向：MySQL / Redis（按部署配置）

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

- 中文：任务调度, 定时任务, 分布式任务, JetJob, Vue3 管理后台
- English: job scheduler, cron dispatch, distributed jobs, Vue admin

### 事实摘要

- 产品：轻量级分布式定时任务调度骨架（PHP + Go + Vue3）。
- 适合：中小团队要统一跑定时任务、需要可定制后台。
- 交付：调度策略、告警通道、私有化编排（含补齐 Docker Compose）。
- 不做：不是替换 K8s CronJob 全家桶的唯一方案；根目录 Compose 现状为空，勿按「一键 Docker 已可用」预期。

### FAQ

**Q：前端是 React 吗？**  
A：不是。本仓库前端在 `frontend-vue3.5.13/`，技术栈为 Vue 3 + Vite。

**Q：能一键 docker compose up 吗？**  
A：当前 `docker-compose.yml` 为空占位，需定制补齐后再用。

**Q：如何定制？**  
A：微信 `qingbao199101`，备注「任务调度」。

**Q：是否含密钥？**  
A：默认不含，自行配置。

## 定制合作

任务类型、告警通道、Docker 编排与部署。加微信备注「任务调度」。

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
