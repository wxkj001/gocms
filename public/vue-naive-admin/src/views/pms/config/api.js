import { request } from '@/utils'

export default {
  list: data => request.get('/admin/config/list', data),
  update: data => request.patch(`/admin/config/${data.config_key}`, data),
}
