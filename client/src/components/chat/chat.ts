import Vue from 'vue';
import Component from 'vue-class-component';
import { Chat } from '../../model';
import io from 'socket.io-client';

const POLL_INTERVAL = 1000; // 1 second
const TIMEOUT = 60000; // 1 minute

@Component({
  template: require('./chat.html')
})
export class ChatComponent extends Vue {
  public id: string = '';
  public messages: string[] =  [];

  mounted() {
    this.id = this.$route.params.chatId;

    const config = {
      host: process.env.API_HOST_NAME,
      port: process.env.API_PORT,
      path: `/api/chats/${this.id}/ws`,
    };
    let socket = io.connect(`${process.env.API_URL}/api/chats/${this.id}/ws`, config);

    socket.on('connect', function() {
      console.log('connected :D');
      socket.emit('message', 'Hola');
      socket.on('welcome', function(data) {
        console.log('Welcome');
        console.log(data);
      });
      socket.on('diconnect', function() {
        console.log('diconnected D:');
      });
    });

  }

}
