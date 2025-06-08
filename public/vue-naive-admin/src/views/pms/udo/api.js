
import { request } from '@/utils'
import axios from 'axios'

export default {
  getUdoObject: () => request.get('/admin/udo/object/list'),
  addUdoObject: (udo) => request.post('/admin/udo/object/', udo),
  addUdoField: (udo) => request.post('/admin/udo/field/', udo),
  
}
