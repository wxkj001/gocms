import { request } from '@/utils'

export default {
  fields: (code, params = {}) => request.get(`/admin/udo/field/code/${code}`, { params }),
  list: (code, params = {}) => request.get(`/admin/udo/data/${code}/list`, { params }),
  create: (code, params = {}) => request.post(`/admin/udo/data/${code}`, params),
  update: (code,id, params = {}) => request.put(`/admin/udo/data/${code}/${id}`, params),
  delete: (code, id) => request.delete(`/admin/udo/data/${code}/${id}`),
}
