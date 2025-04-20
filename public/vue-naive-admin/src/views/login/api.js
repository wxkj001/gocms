/**********************************
 * @Author: Ronnie Zhang
 * @LastEditor: Ronnie Zhang
 * @LastEditTime: 2023/12/05 21:28:30
 * @Email: zclzone@outlook.com
 * Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 **********************************/

import { request } from '@/utils'

export default {
  toggleRole: data => request.post('/auth/role/toggle', data),
  login: (data, captchaId) => request.post('/admin/login', data, { needToken: false, headers: { 'X-Captcha-ID': captchaId } }),
  getUser: () => request.get('/admin/user/detail'),
  getCaptcha: () => request.get(`/admin/captcha?${Date.now()}`, { responseType: 'blob' }),
}
