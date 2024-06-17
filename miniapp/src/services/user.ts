import { request } from "src/lib/request";
import {
  LoginRequestReq,
  LoginRequestRsp,
  GetUserInfoRequestRsp,
  GetUserHammerRequestRsp,
  GetBoostDetailRequestRsp,
  BoostRequestRsp,
  InviteRequestRsp,
} from "./user.d";

export const login = async (params: LoginRequestReq) => {
  return request<LoginRequestReq, LoginRequestRsp>("user.login", params);
};

export const getUserInfo = async () => {
  return request<{}, GetUserInfoRequestRsp>("user.getinfo", {});
};

export const getUserHammer = async () => {
  return request<{}, GetUserHammerRequestRsp>("user.gethammer", {});
};

export const getBoostInfo = async () => {
  return request<{}, GetBoostDetailRequestRsp>("user.getboost", {});
};

export const boost = async () => {
  return request<{}, BoostRequestRsp>("user.boost", {});
};

export const invite = async () => {
  return request<{}, InviteRequestRsp>("user.inviteinfo", {});
};
