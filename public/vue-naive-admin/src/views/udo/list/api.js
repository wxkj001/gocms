import { request } from '@/utils'

export default {
  fields: (code, params = {}) => request.get(`/admin/udo/field/code/${code}`, { params }),
  list: (code, params = {}) => request.get(`/admin/udo/data/${code}/list`, { params }),
}
