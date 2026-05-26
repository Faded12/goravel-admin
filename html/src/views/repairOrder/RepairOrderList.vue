<template>
  <div class="repair-order-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ $t('menu.repairOrder') }}</span>
        </div>
      </template>

      <SearchForm
        :model="searchForm"
        :fields="searchFields"
        :initial-values="initialSearchForm"
        i18n-prefix="repairOrder"
        @search="handleSearch"
        @reset="handleReset"
      >
        <template #extra-buttons>
          <el-button 
            type="primary" 
            @click="handleImport"
          >
            <el-icon><Upload /></el-icon>
            {{ $t('common.import') }}
          </el-button>
          <el-button 
            type="success" 
            :disabled="isExporting"
            :loading="isExporting"
            @click="handleExport"
          >
            <el-icon><Download /></el-icon>
            {{ $t('common.export') }}
          </el-button>
        </template>
      </SearchForm>

      <VxeTable
        ref="tableRef"
        :data="tableData"
        :loading="loading"
        :columns="tableColumns"
        :height="600"
        @sort-change="handleSortChange"
      />
      <Pagination
        v-model="pagination"
        :auto-load="true"
        :on-page-change="loadData"
      />

      <input
        ref="fileInputRef"
        type="file"
        accept=".csv"
        style="display: none"
        @change="handleFileChange"
      />
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Upload, Download } from '@element-plus/icons-vue'
import SearchForm from '../../components/SearchForm.vue'
import Pagination from '../../components/Pagination.vue'
import VxeTable from '../../components/VxeTable.vue'
import { useListPage } from '../../composables/useListPage'
import { getRepairOrderList, exportRepairOrder, importRepairOrder } from '../../api/repairOrder'

const { t } = useI18n()
const tableRef = ref(null)

const isImporting = ref(false)
const isExporting = ref(false)
const fileInputRef = ref(null)

const initialSearchForm = {
  works_order_number: '',
  notification_number: '',
  contract_number: '',
  works_order_status: '',
  start_date: '',
  end_date: ''
}

const {
  pagination,
  tableData,
  loading,
  searchForm,
  loadData,
  handleSearch,
  handleReset,
  handleSortChange
} = useListPage({
  fetchApi: getRepairOrderList,
  initialSearchForm,
  fieldMapping: {},
  defaultSort: 'id:desc',
  tableRef: computed(() => tableRef.value?.tableRef)
})

loadData()

const searchFields = computed(() => [
  {
    prop: 'works_order_number',
    label: t('repairOrder.works_order_number'),
    type: 'input',
    width: '180px'
  },
  {
    prop: 'notification_number',
    label: t('repairOrder.notification_number'),
    type: 'input',
    width: '180px'
  },
  {
    prop: 'contract_number',
    label: t('repairOrder.contract_number'),
    type: 'input',
    width: '150px'
  },
  {
    prop: 'works_order_status',
    label: t('repairOrder.works_order_status'),
    type: 'input',
    width: '180px'
  },
  {
    prop: 'start_date',
    label: t('repairOrder.start_date'),
    type: 'date',
    width: '150px'
  },
  {
    prop: 'end_date',
    label: t('repairOrder.end_date'),
    type: 'date',
    width: '150px'
  }
])

const tableColumns = computed(() => [
  {
    field: 'id',
    title: 'ID',
    width: 70,
    sortable: true
  },
  {
    field: 'works_order_number',
    title: t('repairOrder.works_order_number'),
    width: 150
  },
  {
    field: 'notification_number',
    title: t('repairOrder.notification_number'),
    width: 150
  },
  {
    field: 'works_order_description',
    title: t('repairOrder.works_order_description'),
    width: 300
  },
  {
    field: 'works_order_type',
    title: t('repairOrder.works_order_type'),
    width: 120
  },
  {
    field: 'contract_number',
    title: t('repairOrder.contract_number'),
    width: 120
  },
  {
    field: 'works_order_status',
    title: t('repairOrder.works_order_status'),
    width: 180
  },
  {
    field: 'functional_location',
    title: t('repairOrder.functional_location'),
    width: 140
  },
  {
    field: 'functional_location_description',
    title: t('repairOrder.functional_location_description'),
    width: 200
  },
  {
    field: 'start_date',
    title: t('repairOrder.start_date'),
    width: 110
  },
  {
    field: 'finish_date',
    title: t('repairOrder.finish_date'),
    width: 110
  },
  {
    field: 'estimated_total_costs',
    title: t('repairOrder.estimated_total_costs'),
    width: 130
  },
  {
    field: 'project_officer_post_description',
    title: t('repairOrder.project_officer_post_description'),
    width: 180
  }
])

const handleExport = async () => {
  if (isExporting.value) {
    return
  }

  isExporting.value = true
  
  try {
    const response = await exportRepairOrder(searchForm)
    const exportId = response.data?.export_id || response.data?.data?.export_id
    
    if (!exportId) {
      ElMessage.error(t('common.output_failed'))
      isExporting.value = false
      return
    }

    ElMessage.success(t('common.queued') || response.data?.message)
    
  } catch (error) {
    ElMessage.error(error.message || t('common.output_failed'))
  } finally {
    isExporting.value = false
  }
}

const handleImport = () => {
  if (fileInputRef.value) {
    fileInputRef.value.click()
  }
}

const handleFileChange = async (event) => {
  const file = event.target.files?.[0]
  if (!file) {
    return
  }

  if (!file.name.toLowerCase().endsWith('.csv')) {
    ElMessage.error(t('common.invalid_file_type') || '文件类型错误，请上传CSV文件')
    if (fileInputRef.value) {
      fileInputRef.value.value = ''
    }
    return
  }

  if (isImporting.value) {
    return
  }

  isImporting.value = true

  try {
    const response = await importRepairOrder(file)
    const result = response.data?.data || response.data

    if (result.success_count > 0) {
      ElMessage.success(
        t('common.import_success') || 
        `导入成功：成功 ${result.success_count} 条，失败 ${result.failed_count} 条`
      )
      
      if (result.failed_count > 0 && result.errors && result.errors.length > 0) {
        const errorMsg = result.errors.slice(0, 10).join('\n')
        if (result.errors.length > 10) {
          ElMessage.warning(`部分导入失败，前10条错误：\n${errorMsg}\n...`)
        } else {
          ElMessage.warning(`部分导入失败：\n${errorMsg}`)
        }
      }

      await loadData()
    } else {
      ElMessage.warning(t('common.import_no_data') || '没有成功导入任何数据')
      if (result.errors && result.errors.length > 0) {
        const errorMsg = result.errors.slice(0, 10).join('\n')
        ElMessage.error(`导入失败：\n${errorMsg}`)
      }
    }
  } catch (error) {
    ElMessage.error(error.message || t('common.import_failed'))
  } finally {
    isImporting.value = false
    if (fileInputRef.value) {
      fileInputRef.value.value = ''
    }
  }
}
</script>

<style scoped>
.repair-order-list {
  padding: 16px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
