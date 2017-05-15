import Vue from 'vue';
import Component from 'vue-class-component';
import { Chat, Message } from '../../model';
import io from 'socket.io-client';
import { ChatService, CHAT_ACCEPTED } from '../../services';

const POLL_INTERVAL = 1000; // 1 second
const TIMEOUT = 600000; // 10 minutes

@Component({
  template: require('./chat.html')
})
export class ChatComponent extends Vue {
  public id: string = '';
  public chat: Chat = null;
  public messages: Message[] =  [];
  public currMessage: string = '';
  public disabled: boolean = true;

  private socket: SocketIOClient.Socket = null;


  mounted() {
    console.log('ChatComponent Mounted!!!');
    this.id = this.$route.params.chatId;

    ChatService.pollChatStatus(this.id, POLL_INTERVAL, TIMEOUT)
      .then((chat: Chat) => {
        this.chat = chat;
        console.log(chat);
        if (chat.state === CHAT_ACCEPTED) {
          this.initSocket();
          this.disabled = false;
        } else {
          console.log('The chat has been declined');
        }
        console.log(chat);
      }, (error) => {
        console.log(error);
      });
  }

  initSocket() {
    const config = {
      host: process.env.API_HOST_NAME,
      port: process.env.API_PORT,
      path: `/api/chats/${this.id}/ws`,
      reconnectionAttempts: 10,
      'force new connection': true,
    };

    let socket = this.socket = io.connect(`${process.env.API_URL}/api/chats/${this.id}/ws`, config);

    socket.on('connect', () => {
      console.log('connected :D');
    });

    socket.on('visitorMessage', (message) => {
      this.messages.push(message);
      this.scrollToEnd();
    });

    socket.on('hostMessage', (message) => {
      this.messages.push(message);
      this.scrollToEnd();
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
    this.messages.push({
      chatId: this.id,
      source: 'visitor',
      fromName: 'You',
      content: this.currMessage,
    });
    this.currMessage = '';
    this.scrollToEnd();
  }

  onRequestTerminate() {
    this.socket.close();
    ChatService.terminateChat(this.id);
    this.$router.push({name: 'chats'});
  }

  scrollToEnd() {
    setTimeout(() => {
      const container = document.querySelector('body');
      container.scrollTop = container.scrollHeight;
    }, 100);
  }

}
