/**********************************
 * @Author: Ronnie Zhang
 * @LastEditor: Ronnie Zhang
 * @LastEditTime: 2023/12/05 21:29:27
 * @Email: zclzone@outlook.com
 * Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 **********************************/

import { request } from '@/utils'

export default {
  create: data => request.post('/admin/role', data),
  read: (params = {}) => request.get('/admin/role/list', { params }),
  update: data => request.patch(`/admin/role/${data.id}`, data),
  delete: id => request.delete(`/admin/role/${id}`),

  getAllPermissionTree: () => request.get('/admin/permission/tree'),
  getAllUsers: (params = {}) => request.get('/admin/user', { params }),
  addRoleUsers: (roleId, data) => request.patch(`/admin/role/users/add/${roleId}`, data),
  removeRoleUsers: (roleId, data) => request.patch(`/admin/role/users/remove/${roleId}`, data),
}
