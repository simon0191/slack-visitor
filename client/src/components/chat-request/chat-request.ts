import Vue from 'vue';
import Component from 'vue-class-component';
import { ChatService } from '../../services';
import { Chat } from '../../model';

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
        this.$router.push({name: 'chat', params: {chatId: chat.id}});
      });

  }
}
