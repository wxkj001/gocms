/**********************************
 * @Author: Ronnie Zhang
 * @LastEditor: Ronnie Zhang
 * @LastEditTime: 2023/12/05 21:29:51
 * @Email: zclzone@outlook.com
 * Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 **********************************/

import { request } from '@/utils'

export default {
  create: data => request.post('/admin/user', data),
  read: (params = {}) => request.get('/admin/user/list', { params }),
  update: data => request.patch(`/admin/user/${data.id}`, data),
  delete: id => request.delete(`/admin/user/${id}`),
  resetPwd: (id, data) => request.patch(`/admin/user/password/reset/${id}`, data),

  getAllRoles: () => request.get('/admin/role/tree?enable=1'),
}
