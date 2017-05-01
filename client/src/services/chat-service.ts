import { Chat, Error } from '../model';
import axios from 'axios';
import { AxiosResponse } from 'axios';

export const CHAT_PENDING = 'pending';
export const CHAT_ACCEPTED = 'accepted';

export class ChatService {
  public static create(visitorName: string, subject: string): PromiseLike<Chat> {

    return axios.post(`${process.env.API_URL}/api/chats`, {visitorName, subject})
        .then<Chat>((response: AxiosResponse) => {
          const chat = response.data as Chat;
          return chat;
        })
        .catch<Error>((error) => {
          // TODO: properly handle different type of errors
          throw {
            id: 'chats.create.unkown_error',
            data: {
              parent: error
            }
          };
        });
  }

  public static get(id: string): PromiseLike<Chat> {

    return axios.get(`${process.env.API_URL}/api/chats/${id}`)
      .then<Chat|any>((response: AxiosResponse) => {
        const chat = response.data as Chat;
        return chat;
      })
      .catch((error) => {
        throw {
          id: 'chats.get.unkown_error',
          data: {
            parent: error
          }
        };
      });
  }

  public static pollChatStatus(chatId: string, pollInterval: number, timeout: number): Promise<Chat> {
    return new Promise<Chat>((resolve: (chat: Chat) => void, reject: (error: any) => void) => {

      const intervalId = setInterval(() => {
        ChatService.get(chatId).then((chat: Chat) => {
          if (chat.state !== CHAT_PENDING) {
            resolve(chat);
            clearInterval(intervalId);
            clearTimeout(timeoutId);
          }
        });
      }, pollInterval);

      const timeoutId = setTimeout(() => {
        reject({id: 'chats.poll_chat_status.timeout'});
        clearInterval(intervalId);
      }, timeout);

    });
  }
}
