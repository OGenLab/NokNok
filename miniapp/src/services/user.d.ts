export interface LoginRequestReq {
  token: string;
  referralCode?: string;
  channel?: string;
}

export interface LoginRequestRsp extends BaseResponse {}

export interface BaseResponse {
  errCode: number;
  errMsg: string;
}

export type AccountInfo = {
  isPremium: boolean;
  boost: number;
  nextBoost: number;
  coins: number;
  invitedCount: number;
  invitedPremiumCount: number;
};

export type DrawInfo = {
  prizeType: number;
  prizeId: number;
  count: number;
  threshold: number;
  lastDrawTime: number;
};

export type NokInfo = {
  stamina: number;
  hitCount: number;
  lastHitTime: number;
};

export type TaskInfo = {
  id: number;
  type: number;
  schedule: number;
};

export type GetUserInfoRequestRsp = {
  accountInfo: AccountInfo;
  drawInfo: DrawInfo[];
  nokInfo: NokInfo;
  taskInfos: TaskInfo[];
  metaInfo: [];
} & BaseResponse;

export type HammerInfo = {
  hammerId: number;
  equipmentCount: number;
  count: number;
};

export type GetUserHammerRequestRsp = {
  hammers: HammerInfo[];
} & BaseResponse;

export type BoostInfo = {
  boost: number;
  neededCoins: number;
  consumeSamina: number;
  coinsRate: number;
};

export type GetBoostDetailRequestRsp = {
  boosts: BoostInfo[];
} & BaseResponse;

export type BoostRequestRsp = {
  coins: number;
  account: {
    boost: number;
    coins: number;
  };
} & BaseResponse;

export type InviteRequestRsp = {
  referralCode: string;
} & BaseResponse;
