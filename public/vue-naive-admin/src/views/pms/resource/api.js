/**********************************
 * @Author: Ronnie Zhang
 * @LastEditor: Ronnie Zhang
 * @LastEditTime: 2024/04/01 15:52:04
 * @Email: zclzone@outlook.com
 * Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 **********************************/

import { request } from '@/utils'
import axios from 'axios'

export default {
  getMenuTree: () => request.get('/admin/permission/list'),
  getButtons: ({ parentId }) => request.get(`/admin/permission/button/${parentId}`),
  getComponents: () => axios.get(`${import.meta.env.VITE_PUBLIC_PATH}components.json`),
  addPermission: data => request.post('/admin/permission', data),
  savePermission: (id, data) => request.patch(`/admin/permission/${id}`, data),
  deletePermission: id => request.delete(`/admin/permission/${id}`),
}
