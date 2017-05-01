import Vue from 'vue';
import Component from 'vue-class-component';
import { ChatService, CHAT_ACCEPTED } from '../../services';
import { Chat } from '../../model';

const POLL_INTERVAL = 1000; // 1 second
const TIMEOUT = 60000; // 1 minute

@Component({
  template: require('./chat-request.html')
})
export class ChatRequestComponent extends Vue {
  public visitorName: string = '';
  public subject: string = '';

  public onSubmit() {
    ChatService
      .create(this.visitorName, this.subject)
      .then((chat: Chat) => {
        return ChatService.pollChatStatus(chat.id, POLL_INTERVAL, TIMEOUT);
      })
      .then((chat: Chat) => {
        console.log(chat);
        console.log(chat);
        if (chat.state === CHAT_ACCEPTED) {
          this.$router.push({name: 'chat', params: {chatId: chat.id}});
        } else {
          console.log('The chat has been declined');
        }
        console.log(chat);
      }, (error) => {
        console.log(error);
      });

  }
}
