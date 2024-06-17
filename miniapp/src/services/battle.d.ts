import { BaseResponse } from "./user.d";

export type Gopher = {
  id: number;
  season: number;
  name: string;
  imageUrl: string;
  probability: number;
};

export type GetGophersRequestRsp = {
  gophers: Gopher[];
} & BaseResponse;

export type Hammer = {
  id: number;
  season: number;
  name: string;
  image: string;
  probability: number;
};

export type GetHammersRequestRsp = {
  hammers: Hammer[];
} & BaseResponse;

export type GetTasksRequestReq = {
  type?: number;
};

export type TaskInfo = {
  id: number;
  season: number;
  type: number;
  threshold: number;
  reward: number;
  startTime: number;
  endTime: number;
};

export type GetTasksRequestRsp = {
  taskInfos: TaskInfo[];
} & BaseResponse;

export type GetGopherLocationRequestRsp = {
  gopherLocations: [];
} & BaseResponse;

export type HammerGopherRequestReq = {
  round: string;
  index: number;
};

type NokInfo = {
  stamina: number;
  hitCount: number;
  lastHitTime: number;
};

export type HammerGopherRequestRsp = {
  coins: number;
  nokInfo: NokInfo[];
  account: {
    coins: number;
  };
} & BaseResponse;

export type GetHammerBoardRequestReq = {
  type: number;
  limit?: number;
  offset?: number;
};

type BoardInfo = {
  rank: number;
  id: number;
  count: number;
};

export type GetHammerBoardRequestRsp = {
  leaderboardInfo: BoardInfo[];
} & BaseResponse;

export type PrizeDrawRequestReq = {
  prizeType: number;
};

export type PrizeDrawRequestRsp = {
  drawInfo: {
    prizeId: number;
    count: number;
    threshold: number;
    lastDrawTime: number;
  };
} & BaseResponse;

export type EquipHammerRequestReq = {
  hammerId: number;
};

export interface EquipHammerRequestRsp extends BaseResponse {}
