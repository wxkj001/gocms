import { request } from '@/utils'

export default {
  getUdoObject: () => request.get('/admin/udo/object/list'),
  getUdoObjectDel: id => request.delete(`/admin/udo/object/${id}`),
  addUdoObject: udo => request.post('/admin/udo/object/', udo),
  addUdoField: udo => request.post('/admin/udo/field/', udo),
  getUdoFieldList: data => request.get(`/admin/udo/field/list`, { params: data }),
  getUdoFieldUp: data => request.put(`/admin/udo/field`, data),
  getUdoFieldDel: id => request.delete(`/admin/udo/field/${id}`),

}
