interface Callbacks {
  [key: number]: (...args: any[]) => void;
}

interface Handlers {
  [key: number]: (body: any) => void;
}

interface Dict {
  [key: string]: number;
}

interface Abbrs {
  [key: number]: string;
}

interface Params {
  host: string;
  port?: number;
  path?: string;
  user?: any;
  encrypt?: boolean;
  reconnect?: boolean;
  maxReconnectAttempts?: number;
  handshakeCallback?: (user: any) => void;
  encode?: (reqId: number, route: string, msg: any) => Uint8Array;
  decode?: (data: Uint8Array) => any;
}

const nano: {
  init: (params: Params, cb: (socket: WebSocket) => void) => void;
  disconnect: () => void;
  request: (route: string, msg: any, cb: (data: any) => void) => void;
  notify: (route: string, msg: any) => void;
  on: (event: string, fn: (...args: any[]) => void) => void;
  off: (event?: string, fn?: (...args: any[]) => void) => void;
  once: (event: string, fn: (...args: any[]) => void) => void;
  emit: (event: string, ...args: any[]) => void;
  listeners: (event: string) => Function[];
  hasListeners: (event: string) => boolean;
} = Object.create(EventEmitter.prototype);

export default nano;
