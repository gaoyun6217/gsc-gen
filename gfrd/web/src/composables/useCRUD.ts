import { ref, type Ref } from 'vue'
import { showMessage } from './status'

interface CRUDOptions<T, P = any> {
  api: {
    list: (params: P) => Promise<{ list: T[]; total: number }>
    get: (id: number) => Promise<T>
    create: (data: Partial<T>) => Promise<void>
    update: (id: number, data: Partial<T>) => Promise<void>
    delete: (id: number | number[]) => Promise<void>
  }
  message?: {
    success: string
    error: string
  }
}

export function useCRUD<T = any, P = any>(options: CRUDOptions<T, P>) {
  const loading = ref(false)
  const data = ref<T[]>([])
  const total = ref(0)
  const current = ref<T | null>(null)

  const fetchList = async (params: P) => {
    loading.value = true
    try {
      const res = await options.api.list(params)
      data.value = res.list
      total.value = res.total
    } catch (error: any) {
      showMessage(options.message?.error || '获取列表失败')
    } finally {
      loading.value = false
    }
  }

  const handleCreate = async (formData: Partial<T>) => {
    await options.api.create(formData)
    showMessage(options.message?.success || '创建成功')
    fetchList({} as P)
  }

  const handleUpdate = async (id: number, formData: Partial<T>) => {
    await options.api.update(id, formData)
    showMessage(options.message?.success || '更新成功')
    fetchList({} as P)
  }

  const handleDelete = async (id: number | number[]) => {
    await options.api.delete(id)
    showMessage(options.message?.success || '删除成功')
    fetchList({} as P)
  }

  return {
    loading,
    data,
    total,
    current,
    fetchList,
    handleCreate,
    handleUpdate,
    handleDelete,
  }
}
