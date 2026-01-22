---
inclusion: manual
---

# 表格列宽度标准

## 列宽度规范

### 长文本列（域名、URL、任务名称等）
- `size: 350` - 默认宽度
- `minSize: 250` - 最小宽度，保证可读性
- 无 `maxSize` - 允许用户自由拖拽扩展
- 样式：`break-all leading-relaxed whitespace-normal`

### 中等文本列（引擎名称、状态等）
- `size: 120-150`
- `minSize: 80-100`
- 通常使用 Badge 显示

### 短文本/数字列（日期、计数、进度等）
- `size: 100-150`
- `minSize: 80-120`

### 固定列（选择框、操作按钮）
- `size/minSize/maxSize` 相同值
- `enableResizing: false`
- 选择框：40px
- 操作按钮：60-120px

## 单元格样式

### 多行文本（推荐）
```tsx
<div className="flex-1 min-w-0">
  <span className="text-sm font-medium break-all leading-relaxed whitespace-normal">
    {text}
  </span>
</div>
```

### 可点击链接
```tsx
<button className="text-sm font-medium hover:text-primary hover:underline underline-offset-2 transition-colors cursor-pointer text-left break-all leading-relaxed whitespace-normal">
  {text}
</button>
```

### Badge 列表（横向自动换行）
```tsx
<div className="flex flex-wrap items-center gap-1.5">
  {badges}
</div>
```

## 列定义示例

```tsx
// 长文本列（域名/URL/名称）
{
  accessorKey: "targetName",
  size: 350,
  minSize: 250,
  header: ({ column }) => <DataTableColumnHeader column={column} title="Target" />,
  cell: ({ row }) => {
    const value = row.getValue("targetName") as string
    return (
      <div className="flex-1 min-w-0">
        <span className="text-sm font-medium break-all leading-relaxed whitespace-normal">
          {value}
        </span>
      </div>
    )
  },
}

// 中等文本列（Badge）
{
  accessorKey: "engineName",
  size: 120,
  minSize: 80,
  header: ({ column }) => <DataTableColumnHeader column={column} title="Engine" />,
  cell: ({ row }) => (
    <Badge variant="secondary">{row.getValue("engineName")}</Badge>
  ),
}

// 日期列
{
  accessorKey: "createdAt",
  size: 150,
  minSize: 120,
  header: ({ column }) => <DataTableColumnHeader column={column} title="Created At" />,
  cell: ({ row }) => (
    <span className="text-sm text-muted-foreground">
      {formatDate(row.getValue("createdAt"))}
    </span>
  ),
}

// 固定操作列
{
  id: "actions",
  size: 80,
  minSize: 80,
  maxSize: 80,
  enableResizing: false,
  cell: ({ row }) => <ActionButtons row={row} />,
}
```

## 注意事项

1. 不要使用折叠/省略号 + Popover，改用多行显示
2. 长文本列不设置 `maxSize`，让用户自由调整
3. 使用 `flex-wrap` 让 Badge 列表自动换行
4. 保持列宽一致性，相同类型的列使用相同的宽度配置
