export function showMessage(message: string): void {
  console.error(message)
  // 实际项目中这里会使用 naive-ui 的 message 组件
  // window.$message?.error(message)
}
