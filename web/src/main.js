import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";
import installVant from "./plugins/vant";

const app = createApp(App);
installVant(app);
app.use(store).use(router).mount("#app");
