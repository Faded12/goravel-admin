import request from '@/utils/request'

export function getRepairOrderList(params) {
  return request({
    url: '/repair-orders',
    method: 'get',
    params,
  })
}

export function getRepairOrder(id) {
  return request({
    url: `/repair-orders/${id}`,
    method: 'get',
  })
}

export function exportRepairOrder(params) {
  return request({
    url: '/repair-orders/export',
    method: 'post',
    data: params
  })
}

export function importRepairOrder(file) {
  const formData = new FormData()
  formData.append('file', file)
  return request({
    url: '/repair-orders/import',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}
