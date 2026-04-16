# TypeScript 三大意义

## 一、核心价值

前端开发者用 TS 的三大意义：

1. **发往后端的参数**：必须合规，绝不乱传
2. **后端返回的数据**：我明确知道结构，访问放心
3. **类型推断**：少写代码，自动提示，开发飞快

---

## 二、发往后端的参数 → 请求契约

> 这就是 TS 的 **「请求契约」**

定义好 `CourseQueryParams` / `CourseForm`，搜索、分页、提交表单时，少传字段、传错类型、传错名字，直接报错

**彻底杜绝：** page 传成字符串、keyword 漏传、接口 400/500 这种低级错误

**对应你的代码：**

```ts
// 定义好了，发请求时必须严格遵守
interface CourseQueryParams {
  keyword: string
  page: number
  pageSize: number
}

// 写错直接爆红，根本发不出错误请求
setSearchParams({ keyword: '', page: '1' }) // ❌ 报错！page 必须是数字
```

---

## 三、后端返回的数据 → 响应安全

> 这就是 TS 的 **「响应安全」**

- 不用猜后端返回什么字段
- 不用写 `res?.data?.list?.[0]?.name` 怕报错
- 表格渲染、取值、循环，想点什么就点什么，100% 放心

**对应你的代码：**

```ts
interface CourseItem {
  id: string
  name: string
  teacher: string
}

// courseList 我明确知道是 CourseItem[]，渲染时随便访问
courseList.map(item => item.name) // ✅ 放心访问
```

---

## 四、类型推断 → 效率加成

> 这是 TS 的 **「效率加成」**

- 不用手动写类型，TS 自动猜出来
- VSCode 自动补全字段、提示错误
- 不用记接口字段、不用查文档

> 这就是你说的：再多一点就是类型推断！
