import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";
import installVant from "./plugins/vant";

const app = createApp(App);
installVant(app);
app.use(router).mount("#app");
