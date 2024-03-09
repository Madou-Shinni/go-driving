<template>
  <div class="bg-white p-8 shadow-md rounded">
    <h1>微信配置</h1>
    <hr class="my-4">
    <el-form :model="wechatConfig" :rules="rules" ref="wechatConfigForm" label-width="120px">
      <el-form-item label="小程序配置" prop="miniProgramEnabled">
        <el-switch v-model="wechatConfig.miniProgramEnabled" active-text="启用" inactive-text="禁用" />
      </el-form-item>
      <el-form-item v-if="wechatConfig.miniProgramEnabled">
        <el-form-item label="小程序AppID" prop="miniProgram.appId">
          <el-input v-model="wechatConfig.miniProgram.appId"></el-input>
        </el-form-item>
        <el-form-item label="小程序Secret" prop="miniProgram.appSecret">
          <el-input v-model="wechatConfig.miniProgram.appSecret"></el-input>
        </el-form-item>
        <el-form-item style="margin-left: 20px;">
          <el-button type="primary" @click="getMiniProgramAccessToken">获取小程序AccessToken</el-button>
        </el-form-item>
      </el-form-item>

      <el-form-item label="公众号配置" prop="officialAccountEnabled">
        <el-switch v-model="wechatConfig.officialAccountEnabled" active-text="启用" inactive-text="禁用" />
      </el-form-item>
      <el-form-item v-if="wechatConfig.officialAccountEnabled">
        <el-form-item label="公众号AppID" prop="officialAccount.appId">
          <el-input v-model="wechatConfig.officialAccount.appId"></el-input>
        </el-form-item>
        <el-form-item label="公众号Secret" prop="officialAccount.appSecret">
          <el-input v-model="wechatConfig.officialAccount.appSecret"></el-input>
        </el-form-item>
        <el-form-item style="margin-left: 20px;">
          <el-button type="primary" @click="getOfficialAccountAccessToken">获取公众号AccessToken</el-button>
        </el-form-item>
      </el-form-item>

      <el-form-item>
        <el-button type="primary" @click="updateConfig">更新配置</el-button>
      </el-form-item>
    </el-form>

    <!-- 弹窗 -->
    <el-dialog v-model:="dialog.visible" title="提示" :center="true" :before-close="handleClose">
      <span class="text-3xl">{{ dialog.content }}</span>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { ElSwitch, ElForm, ElFormItem, ElInput, ElButton, ElMessage } from 'element-plus';
import {getAccessToken, getWechatConfig, updateWechatConfig} from "@/plugin/wechat/api/wechat";

defineOptions({
  name: 'Wechat',
})

const wechatConfigForm = ref(null)

const wechatConfig = ref({
  miniProgramEnabled: false,
  miniProgram: {
    appId: '',
    appSecret: ''
  },
  officialAccountEnabled: false,
  officialAccount: {
    appId: '',
    appSecret: ''
  }
});

const dialog = ref({
  visible: false,
  content: ''
})

const handleClose = () => {
  dialog.value.visible = false
  dialog.value.content = ''
}

const rules = {
  miniProgram: {
    appId: [{ required: true, message: '请输入小程序AppID', trigger: 'blur' }],
    appSecret: [{ required: true, message: '请输入小程序Secret', trigger: 'blur' }]
  },
  officialAccount: {
    appId: [{ required: true, message: '请输入公众号AppID', trigger: 'blur' }],
    appSecret: [{ required: true, message: '请输入公众号Secret', trigger: 'blur' }]
  }
};

// 获取配置
const getConfig = async() => {
  const res = await getWechatConfig()
  if (res.code === 0) {
    wechatConfig.value = res.data
  }
}
getConfig()

// 修改配置
const updateConfig = async() => {
  wechatConfigForm.value.validate(async valid => {
    if (valid) {
      const res = await updateWechatConfig(wechatConfig.value)
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '更新成功!'
        })
      }else {
        ElMessage({
          type: 'error',
          message: '更新失败'
        })
      }
    }
  })
};

const getMiniProgramAccessToken = async() => {
  const res = await getAccessToken('miniProgram')
  if (res.code === 0) {
    // 实现获取小程序 access token 的逻辑
    console.log('获取小程序AccessToken');
    dialog.value.visible = true
    dialog.value.content = res.data.accessToken
  }
}

const getOfficialAccountAccessToken = async() => {
  const res = await getAccessToken('officialAccount')
  if (res.code === 0) {
    // 实现获取小程序 access token 的逻辑
    dialog.value.visible = true
    dialog.value.content = res.data.accessToken
  }
}

</script>
