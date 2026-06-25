import request from '../utils/request'

// 一字头维修订单 API
export function getYiOrderList(params) {
  return request({
    url: '/yi-orders',
    method: 'get',
    params
  })
}

export function getYiOrderById(id) {
  return request({
    url: `/yi-orders/${id}`,
    method: 'get'
  })
}

export function createYiOrder(data) {
  return request({
    url: '/yi-orders',
    method: 'post',
    data
  })
}

export function updateYiOrder(id, data) {
  return request({
    url: `/yi-orders/${id}`,
    method: 'put',
    data
  })
}

export function deleteYiOrder(id) {
  return request({
    url: `/yi-orders/${id}`,
    method: 'delete'
  })
}

export function exportYiOrder(params) {
  return request({
    url: '/yi-orders/export',
    method: 'post',
    params
  })
}

export function importYiOrder(data) {
  return request({
    url: '/yi-orders/import',
    method: 'post',
    data
  })
}
