import { createPromise } from "src/utils/helper";
import nano from "./nano";
import { BaseResponse } from "src/services/user.d";

// docs: https://github.com/nano-ecosystem/nano-websocket-client

const host = "9.134.91.143";
const port = 8790;

const [initid, resolveInitid] = createPromise();
nano.init(
  {
    host: host,
    port: port,
  },
  function () {
    console.log("init nano success");
    resolveInitid(true);
  }
);

export const request = async <T, J extends BaseResponse>(
  route: string,
  msg: T
): Promise<J> => {
  await initid;
  return new Promise((resolve, reject) => {
    try {
      nano.request(route, msg, (data: J) => {
        const { errCode, errMsg, ...rest } = data;

        if (errCode === 0) {
          resolve(rest as J);
        } else {
          reject(new Error(errMsg));
        }
      });
    } catch (e) {
      reject(e);
    }
  });
};