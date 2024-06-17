import { request } from "src/lib/request";
import {
  GetGophersRequestRsp,
  GetHammersRequestRsp,
  GetTasksRequestReq,
  GetTasksRequestRsp,
  GetGopherLocationRequestRsp,
  HammerGopherRequestReq,
  HammerGopherRequestRsp,
  GetHammerBoardRequestReq,
  GetHammerBoardRequestRsp,
  PrizeDrawRequestReq,
  PrizeDrawRequestRsp,
  EquipHammerRequestReq,
  EquipHammerRequestRsp,
} from "./battle.d";

export const getGophers = async () => {
  return request<{}, GetGophersRequestRsp>("battle.gophers", {});
};

export const getHammers = async () => {
  return request<{}, GetHammersRequestRsp>("battle.entities", {});
};

export const getTasks = async (params?: GetTasksRequestReq) => {
  return request<GetTasksRequestReq, GetTasksRequestRsp>(
    "battle.tasks",
    params ? params : {}
  );
};

export const getGopherLocations = async () => {
  return request<{}, GetGopherLocationRequestRsp>("battle.getgophers", {});
};

export const hammerGopher = async (params: HammerGopherRequestReq) => {
  return request<HammerGopherRequestReq, HammerGopherRequestRsp>(
    "battle.hammergopher",
    params
  );
};

export const getHammerBoard = async (params: GetHammerBoardRequestReq) => {
  return request<GetHammerBoardRequestReq, GetHammerBoardRequestRsp>(
    "battle.hammerleaderboard",
    params
  );
};

export const prizeDraw = async (params: PrizeDrawRequestReq) => {
  return request<PrizeDrawRequestReq, PrizeDrawRequestRsp>(
    "battle.prizedraw",
    params
  );
};

export const equipHammer = async (params: EquipHammerRequestReq) => {
  return request<EquipHammerRequestReq, EquipHammerRequestRsp>(
    "battle.equiphammer",
    params
  );
};
