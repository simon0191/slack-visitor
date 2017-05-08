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
  public currMessage: string = '';

  private socket: SocketIOClient.Socket = null;


  mounted() {
    this.id = this.$route.params.chatId;

    const config: SocketIOClient.ConnectOpts = {
      host: process.env.API_HOST_NAME,
      port: process.env.API_PORT,
      path: `/api/chats/${this.id}/ws`,
      reconnectionAttempts: 10,
    };

    let socket = this.socket = io.connect(`${process.env.API_URL}/api/chats/${this.id}/ws`, config);

    socket.on('connect', () => {
      console.log('connected :D');
    });

    socket.on('visitorMessage', (message) => {
      this.messages.push(message);
    });

    socket.on('hostMessage', (message) => {
      this.messages.push(message);
    });

    socket.on('received', (message) => {
      this.messages.push(message);
    });

    socket.on('welcome', (data) => {
      console.log('Welcome');
      console.log(data);
    });

    socket.on('diconnect', () => {
      console.log('diconnected D:');
    });
  }

  onSubmit() {
    this.socket.emit('visitorMessage', this.currMessage);
    this.messages.push(this.currMessage);
    this.currMessage = '';
  }

}
