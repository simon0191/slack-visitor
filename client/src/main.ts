import * as Vue from 'vue';
import VueRouter from 'vue-router';
import { ChatRequestComponent } from './components/chat-request';
import { NavbarComponent } from './components/navbar';
import { ChatComponent } from './components/chat';

// register the plugin
Vue.use(VueRouter);

let router = new VueRouter({
  routes: [
    { name: 'chats', path: '/', component: ChatRequestComponent },
    { name: 'chat', path: '/chat/:chatId', component: ChatComponent },
  ]
});

new Vue({
  el: '#app-main',
  router: router,
  components: {
    'navbar': NavbarComponent
  }
});
