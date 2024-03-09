import service from '@/utils/request'

export const getWechatConfig = (data) => {
  return service({
    url: '/wechat/private/config',
    method: 'get',
    data
  })
}

export const updateWechatConfig = (data) => {
  return service({
    url: '/wechat/private/config',
    method: 'put',
    data
  })
}

// 获取令牌
// type = 'miniProgram' | 'officialAccount'
export const getAccessToken = (type = 'miniProgram' | 'officialAccount') => {
  console.log(type)
  return service({
    url: '/wechat/private/token',
    method: 'get',
    params: {
      type
    }
  })
}

