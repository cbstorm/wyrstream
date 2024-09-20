export interface INATSConfig {
  LoadNATSCoreConfig: () => void;
  NATS_CORE_HOST: string;
  NATS_CORE_PORT: string;
  NATS_CORE_USERNAME: string;
  NATS_CORE_PASSWORD: string;
}

export interface ISubscriber<T> {
  Start(): Promise<void>;
  Stop(): Promise<void>;
  Process(msg: IRequestMessage<T>): Promise<void>;
}

export interface IRequestMessage<T> {
  Data(): T;
}

export interface IResponseMessage<T> {
  Data(): T;
  Error(): any;
}
