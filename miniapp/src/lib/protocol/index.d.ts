export interface Package {
  TYPE_HANDSHAKE: number;
  TYPE_HANDSHAKE_ACK: number;
  TYPE_HEARTBEAT: number;
  TYPE_DATA: number;
  TYPE_KICK: number;
}

export interface Message {
  TYPE_REQUEST: number;
  TYPE_NOTIFY: number;
  TYPE_RESPONSE: number;
  TYPE_PUSH: number;
}

export function strencode(str: string): Uint8Array;
export function strdecode(buffer: Uint8Array): string;

export function encode(type: number, body?: Uint8Array): Uint8Array;
export function decode(
  buffer: Uint8Array
):
  | { type: number; body?: Uint8Array }
  | Array<{ type: number; body?: Uint8Array }>;

export function encode(
  id: number,
  type: number,
  compressRoute: boolean,
  route: number | string,
  msg: Uint8Array
): Uint8Array;
export function decode(buffer: Uint8Array): {
  id: number;
  type: number;
  compressRoute: boolean;
  route: number | string;
  body: Uint8Array;
};

export const Package: Package;
export const Message: Message;
