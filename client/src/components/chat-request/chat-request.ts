import Vue from 'vue';
import Component from 'vue-class-component';
import { ChatService } from '../../services';
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
        return ChatService.pollChatStatus(chat.id, POLL_INTERVAL, TIMEOUT)
      })
      .then((chat: Chat) => {
        console.log(chat);
      }, (error) => {
        console.log(error);
      });

  }
}
