import request from '../utils/request'

// 九字头维修订单 API
export function getJiuOrderList(params) {
  return request({
    url: '/jiu-orders',
    method: 'get',
    params
  })
}

export function getJiuOrderById(id) {
  return request({
    url: `/jiu-orders/${id}`,
    method: 'get'
  })
}

export function createJiuOrder(data) {
  return request({
    url: '/jiu-orders',
    method: 'post',
    data
  })
}

export function updateJiuOrder(id, data) {
  return request({
    url: `/jiu-orders/${id}`,
    method: 'put',
    data
  })
}

export function deleteJiuOrder(id) {
  return request({
    url: `/jiu-orders/${id}`,
    method: 'delete'
  })
}

export function exportJiuOrder(params) {
  return request({
    url: '/jiu-orders/export',
    method: 'post',
    params
  })
}

export function importJiuOrder(data) {
  return request({
    url: '/jiu-orders/import',
    method: 'post',
    data
  })
}
