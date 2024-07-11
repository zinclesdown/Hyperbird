import { defineStore } from 'pinia';

// 用户状态
export const userState = defineStore('counter', {
  state: () => ({
    curUserName: '',
    curUserToken: '',
  }),
  getters: {
    GetCurLoginStatus: (state) => {
      if (state.curUserName === '' || state.curUserToken === '') {
        return false;
      } else {
        return true;
      }
    },
  },
  actions: {
    SetName() {
      this.curUserName = 'test Name!';
    },

    DebugFakeLogin() {
      this.curUserName = 'debugtestffakeLoginuser';
      this.curUserToken = 'dhsaidhsdsocfdsojisjfoiljifjsd';
      console.log('尝试调试登陆了：DebugFakeLogin: ', this.curUserName, this.curUserToken);
    },
  },
});
