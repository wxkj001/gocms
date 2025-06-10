import axios from 'axios'
// import { request } from '@/utils'

export default {
  hot: data => axios.post(`https://api-hot.imsyy.top/${data}`),
//   changePassword: data => request.post(`https://api-hot.imsyy.top/${data}`, data),
}
