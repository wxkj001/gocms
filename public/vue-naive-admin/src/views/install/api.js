/**********************************
 * @Author: Ronnie Zhang
 * @LastEditor: Ronnie Zhang
 * @LastEditTime: 2023/12/05 21:28:30
 * @Email: zclzone@outlook.com
 * Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 **********************************/

import { request } from '@/utils'

export default {
  register: (data, captchaId) => request.post('/install', data, { needToken: false, headers: { 'X-Captcha-ID': captchaId } }),
  getCaptcha: () => request.get(`/admin/captcha?${Date.now()}`, { responseType: 'blob' }),
}
