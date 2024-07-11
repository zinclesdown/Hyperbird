// 用户鉴权/登陆/刷新相关方法。
import { userState } from './stores/user-token';

const userStore = userState();

export function printUserStatus() {
  console.log('当前登陆状态：', userStore.GetCurLoginStatus);
}

export function DebugFakeLogin() {
  userStore.DebugFakeLogin();
}
