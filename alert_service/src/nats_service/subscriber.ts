import { IRequestMessage, ISubscriber } from './nats_service.interface';

export class Subscriber<T> implements ISubscriber<T> {
  async Start() {}
  async Stop() {}
  async Process(msg: IRequestMessage<T>) {}
}
